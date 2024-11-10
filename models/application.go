package models

import (
	"context"
	"net/http"
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

// defines the base query function for all models in the application
type BuildQueryFunc func(db *gorm.DB, r *http.Request) *gorm.DB

func DefaultQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db
}

func GetByIdQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	id := utils.IntPathValue(r, "id")
	return db.Where("id = ?", id)
}

func DefaultUpdateQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return baseUpdateQuery(GetByIdQuery(db, r))
}

func baseUpdateQuery(db *gorm.DB) *gorm.DB {
	return db.Omit("id", "created_at")
}

func Paginate(db *gorm.DB, pageNo int, perPage int, orderBy string, orderDirection string) *gorm.DB {
	return db.Offset((pageNo - 1) * perPage).Limit(perPage).Order(orderBy + " " + orderDirection)
}

func doNothing(db *gorm.DB) *gorm.DB {
	return db
}

func getByIds(r *http.Request) func(db *gorm.DB) *gorm.DB {
	queryParams := r.URL.Query()
	if queryParams.Has("ids") {
		ids := queryParams.Get("ids")
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("id in (?)", ids)
		}
	}

	return doNothing
}

// defines the base response function for all models in the application
type BuildResponseFunc[T any, V any] func(T) V

type DefaultResponse struct {
	Data any `json:"data"`
}

func BuildDefaultResponse[T any](model T) T {
	return model
}

// fetches the database connection from the request context
func Db(r *http.Request) *gorm.DB {
	return r.Context().Value(utils.DbContextKey).(*gorm.DB)
}
