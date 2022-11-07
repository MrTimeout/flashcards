package main

import (
	"encoding/xml"
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	errInvalidDirection = errors.New("invalid direction for the value supplied")
	errInvalidOrderBy   = errors.New("invalid order by for the value supplied")
	// ErrContentTypeNotAllowed is used when request contains an incorrect Content-Type.
	errContentTypeNotAllowed = errors.New("content type not allowed")

	// Negotiate is used to express which Accept and Content-Type MIME types are allowed.
	Negotiate = []string{gin.MIMEJSON, gin.MIMEXML}
)

const (
	limitQuery   = "limit"
	skipQuery    = "skip"
	orderByQuery = "order_by"

	limitDefault = 50
	limitMax     = 100
	limitMin     = 1

	skipDefault = 0
	skipMax     = math.MaxInt32
	skipMin     = 0
)

type Direction int

const (
	asc Direction = iota
	desc
)

// Type returns the string representation of the inner type
// of direction
func (d Direction) Type() string {
	return reflect.Int.String()
}

// Set is used to set the value to a direction type, checking
// if the input is inside the boundaries. If the value is not correct,
// by default, will be assigned asc
func (d *Direction) Set(src int) error {
	if !d.unmarshal(src) {
		return errInvalidDirection
	}
	return nil
}

func (d *Direction) unmarshal(src int) bool {
	switch src {
	case 0:
		*d = asc
	case 1:
		*d = desc
	default:
		return false
	}
	return true
}

func (d *Direction) unmarshalText(src string) bool {
	switch src {
	case "asc", "ASC":
		*d = asc
	case "desc", "DESC":
		*d = desc
	default:
		return false
	}
	return true
}

// String returns the string representation of the type direction
func (d Direction) String() string {
	var result string

	switch d {
	case asc:
		result = "ASC"
	case desc:
		result = "DESC"
	}

	return result
}

type WrapperRequest[T OrderByAllower] struct {
	Limit   int
	Skip    int
	OrderBy []orderBy
	Body    T
}

type orderBy struct {
	Field     string
	Direction Direction
}

func (o orderBy) String() string {
	return o.Field + " " + o.Direction.String()
}

type WrapperResponse struct {
	XMLName xml.Name `json:"-" xml:"Response"`
	Code    int      `json:"code" xml:"Code"`
	Msg     string   `json:"message" xml:"Message"`
}

func ParseRequest[T OrderByAllower](qParser QueryParser, body T) WrapperRequest[T] {
	return WrapperRequest[T]{
		Limit:   ParseNumber(qParser.Query(limitQuery), limitDefault, Boundaries(limitMin, limitMax)),
		Skip:    ParseNumber(qParser.Query(skipQuery), skipDefault, Boundaries(skipMin, skipMax)),
		OrderBy: parseArrOrderBy(qParser.QueryArray(orderByQuery)),
		Body:    body,
	}
}

func parseArrOrderBy(arr []string) []orderBy {
	var (
		result = make([]orderBy, len(arr))
		j      = 0
	)

	for i := range arr {
		if o, err := parseOrderBy(arr[i]); err == nil {
			result[j] = o
			j++
		}
	}

	return result[:j]
}

func parseOrderBy(input string) (orderBy, error) {
	var (
		splitted = strings.Split(input, " ")

		d Direction
	)

	if len(splitted) != 2 {
		return orderBy{}, errInvalidOrderBy
	}

	d.unmarshalText(splitted[1])

	return orderBy{
		Field:     splitted[0],
		Direction: d,
	}, nil
}

type Number interface {
	int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | int | uint
}

type Filter[T Number] func(input T) bool

func ParseNumber[T Number](input string, d T, filter ...Filter[T]) T {
	i, err := strconv.Atoi(input)
	if err != nil {
		return d
	}

	result := T(i)

	for _, f := range filter {
		if !f(result) {
			return d
		}
	}

	return result
}

func Boundaries[T Number](min, max T) Filter[T] {
	return func(input T) bool {
		return min <= input && input <= max
	}
}

type RepositoryBuilder interface {
	Build() *gorm.DB
}

type OrderByAllower interface {
	OrderByColumnsAllowed() map[string]any
}

func (w WrapperRequest[T]) ToScope(db *gorm.DB) *gorm.DB {
	return w.limit(w.skip(w.orderBy(db)))
}

func (w WrapperRequest[T]) orderBy(db *gorm.DB) *gorm.DB {
	var cols = w.Body.OrderByColumnsAllowed()

	for _, each := range w.OrderBy {
		if _, ok := cols[each.Field]; ok {
			db = db.Order(each.String())
		}
	}

	return db
}

func (w WrapperRequest[T]) limit(db *gorm.DB) *gorm.DB {
	return db.Limit(w.Limit)
}

func (w WrapperRequest[T]) skip(db *gorm.DB) *gorm.DB {
	return db.Offset(w.Skip)
}

type ParamParser interface {
	Param(string) string
}

type QueryParser interface {
	Query(string) string
	DefaultQuery(string, string) string
	QueryArray(string) []string
}

func parseErr(err error) string {
	if strings.Contains(err.Error(), "duplicate key value") {
		return "The entry that you are trying to insert, already exists"
	} else {
		return "Internal error"
	}
}

func ErrRes(g *gin.Context, err error, statusCode int) {
	g.Negotiate(statusCode, gin.Negotiate{
		Offered: Negotiate,
		Data: WrapperResponse{
			Code: statusCode,
			Msg:  parseErr(err),
		},
	})
}
