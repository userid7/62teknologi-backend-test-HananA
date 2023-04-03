// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"time"

	"backend-test/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Business -.
	Business interface {
		Create(context.Context, entity.Business) error
		Read(context.Context, string) (entity.Business, error)
		Search(context.Context, entity.SearchBusinessParam) ([]entity.Business, error)
		Update(context.Context, string, entity.Business) error
		Delete(context.Context, string) error
	}

	// BusinessRepo -.
	BusinessRepo interface {
		Create(context.Context, entity.Business) error
		ReadById(context.Context, string) (entity.Business, error)
		UpdateById(context.Context, string, entity.Business) error
		DeleteById(context.Context, string) error
		Search(ctx context.Context, limit uint, offset uint, price uint, attributes []string, categories []string, openAt time.Time) ([]entity.Business, error)
	}
)
