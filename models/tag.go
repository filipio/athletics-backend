package model

type Tag struct {
	Id   int `gorm:"type:int;primary_key"`
	Name string
}
