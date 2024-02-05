package models

type Owner struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Address    string   `json:"address"`
	ZipCode    string   `json:"zipCode"`
	LocationId uint     `json:"locationId"`
	Location   Location `json:"location" gorm:"foreignKey:LocationId"`
}
