package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"

	"backend-test/pkg/logger"

	"backend-test/internal/business/usecase"

	"backend-test/internal/entity"
)

type businessRoutes struct {
	b usecase.Business
	l logger.Interface
	v *validator.Validate
}

func newBusinessRoutes(handler *gin.RouterGroup, t usecase.Business, l logger.Interface) {
	r := &businessRoutes{t, l, validator.New()}

	h := handler.Group("/business")
	{
		h.POST("/", r.addBusiness)
		h.PUT("/:id", r.updateBusiness)
		h.DELETE("/:id", r.deleteBusiness)

		h.GET("/:id", r.getBusiness)
		h.GET("/search", r.searchBusiness)

	}
}

type addBusinessRequest struct {
	Alias        string            `json:"alias" binding:"required"`
	Categories   []string          `json:"categories"`
	Coordinates  entity.Cordinates `json:"coordinates" `
	DisplayPhone string            `json:"display_phone"`
	ImageURL     string            `json:"image_url"`
	OpenTime     datatypes.Time    `json:"open_time" binding:"required"`
	CloseTime    datatypes.Time    `json:"close_time" binding:"required"`
	Location     entity.Location   `json:"location" `
	Name         string            `json:"name"`
	Phone        string            `json:"phone" `
	Price        string            `json:"price" binding:"required,gte=0,lte=4"`
	Attributes   []string          `json:"attributes"`
	URL          string            `json:"url"`
}

func (r *businessRoutes) addBusiness(c *gin.Context) {
	var req addBusinessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var cats []entity.Categories
	for _, cat := range req.Categories {
		cats = append(cats, entity.Categories{Alias: cat})
	}

	att := []string{}
	if req.Attributes != nil {
		att = req.Attributes
	}

	if err := r.b.Create(
		c,
		entity.Business{
			Alias:        req.Alias,
			Categories:   cats,
			Coordinates:  req.Coordinates,
			DisplayPhone: req.DisplayPhone,
			ImageURL:     req.ImageURL,
			OpenTime:     req.OpenTime,
			CloseTime:    req.CloseTime,
			Location:     req.Location,
			Name:         req.Name,
			Phone:        req.Phone,
			Price:        req.Price,
			Attributes:   datatypes.JSONType[[]string]{Data: att},
			Transactions: datatypes.JSONType[[]string]{Data: []string{}},
			URL:          req.URL},
	); err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "OK", "data": "business created"})
}

type updateBusinessRequest struct {
	Alias        string            `json:"alias" `
	Categories   []string          `json:"categories"`
	Coordinates  entity.Cordinates `json:"coordinates" `
	DisplayPhone string            `json:"display_phone"`
	ImageURL     string            `json:"image_url"`
	OpenTime     datatypes.Time    `json:"open_time" `
	CloseTime    datatypes.Time    `json:"close_time" `
	Location     entity.Location   `json:"location" `
	Name         string            `json:"name"`
	Phone        string            `json:"phone" `
	Price        string            `json:"price"`
	Attributes   []string          `json:"attributes"`
	URL          string            `json:"url"`
}

func (r *businessRoutes) updateBusiness(c *gin.Context) {
	var req updateBusinessRequest
	paramId := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  err.Error(),
		})
		return
	}

	var cats []entity.Categories
	for _, cat := range req.Categories {
		cats = append(cats, entity.Categories{Alias: cat})
	}

	att := []string{}
	if req.Attributes != nil {
		att = req.Attributes
	}

	if err := r.b.Update(c, paramId, entity.Business{
		Alias:        req.Alias,
		Categories:   cats,
		Coordinates:  req.Coordinates,
		DisplayPhone: req.DisplayPhone,
		ImageURL:     req.ImageURL,
		OpenTime:     req.OpenTime,
		CloseTime:    req.CloseTime,
		Location:     req.Location,
		Name:         req.Name,
		Phone:        req.Phone,
		Price:        req.Price,
		Attributes:   datatypes.JSONType[[]string]{Data: att},
		Transactions: datatypes.JSONType[[]string]{Data: []string{}},
		URL:          req.URL}); err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "OK", "data": "business updated"})
}

func (r *businessRoutes) deleteBusiness(c *gin.Context) {
	paramId := c.Param("id")

	if err := r.b.Delete(c, paramId); err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "OK", "data": "business deleted"})
}

func (r *businessRoutes) getBusiness(c *gin.Context) {
	paramId := c.Param("id")

	business, err := r.b.Read(c, paramId)
	if err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "OK", "data": business})
}

func (r *businessRoutes) searchBusiness(c *gin.Context) {
	var q entity.SearchBusinessQueryParam
	if err := c.ShouldBindQuery(&q); err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	sp := entity.SearchBusinessParam{
		Limit:      q.Limit,
		Offset:     q.Offset,
		Categories: q.Categories(),
		Attributes: q.Attributes(),
		Price:      q.Price,
		OpenAt:     q.OpenAt,
		OpenNow:    q.OpenNow,
	}
	businesses, err := r.b.Search(c, sp)
	if err != nil {
		r.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "OK", "data": businesses, "length": len(businesses)})
}
