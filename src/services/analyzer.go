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

func AnalyzeFranchiseWebsite(franchiseURL string) ([]DomainInfo, FranchiseInfo, error) {
	var domainInfoList []DomainInfo
	var franchiseInfo FranchiseInfo

	iconUri, err := ExtractLogo(franchiseURL)
	franchiseInfo.IconUri = iconUri
	franchiseInfo.WebsiteAvailable = true
	if err != nil {
		franchiseInfo.IconUri = ""
		franchiseInfo.WebsiteAvailable = false
	}

	sslDetail, err := GetSSLInfo(franchiseURL)
	if err != nil {
		return domainInfoList, franchiseInfo, err
	}
	franchiseInfo.Protocol = sslDetail.Protocol

	for _, endpoint := range sslDetail.Endpoints {

		endpointInfo, err := GetDomainInfo(endpoint.IPAddress)
		if err != nil {
			return domainInfoList, franchiseInfo, err
		}

		domainInfo := DomainInfo{
			IpAddress:    endpoint.IPAddress,
			ServerName:   endpoint.ServerName,
			Creation:     endpointInfo.Creation,
			Expiry:       endpointInfo.Expiry,
			RegisteredTo: endpointInfo.RegisteredTo,
		}

		domainInfoList = append(domainInfoList, domainInfo)
	}

	return domainInfoList, franchiseInfo, nil

}

func GetSSLInfo(franchiseURL string) (SSLLabsResponse, error) {
	var result SSLLabsResponse

	domainRegex := regexp.MustCompile(`^(?:https?://)?(?:www\.)?([^/]+)`)
	matches := domainRegex.FindStringSubmatch(franchiseURL)
	if len(matches) < 2 {
		return result, fmt.Errorf("unable to extract domain from the URI")
	}
	domain := matches[1]

	var retries int
	for retries = 0; retries < 3; retries++ {
		log.Printf("api: fetching SSL details, try:%v", retries)
		reqUri := fmt.Sprintf("%s/%s?%s=%s", configs.SSLlabBaseUri, configs.SSLlabAnalyzeUri, configs.SSLlabQueryHost, domain)
		response, err := http.Get(reqUri)
		if err != nil {
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			err = json.NewDecoder(response.Body).Decode(&result)
			if err != nil {
				return result, err
			}
			return result, nil
		}

	}

	if retries == 3 {
		return result, fmt.Errorf("maximum retries reached")
	}
	log.Printf("api: fetched SSL details, host: %v", result.Host)
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
	if err != nil {
		log.Panicf("util: failed to fetch domain/ipaddress related info")
		return endpointInfo, err
	}

	creationDateRegex := regexp.MustCompile(`RegDate:\s+(\d{4}-\d{2}-\d{2})`)
	creationDateMatches := creationDateRegex.FindStringSubmatch(whoisResponse)
	if len(creationDateMatches) == 2 {
		creationDateStr := creationDateMatches[1]
		creationDate, err := time.Parse("2006-01-02", creationDateStr)
		if err != nil {
			return endpointInfo, err
		}
		endpointInfo.Creation = creationDate.Format("2006-01-02")
	}

	expiryDateRegex := regexp.MustCompile(`Updated:\s+(\d{4}-\d{2}-\d{2})`)
	expiryDateMatches := expiryDateRegex.FindStringSubmatch(whoisResponse)
	if len(expiryDateMatches) == 2 {
		expiryDateStr := expiryDateMatches[1]
		expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
		if err != nil {
			return endpointInfo, err
		}
		endpointInfo.Expiry = expiryDate.Format("2006-01-02")
	}

	registrantInfoRegex := regexp.MustCompile(`OrgName:\s+(.+)`)
	registrantInfoMatches := registrantInfoRegex.FindStringSubmatch(whoisResponse)
	if len(registrantInfoMatches) == 2 {
		endpointInfo.RegisteredTo = strings.TrimSpace(registrantInfoMatches[1])
	}

	return endpointInfo, nil
}

func ExtractLogo(franchiseURL string) (string, error) {
	c := colly.NewCollector()
	var websiteUri string = "https://" + franchiseURL
	var iconURL string

	log.Printf("util: extracting icon details from franchise website, uri: %v", franchiseURL)

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
