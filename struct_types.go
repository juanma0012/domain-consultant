package main

import "github.com/jinzhu/gorm"

// Server model for the API server
type ServerJson struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

// Response model for the API response
type ResponseJson struct {
	Domain           string       `json:"domain"`
	Servers          []ServerJson `json:"servers"`
	ServersChanged   bool         `json:"servers_changed"`
	SslGrade         string       `json:"ssl_grade"`
	PreviousSslGrade string       `json:"previous_ssl_grade"`
	Logo             string       `json:"logo"`
	Title            string       `json:"title"`
	IsDown           bool         `json:"is_down"`
	CreatedAt        string       `json:"created_at"`
}

// Endpoint model to get the information that comes from the ssllabs endpoint and whois command
type Endpoint struct {
	IpAddress    string `json:"ipAddress"`
	Grade        string `json:"grade"`
	Country      string `json:"country"`
	Organization string `json:"organization"`
}

// Domain model to get the information that comes from the ssllabs endpoint
type Domain struct {
	Endpoints []Endpoint `json:"endpoints"`
	Status    string     `json:"status"`
	Domain    string     `json:"domain"`
}

// Server ORM to get the information that comes from the database
type Server struct {
	gorm.Model
	Address    string
	SslGrade   string
	Country    string
	Owner      string
	ResponseId int
	Response   Response `gorm:"foreignkey:ResponseId"`
}

// Response ORM to get the information that comes from the database
type Response struct {
	gorm.Model
	Domain           string
	UserSessionId    string
	ServersChanged   bool
	SslGrade         string
	PreviousSslGrade string
	Logo             string
	Title            string
	IsDown           bool
	Servers          []Server `gorm:"foreignkey:ResponseId;association_foreignkey:ID`
}

// History model to store temporal information while the processes are in progress.
type History struct {
	Responses     []Response
	ResponsesJson []ResponseJson
}
