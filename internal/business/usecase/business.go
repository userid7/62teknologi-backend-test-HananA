package usecase

import (
	"backend-test/internal/entity"
	"backend-test/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// BusinessUseCase -.
type BusinessUseCase struct {
	repo BusinessRepo
	l    logger.Interface
	v    *validator.Validate
}

// New -.
func NewBusinessUseCase(r BusinessRepo, l logger.Interface) *BusinessUseCase {
	return &BusinessUseCase{
		repo: r,
		l:    l,
		v:    validator.New(),
	}
}

func (bu *BusinessUseCase) Create(ctx context.Context, b entity.Business) error {
	b.ID = generateRandomToken(16)
	if err := bu.v.Struct(&b); err != nil {
		bu.l.Error(fmt.Errorf("usecase - Create - validate: %w", err))
		return err
	}
	if err := bu.repo.Create(ctx, b); err != nil {
		bu.l.Error(fmt.Errorf("usecase - Create - repo.Create: %w", err))
		return err
	}
	return nil
}

func (bu *BusinessUseCase) Read(ctx context.Context, id string) (entity.Business, error) {
	business, err := bu.repo.ReadById(ctx, id)
	if err != nil {
		bu.l.Error(fmt.Errorf("usecase - Read - repo.ReadById: %w", err))
		return business, err
	}
	return business, nil
}

func (bu *BusinessUseCase) Update(ctx context.Context, id string, b entity.Business) error {
	if err := bu.repo.UpdateById(ctx, id, b); err != nil {
		bu.l.Error(fmt.Errorf("usecase - Create - repo.UpdateById: %w", err))
		return err
	}
	bu.l.Info("usecase - Create - repo.UpdateById: 1 row updated")
	return nil
}

func (bu *BusinessUseCase) Delete(ctx context.Context, id string) error {
	if err := bu.repo.DeleteById(ctx, id); err != nil {
		bu.l.Error(fmt.Errorf("usecase - Create - repo.DeleteById: %w", err))
		return err
	}
	bu.l.Info("usecase - Create - repo.DeleteById: 1 row deleted")
	return nil
}

func (bu *BusinessUseCase) Search(ctx context.Context, sp entity.SearchBusinessParam) ([]entity.Business, error) {
	if sp.Limit >= 100 {
		sp.Limit = 100
	}

	var tm time.Time
	fmt.Println(sp.OpenNow)
	if sp.OpenNow {
		tm = time.Now()
	} else if sp.OpenAt > 0 {
		tm = time.Unix(int64(sp.OpenAt), 0)
	}

	businesses, err := bu.repo.Search(ctx, sp.Limit, sp.Offset, sp.Price, sp.Attributes, sp.Categories, tm)
	if err != nil {
		bu.l.Error(fmt.Errorf("usecase - Search - repo.Search: %w", err))
		return []entity.Business{}, err
	}
	return businesses, nil

}

func generateRandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func unixStrToTime(u string) (time.Time, error) {
	i, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}
