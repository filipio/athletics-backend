package queries

import (
	"fmt"
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func GetPokemonsQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	queryFunctions := []func(db *gorm.DB) *gorm.DB{}
	queryParams := r.URL.Query()

	if queryParams.Has("name") {
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("pokemon_name IN (?)", queryParams.Get("name"))
		})
	}

	userAny := r.Context().Value(utils.UserContextKey).(models.User)
	fmt.Println(userAny.ID, userAny.Email)

	// user := r.Context().Value(utils.UserContextKey).(*models.User)
	// if user != nil {
	// 	fmt.Println(user.ID, user.Email)
	// }

	return db.Scopes(queryFunctions...)
}
