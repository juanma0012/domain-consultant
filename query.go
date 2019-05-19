package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func storeResponse(response ResponseJson, userSessionId string) {
	db, err := gorm.Open("postgres", database_connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Automatically create the "accounts" table based on the Account model.
	db.AutoMigrate(&Server{})
	db.AutoMigrate(&Response{})

	// Insert two rows into the "response" table.
	responseOrm := parseJsonToOrm(response, userSessionId)
	db.Create(&responseOrm)
}

func getHistoryByUser(userSessionId string) []ResponseJson {
	db, err := gorm.Open("postgres", database_connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Automatically create the "accounts" table based on the Account model.
	db.AutoMigrate(&Server{})
	db.AutoMigrate(&Response{})

	var history History
	// Get all matched records
	db.Where(Response{UserSessionId: userSessionId}).Find(&history.Responses)
	history.ResponsesJson = make([]ResponseJson, len(history.Responses))
	for i := 0; i < len(history.Responses); i++ {
		db.Where(Server{ResponseId: int(history.Responses[i].ID)}).Find(&history.Responses[i].Servers)
		history.ResponsesJson[i] = parseOrmToJson(history.Responses[i])
	}
	return history.ResponsesJson
}

func getChangesByDomain(userSessionId string, response *ResponseJson) {
	db, err := gorm.Open("postgres", database_connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&Server{})
	db.AutoMigrate(&Response{})

	var (
		history          History
		previousResponse Response
		start            = time.Now()
	)
	db.Order("created_at desc").Where(Response{UserSessionId: userSessionId, Domain: response.Domain}).Find(&history.Responses)
	history.ResponsesJson = make([]ResponseJson, len(history.Responses))
	for i := 0; i < len(history.Responses); i++ {
		createdDate := history.Responses[i].CreatedAt
		difference := start.Sub(createdDate)
		if int(difference.Hours()) >= 1 {
			db.Where(Server{ResponseId: int(history.Responses[i].ID)}).Find(&history.Responses[i].Servers)
			previousResponse = history.Responses[i]
			break
		}
	}
	response.PreviousSslGrade = previousResponse.SslGrade
	response.ServersChanged = hasServerChanged(response, previousResponse)
}
func hasServerChanged(response *ResponseJson, previousResponse Response) bool {
	if len(previousResponse.Servers) != 0 && len(previousResponse.Servers) != len(response.Servers) {
		return true
	}
	for i := 0; i < len(previousResponse.Servers); i++ {
		for j := 0; j < len(response.Servers); j++ {
			if previousResponse.Servers[i].Address == response.Servers[j].Address &&
				previousResponse.Servers[i].SslGrade != "" &&
				previousResponse.Servers[i].SslGrade != response.Servers[j].SslGrade {
				return true
			}
		}
	}
	return false
}
