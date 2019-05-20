package main

import "github.com/jinzhu/gorm"

type ServerJson struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}
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

type Endpoint struct {
	IpAddress    string `json:"ipAddress"`
	Grade        string `json:"grade"`
	Country      string `json:"country"`
	Organization string `json:"organization"`
}
type Domain struct {
	Endpoints []Endpoint `json:"endpoints"`
	Status    string     `json:"status"`
	Domain    string     `json:"domain"`
}

type Server struct {
	gorm.Model
	Address    string
	SslGrade   string
	Country    string
	Owner      string
	ResponseId int
	Response   Response `gorm:"foreignkey:ResponseId"` // use UserRefer as foreign ke
}
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

type History struct {
	Responses     []Response
	ResponsesJson []ResponseJson
}
