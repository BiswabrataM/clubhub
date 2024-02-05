package models

type Endpoints struct {
	Id           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	FranchiseId  uint   `json:"franchiseId" gorm:"foreignKey:FranchiseId"`
	IpAddress    string `json:"ipAddress"`
	ServerName   string `json:"serverName"`
	Creation     string `json:"creation"`
	Expiry       string `json:"expiry"`
	RegisteredTo string `json:"registeredTo"`
}
