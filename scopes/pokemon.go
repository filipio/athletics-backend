package scopes

import (
	"fmt"
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

type ScopesFunc func(r *http.Request) []func(db *gorm.DB) *gorm.DB

func PokemonScopes(r *http.Request) []func(db *gorm.DB) *gorm.DB {
	result := []func(db *gorm.DB) *gorm.DB{}
	queryParams := r.URL.Query()

	if queryParams.Has("name") {
		result = append(result, func(db *gorm.DB) *gorm.DB {
			return db.Where("pokemon_name IN (?)", queryParams.Get("name"))
		})
	}

	userAny := r.Context().Value(utils.UserContextKey).(models.User)
	fmt.Println(userAny.ID, userAny.Email)

	// user := r.Context().Value(utils.UserContextKey).(*models.User)
	// if user != nil {
	// 	fmt.Println(user.ID, user.Email)
	// }

	return result
}
