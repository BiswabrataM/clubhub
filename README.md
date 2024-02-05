
### Architecture
The project follows a typical Go application structure with separate packages for services, models, and repositories. It uses the GORM library for database interactions and external libraries like Colly for web scraping.


### Workflow

                       +---------------------+
                       |      HTTP Request   |
                       +----------+----------+
                                  v
                       +---------------------+
                       |  server routes      |
                       +----------+----------+
                                  v
                       +---------------------+
                       |  server Handlers    |
                       +----------+----------+
                                  |
                 +----------------+-----------------+
                 v                                  v
        +----------------------+                  +----------------+
        |  Repository Layer    |                  |                |
        | Database Interaction |                  |  Service Layer |
        +--------+-------------+                  +--------+-------+
                 v                                      v
        +----------------+                         +----------------+
        | Model          |                         | Third-party    |
        | (Database      |                         | External, etc. |
        | Structures)    |                         |                |
        +--------+-------+                         +--------+-------+
 



## Database
        +--------------+     +--------------+    +--------------+    +---------------+    +--------------+
        |    Owner     |     |   Location   |    |  Franchise   |    |   Endpoints   |    |   Company    |
        +--------------+     +--------------+    +--------------+    +---------------+    +--------------+
        | Id (PK)      |     | Id (PK)      |    | Id (PK)      |    | Id (PK)       |    | Id (PK)      |
        | FirstName    |     | City         |    | Name         |    | FranchiseId   |    | Name         |
        | LastName     |     | Country      |    | URL          |    | IpAddress     |    | TaxNumber    |
        | Email        |     +--------------+    | Protocol     |    | ServerName    |    | Address      |
        | Phone        |                         | WebsiteAvail |    | Creation      |    | ZipCode      |
        | Address      |                         | IconURL      |    | Expiry        |    | LocationId   |
        | ZipCode      |                         | Address      |    | RegisteredTo  |    | OwnerId      |
        | LocationId   |                         | ZipCode      |    +---------------+    | Location     |
        | Location     |                         | LocationId   |                         | Owner        |
        +--------------+                         | CompanyId    |                         +--------------+
                                                 +--------------+

- Owner (One) - (Many) Company (Owner is associated with multiple companies)
- Location (One) - (Many) Owner (One location can have multiple owners)
- Location (One) - (Many) Company (One location can have multiple companies)
- Company (One) - (Many) Franchise (One company can have multiple franchises)
- Franchise (One) - (Many) Endpoints (One franchise can have multiple endpoints)



## APIs
- GET /hotelchain
- POST /hotelchain
- PATCH /hotelchain



### Run server
## docker
docker build -t <your-image-name> .
docker run -p 8080:8080 <your-image-name>

## terminal
go run server.go




### CURLs
## Check if server health
curl -X GET
'http://localhost:8080/apis/v1/health'
--header 'Accept: /' \


## GET existing data for hotelchain
curl -X GET
'http://localhost:8080/apis/v1/hotelchain/?managementCompanyName=My%20entreprise%20holding&franchiseName=Marriot'
--header 'Accept: /' \


## Add data for hotelchain
curl -X POST
'http://localhost:8080/apis/v1/hotelchain/'
--header 'Accept: /'
--header 'Content-Type: application/json'
--data-raw '{ "company": { "owner": { "first_name": "josh", "last_name": "porch", "contact": { "email": "josh@my-enterprise-holding.org", "phone": "+1 800 465 6574", "location": { "city": "Toronto", "country": "Canada", "address": "14 bulevar", "zip_code": "N6D 92A" } } }, "informacion": { "name": "My entreprise holding", "tax_number": "DD79654121", "location": { "city": "Toronto", "country": "Canada", "address": "78 Rober ST", "zip_code": "F9A 92O" } },

"franchises": [ { "name": "Park royal", "url": "www.park-royalhotels.com", "location": { "city": "Cancun", "Country": "Mexico", "Address": "Libertadores av 40 - 20", "zip_code": "45971" } }, { "name": "Marriot", "url": "www.marriott.com", "location": { "city": "Miami", "Country": "United States", "Address": "35 Tom st 18 bridge av", "zip_code": "115745" } } ] } } '

## Update existing data for hotelchain
curl -X PATCH
'http://localhost:8080/apis/v1/hotelchain/'
--header 'Accept: /'
--header 'Content-Type: application/json'
--data-raw '{ "company": { "id": 1, "name": "UpdatedCompanyName", "taxNumber": "UpdatedTaxNumber", "address": "UpdatedAddress", "zipCode": "UpdatedZipCode", "locationId": 2, "ownerId": 3, "owner": { "id": 1, "firstName": "UpdatedFirstName", "lastName": "UpdatedLastName", "contact": { "email": "updated@example.com", "phone": "9876543210", "location": { "city": "UpdatedCity", "country": "UpdatedCountry", "address": "UpdatedOwnerAddress", "zipCode": "UpdatedOwnerZipCode" } } }, "franchises": [ { "id": 1, "name": "UpdatedFranchiseName", "url": "UpdatedFranchiseURL", "location": { "city": "UpdatedFranchiseCity", "country": "UpdatedFranchiseCountry", "address": "UpdatedFranchiseAddress", "zipCode": "UpdatedFranchiseZipCode" } } ] } }'