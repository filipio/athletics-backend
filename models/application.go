package models

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

// defines the base model for all models in the application
type AppModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m AppModel) GetID() uint {
	return m.ID
}

// used for custom validation logic, which can't be defined in the struct tags
func (m AppModel) Validate(r *http.Request) error {
	return nil
}

func (m AppModel) BeforeCreateCtx(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (m AppModel) AfterCreateCtx(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (m AppModel) BeforeUpdateCtx(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (m AppModel) AfterUpdateCtx(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (m AppModel) BeforeDeleteCtx(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (m AppModel) AfterDeleteCtx(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (m AppModel) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db
}

func (m AppModel) GetQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return GetByIdQuery(db, r)
}

func (m AppModel) UpdateQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return baseUpdateQuery(m.GetQuery(db, r))
}

func (m AppModel) DeleteQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return m.GetQuery(db, r)
}

func GetByIdQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	id := utils.IntPathValue(r, "id")
	return db.Where("id = ?", id)
}

func baseUpdateQuery(db *gorm.DB) *gorm.DB {
	return db.Omit("id", "created_at")
}

func Paginate[T any](query *gorm.DB, r *http.Request, paginationOptions *utils.PaginationParams, mapResponseFunc func(T) any, queryInstance any) (*utils.PaginatedResponse, error) {
	var totalCount int64

	totalCountResult := query.Model(&queryInstance).Count(&totalCount)
	if totalCountResult.Error != nil {
		return nil, totalCountResult.Error
	}

	paginationParams := utils.BuildPaginationParams(r)
	if paginationOptions != nil {
		if paginationOptions.PerPage != 0 {
			paginationParams.PerPage = paginationOptions.PerPage
		}
		if paginationOptions.PageNo != 0 {
			paginationParams.PageNo = paginationOptions.PageNo
		}
		if paginationOptions.OrderBy != "" {
			paginationParams.OrderBy = paginationOptions.OrderBy
		}
		if paginationOptions.OrderDirection != "" {
			paginationParams.OrderDirection = paginationOptions.OrderDirection
		}
	}

	var records []T
	queryResult := PaginateQuery(query, paginationParams).Find(&records)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	var responseRecords []any = make([]any, len(records))
	for i, record := range records {
		if mapResponseFunc != nil {
			responseRecords[i] = mapResponseFunc(record)
		} else {
			responseRecords[i] = record
		}

	}

	return utils.BuildPaginatedResponse(responseRecords, totalCount, paginationParams), nil
}

func PaginateQuery(db *gorm.DB, paginationParams *utils.PaginationParams) *gorm.DB {
	return db.Offset((paginationParams.PageNo - 1) * paginationParams.PerPage).
		Limit(paginationParams.PerPage).Order(paginationParams.OrderBy + " " + paginationParams.OrderDirection)
}

func getByIds(db *gorm.DB, r *http.Request) *gorm.DB {
	queryParams := r.URL.Query()
	if queryParams.Has("ids") {
		ids := strings.Split(queryParams.Get("ids"), ",")
		return db.Where("id in (?)", ids)
	} else {
		return db
	}
}

// fetches the database connection from the request context
func Db(r *http.Request) *gorm.DB {
	return r.Context().Value(utils.DbContextKey).(*gorm.DB)
}
