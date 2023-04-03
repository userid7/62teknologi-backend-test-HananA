package repo

import (
	"backend-test/internal/entity"
	"backend-test/pkg/logger"
	"context"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BusinessRepo -.
type BusinessRepo struct {
	db *gorm.DB
	l  logger.Interface
}

// NewBusinessRepo -.
func NewBusinessRepo(db *gorm.DB, l logger.Interface) *BusinessRepo {
	return &BusinessRepo{
		db: db,
		l:  l,
	}
}

func (br *BusinessRepo) MigrateAndMockData() error {
	if err := br.db.AutoMigrate(&entity.Business{}, &entity.Categories{}); err != nil {
		br.l.Fatal(fmt.Errorf("app - Run - db.AutoMigrate: %w", err))
	}

	MockCategories := []entity.Categories{
		{Alias: "fnb", Name: "fnb"},
		{Alias: "office", Name: "office"},
		{Alias: "factory", Name: "factory"},
	}

	for _, m := range MockCategories {
		fmt.Println(m.Alias)
		if result := br.db.Where(&m).FirstOrCreate(&m); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (br *BusinessRepo) Create(ctx context.Context, b entity.Business) error {
	var cats []entity.Categories
	if result := br.db.Find(&cats, b.Categories); result.Error != nil {
		return result.Error
	}

	bWithoutCat := b
	bWithoutCat.Categories = nil

	if result := br.db.Create(&bWithoutCat); result.Error != nil {
		return result.Error
	}

	if err := br.db.Model(&bWithoutCat).Where("ID = ?", bWithoutCat.ID).Association("Categories").Append(&cats); err != nil {
		return err
	}

	return nil
}
func (br *BusinessRepo) ReadById(ctx context.Context, id string) (entity.Business, error) {
	var business entity.Business

	result := br.db.Model(&entity.Business{}).Preload("Categories").Where("id=?", id).First(&business)
	if result.Error != nil {
		return business, result.Error
	}

	return business, nil
}
func (br *BusinessRepo) UpdateById(ctx context.Context, id string, b entity.Business) error {
	business := &entity.Business{}
	var cats entity.Categories

	var catAliases []string

	for _, v := range b.Categories {
		catAliases = append(catAliases, v.Alias)
	}

	if result := br.db.Where("ID = ?", id).Find(&business); result.Error != nil {
		return result.Error
	}

	if result := br.db.Where("alias IN ?", catAliases).Find(&cats); result.Error != nil {
		return result.Error
	}

	fmt.Println(id)

	tx := br.db.Begin()

	if err := tx.Model(&business).Association("Categories").Replace(&cats); err != nil {
		tx.Rollback()
		fmt.Println("1")
		return err
	}

	// if err := tx.Model(&business).Association("Categories").Append(&cats); err != nil {
	// 	tx.Rollback()
	// 	fmt.Println("2")
	// 	return err
	// }

	b.Categories = nil

	result := tx.Model(&business).Where("ID = ?", id).Updates(&b)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("3")
		return result.Error
	}

	tx.Commit()
	return nil
}
func (br *BusinessRepo) DeleteById(ctx context.Context, id string) error {
	result := br.db.Where("id=?", id).Delete(&entity.Business{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (br *BusinessRepo) Search(ctx context.Context, limit uint, offset uint, price uint, attributes []string, categories []string, openAt time.Time) ([]entity.Business, error) {
	var businesses []entity.Business
	tx := br.db.Model(&entity.Business{}).Preload("Categories").WithContext(ctx)

	if price != 0 {
		tx = tx.Where("LENGTH(price) = ?", price)
	}

	if len(attributes) > 0 {
		tx = tx.Where(datatypes.JSONArrayQuery("attributes").Contains(attributes))
	}

	fmt.Println(len(categories))
	if len(categories) > 0 {
		// tx = tx.Where("uuid IN (SELECT business_uuid FROM business_categories WHERE alias IN ?)", categories)
		tx = tx.Where("uuid IN (?)", br.db.Table("business_categories").
			Joins("left join categories c on categories_id = c.id ").
			Select("business_uuid").
			Where("alias IN ?", categories),
		)
	}

	fmt.Println(openAt)
	if !openAt.IsZero() {
		layoutTime := "15:04:05"
		timeWithoutDate := openAt.Format(layoutTime)
		fmt.Println(timeWithoutDate)
		tx = tx.Where("open_time < ? AND close_time > ?", timeWithoutDate, timeWithoutDate)
	}

	res := tx.Limit(int(limit)).Offset(int(offset)).Find(&businesses)
	if res.Error != nil {
		return businesses, res.Error
	}

	return businesses, nil
}
