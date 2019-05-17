package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	database_connection string = "postgresql://maxroach@localhost:26257/logs?sslmode=disable"
)

/* func getResponseId(userSessionId string) string{
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	// Print out the balances.
	query, _ := fmt.Sprintf(("SELECT response_id FROM History WHERE user_session_id='%s'ORDER BY time DESC LIMIT 1 ", "some_user_session_id")
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("Initial balances:")
	for rows.Next() {
		var responseId string
		if err := rows.Scan(&id); err != nil {
			return ""
		}
		return responseId
	}
}
func addHistoryLog(userSessionId string) string{
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	// Insert two rows into the "accounts" table.
    if _, err := db.Exec(
        "INSERT INTO History (id, balance) VALUES (1, 1000), (2, 250)"); err != nil {
        log.Fatal(err)
    }
}

func addServer(responseId string, server Server) {
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	query, _ := fmt.Sprintf(("INSERT INTO Server (address, ssl_grade, country, owner, respond_id) VALUES (1, 1000), (2, 250) ", "some_user_session_id")
	rows, err := db.Query(query)
    if _, err := db.Exec(
        "INSERT INTO Server () VALUES (1, 1000), (2, 250)"); err != nil {
        log.Fatal(err)
    }
} */

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
	fmt.Printf(query)
	var response_id int
	err = db.QueryRow(query).Scan(&response_id)
	if err != nil {
		return 0
	}
	return response_id
}
