package models

type Location struct {
	Id      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	City    string `json:"city"`
	Country string `json:"country"`
}
