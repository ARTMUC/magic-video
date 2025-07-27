package base

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type BaseCrud[T any] interface {
	Create(*T, WriteOptions) error
	List(ReadOptions, Pagination) (*PaginatedResult[T], error)
	Get(uuid.UUID, ReadOptions) (*T, error)
	Update(uuid.UUID, *T, WriteOptions) error
	Delete(uuid.UUID, WriteOptions) error
}

type baseCrud[T any, R BaseRepository[T]] struct {
	repository R
}

func NewBaseCrud[T any, R BaseRepository[T]](repository R) BaseCrud[T] {
	return &baseCrud[T, R]{repository: repository}
}

func (b *baseCrud[T, R]) Update(id uuid.UUID, input *T, options WriteOptions) error {
	entity, err := b.repository.FindOne(
		ReadOptions{
			Scopes: []Scope{
				WithID(id),
			},
		},
	)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}

	err = copier.Copy(entity, input)
	if err != nil {
		return err
	}

	err = b.repository.Update(options, entity)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCrud[T, R]) Create(entity *T, options WriteOptions) error {
	err := b.repository.Create(options, entity)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCrud[T, R]) Delete(id uuid.UUID, options WriteOptions) error {
	entity, err := b.repository.FindOne(
		ReadOptions{
			Scopes: []Scope{
				WithID(id),
			},
		},
	)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}

	err = b.repository.Delete(options, entity)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCrud[T, R]) Get(id uuid.UUID, options ReadOptions) (*T, error) {
	options.Scopes = append(options.Scopes, WithID(id))
	entity, err := b.repository.FindOne(options)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return entity, nil
}

func (b *baseCrud[T, R]) List(options ReadOptions, pagination Pagination) (*PaginatedResult[T], error) {
	return b.repository.Paginate(options, pagination)
}
