package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	database_connection string = "postgresql://maxroach@localhost:26257/logs?sslmode=disable"
)

func getHistoryByUserAndDomain(response *Response, userSessionId string) {
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		return
	}
	query := fmt.Sprintf("SELECT H.created, H.response_id FROM History H INNER JOIN Response R ON H.response_id = R.response_id WHERE H.user_session_id='%s' AND R.domain='%s';", userSessionId, response.Domain)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			responseId int
			created    string
		)
		if err := rows.Scan(&created, &responseId); err != nil {
			log.Fatal(err)
		}
		test1 := "Fri May 16 22:07:16 -05 2019"
		start := time.Now()
		end, _ := time.Parse(time.UnixDate, test1)
		difference := start.Sub(end)
		fmt.Printf("time = %v\n", created)
		fmt.Printf("difference = %v\n", difference)
		fmt.Println(int(difference.Hours()))
		if int(difference.Hours()) >= 1 {
			// getResponseById(responseId)
			break
		}
	}
}

func getResponseById(responseId int) Response {
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		return Response{}
	}
	//---------PENDING CHANGE THE QUERY, VERIFY IF ITS OK OR CREATE NEW INNER JOINS
	var response Response
	query := fmt.Sprintf("SELECT response_id, response_id FROM History WHERE user_session_id='%s';", string(responseId))
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			response_id int
			created     string
		)
		if err := rows.Scan(&created, &response_id); err != nil {
			log.Fatal(err)
		}
		test1 := "Fri May 16 22:07:16 -05 2019"
		start := time.Now()
		end, _ := time.Parse(time.UnixDate, test1)
		difference := start.Sub(end)
		fmt.Printf("difference = %v\n", difference)
		fmt.Println(int(difference.Hours()))
		if int(difference.Hours()) >= 1 {

			break
		}
	}
	return response
}

func addHistory(responseId int, userSessionId string) {
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		return
	}
	t := time.Now()
	query := fmt.Sprintf("INSERT INTO History (user_session_id, response_id, created) VALUES ('%s',%d,'%s');", userSessionId, responseId, t.Format(time.UnixDate))
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}
func storeRecord(response Response) {
	var responseId = addResponse(response)
	if responseId != 0 {
		for i := 0; i < len(response.Servers); i++ {
			addServer(responseId, response.Servers[i])
		}
		addHistory(responseId, "test_id3")
	}
}
func addServer(responseId int, server Server) {
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		return
	}
	var (
		fields = "response_id, address"
		values = fmt.Sprintf("%d, '%s'", responseId, server.Address)
	)
	if server.SslGrade != "" {
		fields += ", ssl_grade"
		values += fmt.Sprintf(", '%s'", server.SslGrade)
	}
	if server.Country != "" {
		fields += ", country"
		values += fmt.Sprintf(", '%s'", server.Country)
	}
	if server.Owner != "" {
		fields += ", owner"
		values += fmt.Sprintf(", '%s'", server.Owner)
	}
	query := fmt.Sprintf("INSERT INTO Server (%s) VALUES (%s);", fields, values)
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func addResponse(response Response) int {
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		return 0
	}
	var (
		fields = "domain, servers_changed, is_down"
		values = fmt.Sprintf("'%s',%t,%t", response.Domain, response.ServersChanged, response.IsDown)
	)
	if response.SslGrade != "" {
		fields += ", ssl_grade"
		values += fmt.Sprintf(", '%s'", response.SslGrade)
	}
	if response.PreviousSslGrade != "" {
		fields += ", previous_ssl_grade"
		values += fmt.Sprintf(", '%s'", response.PreviousSslGrade)
	}
	if response.Logo != "" {
		fields += ", logo"
		values += fmt.Sprintf(", '%s'", response.Logo)
	}
	if response.Title != "" {
		fields += ", title"
		values += fmt.Sprintf(", '%s'", response.Title)
	}
	query := fmt.Sprintf("INSERT INTO Response (%s) VALUES (%s) RETURNING response_id;", fields, values)
	var response_id int
	err = db.QueryRow(query).Scan(&response_id)
	if err != nil {
		return 0
	}
	return response_id
}
