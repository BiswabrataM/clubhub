package repositories

import (
	"clubhub/src/databases"
	"clubhub/src/databases/models"
)

// func GetOrAddLocation(locationInfo Location) (uint, error) {
func GetOrAddLocation(city string, country string) (uint, error) {
	var location models.Location
	result := databases.DB.Where(models.Location{City: city, Country: country}).FirstOrCreate(&location)
	if result.Error != nil {
		return 0, result.Error
	}
	return location.Id, nil
}
