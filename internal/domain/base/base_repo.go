package base

import (
	"errors"

	"gorm.io/gorm"
)

var ErrDuplicatedKey = errors.New("duplicate key error")
var ErrRecordNotFound = errors.New("not found")

type Scope = func(*gorm.DB) *gorm.DB

type ReadOptions struct {
	Scopes  []Scope
	Preload []string
}

type WriteOptions struct {
	Tx *gorm.DB
}

func (o ReadOptions) Apply(db *gorm.DB) *gorm.DB {
	for _, scope := range o.Scopes {
		db = db.Scopes(scope)
	}
	for _, preload := range o.Preload {
		db = db.Preload(preload)
	}
	return db
}

type Pagination struct {
	Page     int
	PageSize int
	OrderBy  string
}

func (p Pagination) Apply(db *gorm.DB) *gorm.DB {
	page := p.Page
	pageSize := p.PageSize

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	db = db.Limit(pageSize).Offset(offset)

	if p.OrderBy != "" {
		db = db.Order(p.OrderBy)
	}

	return db
}

type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	TotalCount int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
}

type BaseRepository[T any] interface {
	FindOne(ReadOptions) (*T, error)
	FindMany(ReadOptions) ([]T, error)
	Paginate(ReadOptions, Pagination) (*PaginatedResult[T], error)
	Create(WriteOptions, *T) error
	Update(WriteOptions, *T) error
	Delete(WriteOptions, *T) error
}

type BaseRepo[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{db: db}
}

func (r *BaseRepo[T]) FindOne(opts ReadOptions) (*T, error) {
	var model T
	err := opts.Apply(r.db).First(&model).Error
	return &model, wrapError(err)
}

func (r *BaseRepo[T]) FindMany(opts ReadOptions) ([]T, error) {
	var list []T
	err := opts.Apply(r.db).Find(&list).Error
	return list, wrapError(err)
}

func (r *BaseRepo[T]) Paginate(opts ReadOptions, pagination Pagination) (*PaginatedResult[T], error) {
	var list []T
	db := opts.Apply(r.db)

	var count int64
	if err := db.Model(new(T)).Count(&count).Error; err != nil {
		return nil, err
	}

	db = pagination.Apply(db)
	if err := db.Find(&list).Error; err != nil {
		return nil, wrapError(err)
	}

	return &PaginatedResult[T]{
		Data:       list,
		TotalCount: count,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
	}, nil
}

func (r *BaseRepo[T]) Create(opts WriteOptions, model *T) error {
	action := func(tx *gorm.DB) error {
		return wrapError(tx.Create(model).Error)
	}
	if opts.Tx != nil {
		return action(opts.Tx)
	}
	return action(r.db)
}

func (r *BaseRepo[T]) Update(opts WriteOptions, model *T) error {
	action := func(tx *gorm.DB) error {
		return wrapError(tx.Save(model).Error)
	}
	if opts.Tx != nil {
		return action(opts.Tx)
	}
	return action(r.db)
}

func (r *BaseRepo[T]) Delete(opts WriteOptions, model *T) error {
	action := func(tx *gorm.DB) error {
		return wrapError(tx.Delete(model).Error)
	}
	if opts.Tx != nil {
		return action(opts.Tx)
	}
	return action(r.db)
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrRecordNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrDuplicatedKey
	}

	return err
}
