package service

import (
	"errors"

	"github.com/ARTMUC/magic-video/internal/repository"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var ErrRecordNotFound = errors.New("record not found")

type BaseCrud[T any] interface {
	Create(*T, repository.WriteOptions) error
	List(repository.ReadOptions, repository.Pagination) (*repository.PaginatedResult[T], error)
	Get(uuid.UUID, repository.ReadOptions) (*T, error)
	Update(uuid.UUID, *T, repository.WriteOptions) error
	Delete(uuid.UUID, repository.WriteOptions) error
}

type baseCrud[T any, R repository.BaseRepository[T]] struct {
	repository R
}

func newBaseCrud[T any, R repository.BaseRepository[T]](repository R) BaseCrud[T] {
	return &baseCrud[T, R]{repository: repository}
}

func (b *baseCrud[T, R]) Update(UUID uuid.UUID, input *T, options repository.WriteOptions) error {
	entity, err := b.repository.FindOne(
		repository.ReadOptions{
			Scopes: []repository.Scope{
				repository.WithUUID(UUID),
			},
		},
	)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
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

func (b *baseCrud[T, R]) Create(entity *T, options repository.WriteOptions) error {
	err := b.repository.Create(options, entity)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCrud[T, R]) Delete(UUID uuid.UUID, options repository.WriteOptions) error {
	entity, err := b.repository.FindOne(
		repository.ReadOptions{
			Scopes: []repository.Scope{
				repository.WithUUID(UUID),
			},
		},
	)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
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

func (b *baseCrud[T, R]) Get(UUID uuid.UUID, options repository.ReadOptions) (*T, error) {
	options.Scopes = append(options.Scopes, repository.WithUUID(UUID))
	entity, err := b.repository.FindOne(options)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return entity, nil
}

func (b *baseCrud[T, R]) List(options repository.ReadOptions, pagination repository.Pagination) (*repository.PaginatedResult[T], error) {
	return b.repository.Paginate(options, pagination)
}
