// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Business struct {
	UUID         uint                         `json:"-" gorm:"primaryKey"`
	Alias        string                       `json:"alias" gorm:"unique"`
	Categories   []Categories                 `json:"categories" gorm:"many2many:business_categories" validate:"dive"`
	Coordinates  Cordinates                   `json:"coordinates" gorm:"embedded;embeddedPrefix:cord_"`
	DisplayPhone string                       `json:"display_phone"`
	Distance     float64                      `json:"distance" gorm:"-"`
	ID           string                       `json:"id" gorm:"unique;not null"`
	ImageURL     string                       `json:"image_url"`
	OpenTime     datatypes.Time               `json:"-"  gorm:"coloumn:open_time"`
	CloseTime    datatypes.Time               `json:"-"  gorm:"coloumn:close_time"`
	IsOpen       bool                         `json:"is_open"  gorm:"-"`
	Location     Location                     `json:"location" gorm:"embedded;embeddedPrefix:loc_"`
	Name         string                       `json:"name"`
	Phone        string                       `json:"phone"`
	Price        string                       `json:"price" validate:"gte=0,lte=4"`
	Rating       int                          `json:"rating" validate:"gte=0,lte=5"`
	ReviewCount  int                          `json:"review_count"`
	Transactions datatypes.JSONType[[]string] `json:"transactions"`
	Attributes   datatypes.JSONType[[]string] `json:"attributes"`
	URL          string                       `json:"url"`
	CreatedAt    time.Time                    `json:"-"`
	UpdatedAt    time.Time                    `json:"-"`
	DeletedAt    gorm.DeletedAt               `json:"-" gorm:"index"`
}

type BusinessCategories struct {
	gorm.Model
	BusinessID   uint
	CategoriesID uint
}

type Categories struct {
	ID    uint   `json:"-" gorm:"primarykey"`
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type Cordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Location struct {
	Address1       string                     `json:"address1"`
	Address2       string                     `json:"address2"`
	Address3       string                     `json:"address3"`
	City           string                     `json:"city"`
	Country        string                     `json:"country"`
	DisplayAddress datatypes.JSONType[string] `json:"display_address"`
	State          string                     `json:"state"`
	ZipCode        string                     `json:"zip_code"`
}

type SearchBusinessResponse struct {
	Businesses []Business `json:"businesses"`
	Total      int        `json:"total"`
}

type SearchBusinessQueryParam struct {
	Limit         uint    `form:"limit"`
	Offset        uint    `form:"offset"`
	CategoriesStr string  `form:"categories"`
	AttributesStr string  `form:"attributes"`
	Radius        float64 `form:"radius"`
	Price         uint    `form:"price"`
	OpenAt        uint    `form:"open_at"`
	OpenNow       bool    `form:"open_now"`
}

func (q *SearchBusinessQueryParam) Categories() []string {
	if q.CategoriesStr != "" {
		return strings.Split(q.CategoriesStr, ",")
	}
	return []string{}
}
func (q *SearchBusinessQueryParam) Attributes() []string {
	if q.AttributesStr != "" {
		return strings.Split(q.AttributesStr, ",")
	}
	return []string{}

}

type SearchBusinessParam struct {
	Limit      uint
	Offset     uint
	Categories []string
	Attributes []string
	Radius     float64
	Price      uint
	OpenAt     uint
	OpenNow    bool
}
