package models

import (
	"net/http"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Athlete struct {
	AppModel
	FirstName   *string         `json:"first_name"`
	LastName    *string         `json:"last_name"`
	Birthday    *datatypes.Date `json:"birthday"`
	Country     *string         `json:"country"`
	Gender      string          `json:"gender" gorm:"not null"`
	Disciplines []Discipline    `json:"disciplines" gorm:"many2many:athletes_disciplines;constraint:OnDelete:CASCADE"`
}

func (m Athlete) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = db.Preload("Disciplines")
	db = getByIds(db, r)

	queryParams := r.URL.Query()

	if queryParams.Has("search") {
		searchTerm := "%" + strings.ToLower(queryParams.Get("search")) + "%"
		db = db.Where("first_name || ' ' || last_name LIKE ? OR last_name || ' ' || first_name LIKE ?", searchTerm, searchTerm)
	}

	if queryParams.Has("discipline_ids") {
		disciplineIds := strings.Split(queryParams.Get("discipline_ids"), ",")

		db = db.Joins("JOIN athletes_disciplines ON athletes_disciplines.athlete_id = athletes.id").
			Where("athletes_disciplines.discipline_id IN (?)", disciplineIds)
	}

	if queryParams.Has("country") {
		country := strings.ToUpper(queryParams.Get("country"))
		db = db.Where("country = ?", country)
	}

	if queryParams.Has("gender") {
		gender := strings.ToLower(queryParams.Get("gender"))
		db = db.Where("gender = ?", gender)
	}

	return db
}

func (m Athlete) GetQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = db.Preload("Disciplines")
	db = GetByIdQuery(db, r)
	return db
}

type AthleteResponse struct {
	ID          uint     `json:"id"`
	FirstName   *string  `json:"first_name"`
	LastName    *string  `json:"last_name"`
	Gender      string   `json:"gender"`
	Country     *string  `json:"country"`
	Disciplines []string `json:"disciplines"`
}

func (m Athlete) BuildResponse() any {
	disciplines := make([]string, len(m.Disciplines))
	for i, discipline := range m.Disciplines {
		disciplines[i] = discipline.Name
	}
	athleteResponse := AthleteResponse{
		ID:          m.ID,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Gender:      m.Gender,
		Country:     m.Country,
		Disciplines: disciplines,
	}

	return athleteResponse
}
