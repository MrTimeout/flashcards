package main

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// CategoriesPath retrieves all the categories inside the db.
	// /api/categories
	CategoriesPath = "/categories"
	// CategoryByNamePath is used by the handler GetCategoryByName to
	// return the category by name.
	// /api/categories/:category-name
	CategoryByNamePath = ":" + CategoryNameParam

	// CategoryNameParam is the category name value
	CategoryNameParam = "category-name"
)

type category struct {
	XMLName     xml.Name `gorm:"-" json:"-" xml:"Category"`
	ID          int      `gorm:"column:category_id;primaryKey" json:"-" xml:"-"`
	Name        string   `gorm:"column:name;not null;unique" json:"name" xml:"Name"`
	Description string   `gorm:"column:description;not null" json:"description" xml:"Description"`
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
	category, err := doGetCategories(c.Request.Context(), ParseRequest(c, category{Name: c.Param(CategoryNameParam)}))
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
	if rows, err = doDelCategory(c.Request.Context(), category{Name: c.Param(CategoryNameParam)}); err != nil {
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

func whereCategories(db *gorm.DB, fc category) *gorm.DB {
	if fc.ID != 0 {
		db = db.Where(fc.TableName()+".id = ?", fc.ID)
	}

	if fc.Description != "" {
		db = db.Where(fc.TableName()+".description like ?", fc.Description)
	}

	if fc.Name != "" {
		db = db.Where(fc.TableName()+".name = ?", fc.Name)
	}

	return db
}

func selectWhereCategories(db *gorm.DB, fc category, projection ...string) *gorm.DB {
	return db.Table(fc.TableName()).Select(projection).Where(&fc)
}
