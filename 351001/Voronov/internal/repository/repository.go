package repository

type CRUDRepository[T any] interface {
	FindByID(id int64) (*T, error)
	FindAll() ([]*T, error)
	Create(entity *T) (*T, error)
	Update(id int64, entity *T) (*T, error)
	Delete(id int64) error
}
