package services

import (
	"clubhub/configs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/likexian/whois"
)

type SSLLabsResponse struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	Protocol  string     `json:"protocol"`
	Status    string     `json:"status"`
	StartTime int64      `json:"startTime"`
	TestTime  int64      `json:"testTime"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	IPAddress         string `json:"ipAddress"`
	ServerName        string `json:"serverName"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
	Progress          int    `json:"progress"`
	Duration          int    `json:"duration"`
	Delegation        int    `json:"delegation"`
}

type DomainInfo struct {
	Id           uint
	IpAddress    string
	ServerName   string
	Creation     string
	Expiry       string
	RegisteredTo string
}

type FranchiseInfo struct {
	IconUri          string
	Protocol         string
	WebsiteAvailable bool
}

func AnalyzeFranchiseWebsite(franchiseURL string) (DomainInfo, FranchiseInfo, error) {
	var domainInfo DomainInfo
	var franchiseInfo FranchiseInfo

	sslDetail, err := GetSSLInfo(franchiseURL)
	if err != nil {
		return domainInfo, franchiseInfo, err
	}

	iconUri, err := ExtractLogo(franchiseURL)
	if err != nil {
		return domainInfo, franchiseInfo, err
	}

	endpointInfo, err := GetDomainInfo(sslDetail.Endpoints[0].IPAddress)
	if err != nil {
		return domainInfo, franchiseInfo, err
	}

	domainInfo = DomainInfo{
		IpAddress:    sslDetail.Endpoints[0].IPAddress,
		ServerName:   sslDetail.Endpoints[0].ServerName,
		Creation:     endpointInfo.Creation,
		Expiry:       endpointInfo.Expiry,
		RegisteredTo: endpointInfo.RegisteredTo,
	}

	franchiseInfo = FranchiseInfo{
		IconUri:          iconUri,
		Protocol:         sslDetail.Protocol,
		WebsiteAvailable: true,
	}

	return domainInfo, franchiseInfo, nil
}

func GetSSLInfo(franchiseURL string) (SSLLabsResponse, error) {
	var result SSLLabsResponse

	reqUri := fmt.Sprintf("%s/%s?%s=%s", configs.SSLlabBaseUri, configs.SSLlabAnalyzeUri, configs.SSLlabQueryHost, franchiseURL)
	response, err := http.Get(reqUri)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

type EndpointInfo struct {
	Creation     string
	Expiry       string
	RegisteredTo string
}

func GetDomainInfo(ipAddress string) (EndpointInfo, error) {
	var endpointInfo EndpointInfo

	whoisResponse, err := whois.Whois(ipAddress)
	if err == nil {
		// fmt.Println(whoisResponse)
	}

	// Define regular expressions for capturing creation date, expiry date, and registrant information
	creationDateRegex := regexp.MustCompile(`RegDate:\s+(\d{4}-\d{2}-\d{2})`)
	expiryDateRegex := regexp.MustCompile(`Updated:\s+(\d{4}-\d{2}-\d{2})`)
	registrantInfoRegex := regexp.MustCompile(`OrgName:\s+(.+)`)

	// Find matches in the WHOIS response
	creationDateMatches := creationDateRegex.FindStringSubmatch(whoisResponse)
	expiryDateMatches := expiryDateRegex.FindStringSubmatch(whoisResponse)
	registrantInfoMatches := registrantInfoRegex.FindStringSubmatch(whoisResponse)

	// Extract matched information
	if len(creationDateMatches) == 2 {
		creationDateStr := creationDateMatches[1]
		creationDate, err := time.Parse("2006-01-02", creationDateStr)
		if err != nil {
			return endpointInfo, err
		}
		endpointInfo.Creation = creationDate.Format("2006-01-02")
	}

	if len(expiryDateMatches) == 2 {
		expiryDateStr := expiryDateMatches[1]
		expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
		if err != nil {
			return endpointInfo, err
		}
		endpointInfo.Expiry = expiryDate.Format("2006-01-02")
	}

	if len(registrantInfoMatches) == 2 {
		endpointInfo.RegisteredTo = strings.TrimSpace(registrantInfoMatches[1])
	}

	return endpointInfo, nil
}

func ExtractLogo(franchiseURL string) (string, error) {
	c := colly.NewCollector()
	var websiteUri string = "https://" + franchiseURL
	var iconURL string

	c.OnHTML("link[rel='icon'], link[rel='apple-touch-icon'], link[rel='shortcut icon']", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "http") {
			iconURL = href
		} else if strings.HasPrefix(href, "/") {
			iconURL = websiteUri + href
		}
	})

	err := c.Visit(websiteUri)
	if err != nil {
		return iconURL, err
	}
	log.Printf("icon extracted from franchiseURL: %v", iconURL)

	return iconURL, nil
}
