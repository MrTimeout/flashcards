package main

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// CategoriesPath retrieves all the categories inside the db.
	// /api/categories
	CategoriesPath = "/categories"
	// CategoryByIdPath is used by the handler GetCategoryById to
	// return the category by id.
	// /api/categories/:category-id
	CategoryByIdPath = ":" + CategoryParam

	// CategoryParam is the category id value
	CategoryParam = "category-id"
)

type category struct {
	XMLName     xml.Name  `gorm:"-" json:"-" xml:"Category"`
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:category_id;primaryKey" json:"-" xml:"-"`
	Name        string    `gorm:"column:name;not null;unique" json:"name" xml:"Name"`
	Description string    `gorm:"column:description;not null" json:"description" xml:"Description"`
	Word        []word
}

// OrderByColumnsAllowed will return the list of columns allowed to order by.
func (category) OrderByColumnsAllowed() map[string]any {
	return map[string]any{"id": struct{}{}, "name": struct{}{}}
}

// TableName is the name of the table in the database.
func (category) TableName() string {
	return "categories"
}

func getCategories(c *gin.Context) {
	categories, err := doGetCategories(c.Request.Context(), ParseRequest(c, category{}))
	if err != nil {
		ErrRes(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: Negotiate,
		Data:    categories,
	})
}

func getCategoryByName(c *gin.Context) {
	category, err := doGetCategories(c.Request.Context(), ParseRequest(c, category{Name: c.Param(CategoryParam)}))
	if err != nil {
		ErrRes(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: Negotiate,
		Data:    category,
	})
}

func addCategory(c *gin.Context) {
	var category category
	if err := c.Bind(&category); err != nil {
		ErrRes(c, err, http.StatusBadRequest)
		return
	}

	if _, err := doAddCategory(c.Request.Context(), &category); err != nil {
		ErrRes(c, err, http.StatusBadRequest)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: Negotiate,
		Data:    category,
	})
}

func delCategory(c *gin.Context) {
	var (
		rows int64
		err  error
	)
	if rows, err = doDelCategory(c.Request.Context(), category{Name: c.Param(CategoryParam)}); err != nil {
		ErrRes(c, err, http.StatusBadRequest)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: Negotiate,
		Data: WrapperResponse{
			Msg:  "category rows deleted " + strconv.Itoa(int(rows)),
			Code: http.StatusOK,
		},
	})
}

func doAddCategory(ctx context.Context, c *category) (rows int64, err error) {
	tx := getInstanceWithCtx(ctx).Create(c)
	return tx.RowsAffected, tx.Error
}

func doDelCategory(ctx context.Context, c category) (rows int64, err error) {
	tx := getInstanceWithCtx(ctx).Where(&c).Delete(c)
	return tx.RowsAffected, tx.Error
}

func doGetCategories(ctx context.Context, wrap WrapperRequest[category]) ([]category, error) {
	var result []category

	tx := wrap.ToScope(getInstanceWithCtx(ctx).Where(wrap.Body)).Find(&result)

	return result, tx.Error
}
