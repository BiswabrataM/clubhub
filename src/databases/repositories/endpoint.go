package repositories

import (
	"clubhub/src/databases"
	"clubhub/src/databases/models"
)

type EndpointInfo struct {
	Id           uint
	FranchiseId  uint
	IpAddress    string
	ServerName   string
	Creation     string
	Expiry       string
	RegisteredTo string
}

// func GetOrAddLocation(locationInfo Location) (uint, error) {
func AddEndpoint(endpointInfo EndpointInfo) (uint, error) {
	var endpoint = models.Endpoints{
		FranchiseId:  endpointInfo.FranchiseId,
		IpAddress:    endpointInfo.IpAddress,
		ServerName:   endpointInfo.ServerName,
		Creation:     endpointInfo.Creation,
		Expiry:       endpointInfo.Expiry,
		RegisteredTo: endpointInfo.RegisteredTo,
	}
	result := databases.DB.Create(&endpoint)
	if result.Error != nil {
		return 0, result.Error
	}
	return endpoint.Id, nil
}
