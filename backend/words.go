package main

import (
	"context"
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// WordsPath retrieves all the words inside the category
	// /api/categories/:category-id/words
	WordsPath = CategoryByIdPath + "/words"
	// WordsByIdPath retrieves all the information about one word inside a category
	WordsByIdPath = ":" + WordIdParam

	// WordIdParam is the name of the word
	WordIdParam = "word-id"
)

type word struct {
	XMLName    xml.Name  `gorm:"-" json:"-" xml:"Word"`
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:word_id;primaryKey" json:"-" xml:"-"`
	Term       string    `gorm:"column:term;not null;unique" json:"term" xml:"term"`
	Def        string    `gorm:"column:def;not null" json:"def" xml:"def"`
	CategoryID uuid.UUID `gorm:"column:category_id;not null" json:"-"`
}

func (word) OrderByColumnsAllowed() map[string]any {
	return map[string]any{"id": struct{}{}, "term": struct{}{}}
}

func (word) TableName() string {
	return "words"
}

func getWords(c *gin.Context) {
	categoryId, err := uuid.Parse(c.Param(CategoryParam))
	if err != nil {
		ErrRes(c, err, http.StatusBadRequest)
		return
	}

	words, err := doGetWords(c.Request.Context(), ParseRequest(c, word{CategoryID: categoryId}))
	if err != nil {
		ErrRes(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: Negotiate,
		Data:    words,
	})
}

func doGetWords(ctx context.Context, wrap WrapperRequest[word]) ([]word, error) {
	var result []word

	tx := wrap.ToScope(getInstanceWithCtx(ctx).Preload("Category").
		Joins("inner join categories on categories.category_id = words.category_id and categories.name like ?", wrap.Body.CategoryID)).
		Find(&result)

	return result, tx.Error
}
