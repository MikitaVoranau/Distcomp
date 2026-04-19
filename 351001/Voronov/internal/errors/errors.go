package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type AppError struct {
	Code       int    `json:"errorCode"`
	Message    string `json:"errorMessage"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

func NewAppError(httpStatus int, code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// FromDBError maps a pgx/pgconn error to an AppError by Postgres error code.
func FromDBError(err error) *AppError {
	// 1. ИСПРАВЛЕНИЕ ЗДЕСЬ:
	// Если ошибка уже является нашим AppError (например, репозиторий вернул ErrDuplicate),
	// просто возвращаем её как есть, ничего не оборачивая!
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr
	}

	// 2. Распаковываем ошибку базы данных
	var pgErr *pgconn.PgError
	if stderrors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrDuplicate
		case "23503": // foreign_key_violation
			return ErrDuplicate
		case "23502": // not_null_violation
			return ErrBadRequest
		default:
			// Дебаг: логируем неизвестные ошибки БД
			return NewAppError(http.StatusInternalServerError, 50001, "DB Error: "+pgErr.Code+" - "+pgErr.Message)
		}
	}

	// 3. Неизвестная ошибка
	return NewAppError(http.StatusInternalServerError, 50001, "Internal: "+err.Error())
}

var (
	ErrNotFound   = NewAppError(http.StatusNotFound, 40401, "Resource not found")
	ErrBadRequest = NewAppError(http.StatusBadRequest, 40001, "Invalid request")
	ErrForbidden  = NewAppError(http.StatusForbidden, 40301, "Forbidden")
	// Убедись, что здесь StatusForbidden (403), так как тесты этого ждут для дубликатов/FK
	ErrDuplicate = NewAppError(http.StatusForbidden, 40301, "Action forbidden or conflict")
	ErrInternal  = NewAppError(http.StatusInternalServerError, 50001, "Internal server error")
)
