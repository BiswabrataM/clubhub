package models

type Company struct {
	Id         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string      `json:"name"`
	TaxNumber  string      `json:"taxNumber"`
	Address    string      `json:"address"`
	ZipCode    string      `json:"zipCode"`
	LocationId uint        `json:"locationId"`
	OwnerId    uint        `json:"ownerId"`
	Location   Location    `json:"location" gorm:"foreignKey:LocationId"`
	Owner      Owner       `json:"owner" gorm:"foreignKey:OwnerId"`
	Franchises []Franchise `json:"franchises" gorm:"foreignKey:CompanyId"`
}
