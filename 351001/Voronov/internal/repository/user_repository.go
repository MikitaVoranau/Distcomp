package repository

import (
	apperrors "Voronov/internal/errors"
	"Voronov/internal/model"
	"context"
	std_errors "errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgUserRepository{pool: pool}
}

func (r *pgUserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := "SELECT id, login, password, firstname, lastname FROM distcomp.tbl_user WHERE id = $1"
	row := r.pool.QueryRow(ctx, query, id)
	var u model.User
	err := row.Scan(&u.ID, &u.Login, &u.Password, &u.Firstname, &u.Lastname)
	if std_errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	query := "SELECT id, login, password, firstname, lastname FROM distcomp.tbl_user WHERE login = $1"
	row := r.pool.QueryRow(ctx, query, login)
	var u model.User
	err := row.Scan(&u.ID, &u.Login, &u.Password, &u.Firstname, &u.Lastname)
	if std_errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) FindAll(ctx context.Context, opts *QueryOptions) ([]*model.User, int64, error) {
	if opts == nil {
		opts = NewQueryOptions()
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM distcomp.tbl_user").Scan(&total); err != nil {
		return nil, 0, err
	}

	orderField := "id"
	orderDir := "ASC"
	if opts.Sort != nil {
		if opts.Sort.Field != "" {
			orderField = opts.Sort.Field
		}
		if opts.Sort.Direction == "DESC" {
			orderDir = "DESC"
		}
	}

	offset := (opts.Pagination.Page - 1) * opts.Pagination.PageSize
	query := fmt.Sprintf(
		"SELECT id, login, password, firstname, lastname FROM distcomp.tbl_user ORDER BY %s %s LIMIT $1 OFFSET $2",
		orderField, orderDir,
	)

	rows, err := r.pool.Query(ctx, query, opts.Pagination.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Login, &u.Password, &u.Firstname, &u.Lastname); err != nil {
			return nil, 0, err
		}
		items = append(items, &u)
	}
	if items == nil {
		items = []*model.User{}
	}
	return items, total, nil
}

func (r *pgUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	fmt.Printf("🛠 [DB] Попытка создания пользователя: %s\n", user.Login)
	query := "INSERT INTO distcomp.tbl_user (login, password, firstname, lastname) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int64
	err := r.pool.QueryRow(ctx, query, user.Login, user.Password, user.Firstname, user.Lastname).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if std_errors.As(err, &pgErr) && pgErr.Code == "23505" {
			fmt.Printf("❌ [DB] Дубликат логина: %s\n", user.Login)
			return nil, apperrors.ErrDuplicate
		}
		return nil, apperrors.FromDBError(err)
	}
	fmt.Printf("✅ [DB] Пользователь создан с ID: %d\n", id)
	user.ID = id
	return user, nil
}

func (r *pgUserRepository) Update(ctx context.Context, id int64, user *model.User) (*model.User, error) {
	query := "UPDATE distcomp.tbl_user SET login = $1, password = $2, firstname = $3, lastname = $4 WHERE id = $5"
	result, err := r.pool.Exec(ctx, query, user.Login, user.Password, user.Firstname, user.Lastname, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if std_errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.ErrDuplicate
		}
		return nil, apperrors.FromDBError(err)
	}
	if result.RowsAffected() == 0 {
		return nil, apperrors.ErrNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *pgUserRepository) Delete(ctx context.Context, id int64) error {
	fmt.Printf("🛠 [DB] Попытка удаления пользователя ID: %d\n", id)
	result, err := r.pool.Exec(ctx, "DELETE FROM distcomp.tbl_user WHERE id = $1", id)
	if err != nil {
		fmt.Printf("❌ [DB] Ошибка при удалении пользователя: %v\n", err)
		return apperrors.FromDBError(err)
	}

	rows := result.RowsAffected()
	fmt.Printf("✅ [DB] Удалено строк: %d\n", rows)

	if rows == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
