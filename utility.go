package main

import (
	"regexp"
	"time"
)

// Defining the global variables to use in different files.
var (
	TextReplacer               = regexp.MustCompile(`\n\[(.+?)\][\ ]+(.+?)`)
	sslGrades                  = map[string]int{"A+": 1, "A": 2, "B": 3, "C": 4, "D": 5, "E": 6, "F": 7}
	database_connection string = "postgresql://maxroach@localhost:26257/logs?sslmode=disable"
)

// Parsing a Json object to a ORM object
func parseJsonToOrm(response ResponseJson, userSessionId string) Response {
	responseOrm := Response{}
	responseOrm.Title = response.Title
	responseOrm.Logo = response.Logo
	responseOrm.SslGrade = response.SslGrade
	responseOrm.IsDown = response.IsDown
	responseOrm.PreviousSslGrade = response.PreviousSslGrade
	responseOrm.ServersChanged = response.ServersChanged
	responseOrm.Domain = response.Domain
	responseOrm.UserSessionId = userSessionId
	responseOrm.Servers = make([]Server, len(response.Servers))
	for i := 0; i < len(response.Servers); i++ {
		responseOrm.Servers[i] = Server{
			Address:  response.Servers[i].Address,
			SslGrade: response.Servers[i].SslGrade,
			Country:  response.Servers[i].Country,
			Owner:    response.Servers[i].Owner}
	}
	return responseOrm
}

// Parsing a ORM object to a Json object
func parseOrmToJson(response Response) ResponseJson {
	responseJson := ResponseJson{}
	responseJson.Title = response.Title
	responseJson.Logo = response.Logo
	responseJson.SslGrade = response.SslGrade
	responseJson.IsDown = response.IsDown
	responseJson.PreviousSslGrade = response.PreviousSslGrade
	responseJson.ServersChanged = response.ServersChanged
	responseJson.Domain = response.Domain
	responseJson.CreatedAt = response.CreatedAt.Format(time.RFC850)
	responseJson.Servers = make([]ServerJson, len(response.Servers))
	for i := 0; i < len(response.Servers); i++ {
		responseJson.Servers[i] = ServerJson{
			Address:  response.Servers[i].Address,
			SslGrade: response.Servers[i].SslGrade,
			Country:  response.Servers[i].Country,
			Owner:    response.Servers[i].Owner}
	}
	return responseJson
}

// Parsing the data getting from the saucelabs and whois to Json object
func parseRawDataToResponse(response *ResponseJson, domain Domain) {
	response.Servers = make([]ServerJson, len(domain.Endpoints))
	for i := 0; i < len(domain.Endpoints); i++ {
		response.Servers[i] = ServerJson{
			Address:  domain.Endpoints[i].IpAddress,
			SslGrade: domain.Endpoints[i].Grade,
			Country:  domain.Endpoints[i].Country,
			Owner:    domain.Endpoints[i].Organization}
		if response.Servers[i].SslGrade != "" && sslGrades[response.SslGrade] < sslGrades[response.Servers[i].SslGrade] {
			response.SslGrade = response.Servers[i].SslGrade
		}
	}
	response.Domain = domain.Domain
}
