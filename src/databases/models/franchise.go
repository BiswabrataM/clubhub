package models

type Franchise struct {
	Id               uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string      `json:"name"`
	URL              string      `json:"url"`
	Protocol         string      `json:"protocol"`
	WebsiteAvailable bool        `json:"websiteAvailable"`
	IconURL          string      `json:"iconUrl"`
	Address          string      `json:"address"`
	ZipCode          string      `json:"zipCode"`
	LocationId       uint        `json:"locationId"`
	CompanyId        uint        `json:"companyId"`
	Endpoints        []Endpoints `json:"endpoints" gorm:"foreignKey:FranchiseId"`
	// Location   Location `json:"location" gorm:"foreignKey:LocationId"`
	// Company    Company  `json:"company" gorm:"foreignKey:CompanyId"`
}
