package repositories

import (
	"clubhub/src/databases"
	"clubhub/src/databases/models"
)

type FranchiseInfo struct {
	Id               uint   `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Protocol         string `json:"protocol"`
	WebsiteAvailable bool   `json:"websiteAvailable"`
	IconURL          string `json:"iconUrl"`
	Address          string `json:"address"`
	ZipCode          string `json:"zipCode"`
	LocationId       uint   `json:"locationId"`
	CompanyId        uint   `json:"companyId"`
}

func CreateFranchise(franchiseInfo FranchiseInfo) (uint, error) {
	var franchise = models.Franchise{
		Name:             franchiseInfo.Name,
		URL:              franchiseInfo.URL,
		Protocol:         franchiseInfo.Protocol,
		WebsiteAvailable: franchiseInfo.WebsiteAvailable,
		IconURL:          franchiseInfo.IconURL,
		Address:          franchiseInfo.Address,
		ZipCode:          franchiseInfo.ZipCode,
		LocationId:       franchiseInfo.LocationId,
		CompanyId:        franchiseInfo.CompanyId,
	}
	result := databases.DB.Create(&franchise)
	if result.Error != nil {
		return 0, result.Error
	}
	return franchise.Id, nil
}
