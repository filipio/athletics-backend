package models

// type Address struct {
// 	Street string `json:"street" validate:"required"`
// 	City   string `json:"city" validate:"min=2,max=100"`
// }

// type Child struct {
// 	Height       *int     `json:"height" validate:"required"`
// 	Weight       *float32 `json:"weight" validate:"required"`
// 	ChildAddress *Address `json:"child_address" validate:"required"`
// }

type Pokemon struct {
	AppModel
	PokemonName string `json:"pokemon_name" validate:"required,oneof=D Pikachu Bulbasaur Charmander Squirtle,min=2,max=100"`
	Age         int    `json:"age" validate:"required,gte=0,lte=130" gorm:"type:int"`
	Email       string `json:"email" validate:"required,email" gorm:"not null;type:varchar(100)"`
	Attack      string `json:"attack" validate:"required,oneof=Thunderbolt Ember Vine_Whip Water_Gun,min=2,max=100" gorm:"not null"`
	// PokemonChild *Child  `json:"pokemon_child" validate:"required"`
}

type Other struct {
	AppModel
	PokemonName *string `json:"pokemon_name" validate:"required,oneof=D Pikachu Bulbasaur Charmander Squirtle,min=2,max=100"`
	Age         int     `json:"age" validate:"required,gte=0,lte=130"`
	Email       string  `json:"email" validate:"required,email"`
	// PokemonChild *Child  `json:"pokemon_child" validate:"required"`
}

type Human struct {
	AppModel
	Name   string  `json:"name" validate:"required,min=2,max=100"`
	Height float32 `json:"height" validate:"required,gte=0,lte=3"`
}
