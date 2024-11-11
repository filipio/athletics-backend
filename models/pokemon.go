package models

import (
	"net/http"

	"gorm.io/gorm"
)

type Pokemon struct {
	AppModel
	PokemonName string `json:"pokemon_name" validate:"required,oneof=D Pikachu Bulbasaur Charmander Squirtle,min=2,max=100"`
	Age         int    `json:"age" validate:"required,gte=0,lte=130" gorm:"type:int"`
	Email       string `json:"email" validate:"required,email" gorm:"not null;type:varchar(100)"`
	Attack      string `json:"attack" validate:"required,oneof=Thunderbolt Ember Vine_Whip Water_Gun,min=2,max=100" gorm:"not null"`
}

func (m Pokemon) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = getByIds(db, r)
	queryParams := r.URL.Query()

	if queryParams.Has("name") {
		db = db.Where("pokemon_name IN (?)", queryParams.Get("name"))
	}

	return db
}

func (m Pokemon) BuildResponse() any {
	return m
}
