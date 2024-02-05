package handlers

import (
	"clubhub/src/databases"
	"clubhub/src/databases/models"
	"clubhub/src/databases/repositories"
	"clubhub/src/services"
	"clubhub/src/utils/dtos"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type CreateRequestPayload struct {
	Company struct {
		Owner struct {
			FirstName string  `json:"first_name"`
			LastName  string  `json:"last_name"`
			Contact   Contact `json:"contact"`
		} `json:"owner"`
		Information struct {
			Name      string   `json:"name"`
			TaxNumber string   `json:"tax_number"`
			Location  Location `json:"location"`
		} `json:"information"`
		Franchises []struct {
			Name     string   `json:"name"`
			URL      string   `json:"url"`
			Location Location `json:"location"`
		} `json:"franchises"`
	} `json:"company"`
}

type UpdateRequestPayload struct {
	Company struct {
		Id        uint     `json:"id"`
		Name      string   `json:"name"`
		TaxNumber string   `json:"taxNumber"`
		Address   string   `json:"address"`
		ZipCode   string   `json:"zipCode"`
		Location  Location `json:"location"`
		OwnerId   uint     `json:"ownerId"`
		Owner     struct {
			Id        uint   `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Contact   Contact
		} `json:"owner"`
		Franchises []struct {
			Id       uint     `json:"id"`
			Name     string   `json:"name"`
			URL      string   `json:"url"`
			Location Location `json:"location"`
		} `json:"franchises"`
	} `json:"company"`
}

type Contact struct {
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Address string `json:"address"`
	ZipCode string `json:"zip_code"`
}

func GetAllHotels(c echo.Context) error {
	franchiseName := c.QueryParam("franchiseName")
	managementCompanyName := c.QueryParam("managementCompanyName")

	companyInfo, err := repositories.GetCompanyInfo(managementCompanyName, franchiseName)
	if err != nil {
		log.Error("Failed to fetch company: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, companyInfo)
}

func CreateHotel(c echo.Context) error {
	var createRequest CreateRequestPayload
	if err := c.Bind(&createRequest); err != nil {
		log.Error("Invalid request body:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	tx := databases.DB.Begin()

	ownerLocationId, err := repositories.GetOrAddLocation(createRequest.Company.Owner.Contact.Location.City, createRequest.Company.Owner.Contact.Location.Country)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to create or retrieve location:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create or retrieve location"})
	}

	ownerInfo := models.Owner{
		FirstName:  createRequest.Company.Owner.FirstName,
		LastName:   createRequest.Company.Owner.LastName,
		Email:      createRequest.Company.Owner.Contact.Email,
		Phone:      createRequest.Company.Owner.Contact.Phone,
		Address:    createRequest.Company.Owner.Contact.Location.Address,
		ZipCode:    createRequest.Company.Owner.Contact.Location.ZipCode,
		LocationId: ownerLocationId,
	}

	ownerId, err := repositories.CreateOwner(ownerInfo)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to create owner:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create owner"})
	}

	companyLocationId, err := repositories.GetOrAddLocation(createRequest.Company.Information.Location.City, createRequest.Company.Information.Location.Country)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to create or retrieve company location:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create or retrieve company location"})
	}

	companyInfo := repositories.CompanyInfo{
		Name:       createRequest.Company.Information.Name,
		TaxNumber:  createRequest.Company.Information.TaxNumber,
		Address:    createRequest.Company.Information.Location.Address,
		ZipCode:    createRequest.Company.Information.Location.ZipCode,
		LocationId: companyLocationId,
		OwnerId:    ownerId,
	}

	fmt.Println("---- >> ----", companyInfo)

	companyId, err := repositories.CreateCompany(companyInfo)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to create company:", err)
		var response = dtos.ErrorResponse{Message: "Failed to create company", Error: err}
		return c.JSON(http.StatusInternalServerError, response)
	}

	for _, franchisePayload := range createRequest.Company.Franchises {
		franchiseLocationId, err := repositories.GetOrAddLocation(franchisePayload.Location.City, franchisePayload.Location.Country)
		if err != nil {
			tx.Rollback()
			log.Error("Failed to create or retrieve franchise location:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create or retrieve franchise location"})
		}

		endpointDetails, fInfo, err := services.AnalyzeFranchiseWebsite(franchisePayload.URL)
		if err != nil {
			tx.Rollback()
			log.Error("Failed to create or retrieve franchise location:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create or retrieve franchise location"})
		}

		franchiseInfo := repositories.FranchiseInfo{
			Name:             franchisePayload.Name,
			URL:              franchisePayload.URL,
			IconURL:          fInfo.IconUri,
			WebsiteAvailable: fInfo.WebsiteAvailable,
			Protocol:         fInfo.Protocol,
			Address:          franchisePayload.Location.Address,
			ZipCode:          franchisePayload.Location.ZipCode,
			LocationId:       franchiseLocationId,
			CompanyId:        companyId,
		}

		franchiseId, err := repositories.CreateFranchise(franchiseInfo)
		if err != nil {
			tx.Rollback()
			log.Error("Failed to create franchise:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create franchise"})
		}

		var endpointInfo = repositories.EndpointInfo{
			FranchiseId:  franchiseId,
			IpAddress:    endpointDetails.IpAddress,
			ServerName:   endpointDetails.ServerName,
			Creation:     endpointDetails.Creation,
			Expiry:       endpointDetails.Expiry,
			RegisteredTo: endpointDetails.RegisteredTo,
		}
		_, err = repositories.AddEndpoint(endpointInfo)
		if err != nil {
			tx.Rollback()
			log.Error("Failed to create franchise:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create endpoint"})
		}
	}

	tx.Commit()

	log.Info("Hotel created successfully")
	return c.JSON(http.StatusCreated, map[string]string{"message": "Hotel created successfully"})
}

func UpdateDetails(c echo.Context) error {
	var updateRequest UpdateRequestPayload

	if err := c.Bind(&updateRequest); err != nil {
		log.Error("Invalid request body:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	tx := databases.DB.Begin()

	companyLocationId, err := repositories.GetOrAddLocation(updateRequest.Company.Location.City, updateRequest.Company.Location.Country)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to update franchise:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update franchise"})
	}

	// Update Company
	companyToUpdate := repositories.CompanyInfo{
		Name:       updateRequest.Company.Name,
		TaxNumber:  updateRequest.Company.TaxNumber,
		Address:    updateRequest.Company.Address,
		ZipCode:    updateRequest.Company.ZipCode,
		LocationId: companyLocationId,
	}

	if err := tx.Model(&models.Company{}).Where("id = ?", updateRequest.Company.Id).Updates(&companyToUpdate).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to update company:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update company"})
	}

	// Update Owner
	ownerLocationId, err := repositories.GetOrAddLocation(updateRequest.Company.Owner.Contact.Location.City, updateRequest.Company.Owner.Contact.Location.Country)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to update franchise:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update franchise"})
	}

	ownerToUpdate := models.Owner{
		FirstName:  updateRequest.Company.Owner.FirstName,
		LastName:   updateRequest.Company.Owner.LastName,
		Email:      updateRequest.Company.Owner.Contact.Email,
		Phone:      updateRequest.Company.Owner.Contact.Phone,
		Address:    updateRequest.Company.Owner.Contact.Location.Address,
		ZipCode:    updateRequest.Company.Owner.Contact.Location.ZipCode,
		LocationId: ownerLocationId,
	}

	if err := tx.Model(&models.Owner{}).Where("id = ?", updateRequest.Company.Owner.Id).Updates(&ownerToUpdate).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to update owner:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update owner"})
	}

	// Update Franchises
	for _, franchise := range updateRequest.Company.Franchises {
		franchiseLocationId, err := repositories.GetOrAddLocation(franchise.Location.City, franchise.Location.Country)
		if err != nil {
			tx.Rollback()
			log.Error("Failed to update franchise:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update franchise"})
		}

		franchiseToUpdate := models.Franchise{
			Name:       franchise.Name,
			URL:        franchise.URL,
			Address:    franchise.Location.Address,
			ZipCode:    franchise.Location.ZipCode,
			LocationId: franchiseLocationId,
		}

		if err := tx.Model(&models.Franchise{}).Where("id = ?", franchise.Id).Updates(&franchiseToUpdate).Error; err != nil {
			tx.Rollback()
			log.Error("Failed to update franchise:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update franchise"})
		}

	}

	tx.Commit()

	log.Info("Details updated successfully")
	return c.JSON(http.StatusOK, map[string]string{"message": "Details updated successfully"})
}
