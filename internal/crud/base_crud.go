package crud

import (
	"errors"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var ErrRecordNotFound = errors.New("record not found")

type BaseCrud[T any] interface {
	Create(*T, base.WriteOptions) error
	List(base.ReadOptions, base.Pagination) (*base.PaginatedResult[T], error)
	Get(uuid.UUID, base.ReadOptions) (*T, error)
	Update(uuid.UUID, *T, base.WriteOptions) error
	Delete(uuid.UUID, base.WriteOptions) error
}

type baseCrud[T any, R base.BaseRepository[T]] struct {
	repository R
}

func newBaseCrud[T any, R base.BaseRepository[T]](repository R) BaseCrud[T] {
	return &baseCrud[T, R]{repository: repository}
}

func (b *baseCrud[T, R]) Update(UUID uuid.UUID, input *T, options base.WriteOptions) error {
	entity, err := b.repository.FindOne(
		base.ReadOptions{
			Scopes: []base.Scope{
				base.WithUUID(UUID),
			},
		},
	)
	if err != nil {
		if errors.Is(err, base.ErrRecordNotFound) {
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

func (b *baseCrud[T, R]) Create(entity *T, options base.WriteOptions) error {
	err := b.repository.Create(options, entity)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCrud[T, R]) Delete(UUID uuid.UUID, options base.WriteOptions) error {
	entity, err := b.repository.FindOne(
		base.ReadOptions{
			Scopes: []base.Scope{
				base.WithUUID(UUID),
			},
		},
	)
	if err != nil {
		if errors.Is(err, base.ErrRecordNotFound) {
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

func (b *baseCrud[T, R]) Get(UUID uuid.UUID, options base.ReadOptions) (*T, error) {
	options.Scopes = append(options.Scopes, base.WithUUID(UUID))
	entity, err := b.repository.FindOne(options)
	if err != nil {
		if errors.Is(err, base.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return entity, nil
}

func (b *baseCrud[T, R]) List(options base.ReadOptions, pagination base.Pagination) (*base.PaginatedResult[T], error) {
	return b.repository.Paginate(options, pagination)
}
