package repositories

import (
	"clubhub/src/databases"
	"clubhub/src/databases/models"
)

func GetCompanyInfo(companyName string, franchiseName string) (models.Company, error) {
	query := databases.DB

	if companyName != "" {
		query = query.Where(&models.Company{Name: companyName})
	}
	if franchiseName != "" {
		query = query.Joins("LEFT JOIN franchises ON companies.id = franchises.company_id").
			Where("franchises.name = ?", franchiseName)
	}

	query = query.Preload("Location").
		Preload("Owner").
		Preload("Owner.Location").
		Preload("Franchises").
		Preload("Franchises.Endpoints")

	var companyInfo models.Company
	if err := query.First(&companyInfo).Error; err != nil {
		return models.Company{}, err
	}

	return companyInfo, nil
}

type CompanyInfo struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	TaxNumber  string `json:"taxNumber"`
	Address    string `json:"address"`
	ZipCode    string `json:"zipCode"`
	LocationId uint   `json:"locationId"`
	OwnerId    uint   `json:"ownerId"`
}

func CreateCompany(companyInfo CompanyInfo) (uint, error) {
	var company = models.Company{
		Name:       companyInfo.Name,
		TaxNumber:  companyInfo.TaxNumber,
		Address:    companyInfo.Address,
		ZipCode:    companyInfo.ZipCode,
		LocationId: companyInfo.LocationId,
		OwnerId:    companyInfo.OwnerId,
	}
	result := databases.DB.Create(&company)
	if result.Error != nil {
		return 0, result.Error
	}
	return company.Id, nil
}
