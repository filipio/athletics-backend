package models

type Pokemon struct {
	AppModel
	PokemonName string `json:"pokemon_name" validate:"required,oneof=D Pikachu Bulbasaur Charmander Squirtle,min=2,max=100"`
	Age         int    `json:"age" validate:"required,gte=0,lte=130" gorm:"type:int"`
	Email       string `json:"email" validate:"required,email" gorm:"not null;type:varchar(100)"`
	Attack      string `json:"attack" validate:"required,oneof=Thunderbolt Ember Vine_Whip Water_Gun,min=2,max=100" gorm:"not null"`
}
