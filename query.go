package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	database_connection string = "postgresql://maxroach@localhost:26257/logs?sslmode=disable"
)

/* func getHistoryByUserAndDomain(response *Response, userSessionId string) {
	db, err := sql.Open("postgres", database_connection)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
			SELECT H.created, H.response_id
			FROM History H
			INNER JOIN Response R ON H.response_id = R.response_id
			WHERE H.user_session_id='%s' AND R.domain='%s';`,
		userSessionId, response.Domain)
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
*/

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
	// db.Where("user_session_id = ?", "userSessionId").Find(&users)
	db.Where(Response{UserSessionId: userSessionId}).Find(&history.Responses)
	history.ResponsesJson = make([]ResponseJson, len(history.Responses))
	for i := 0; i < len(history.Responses); i++ {
		db.Where(Server{ResponseId: int(history.Responses[i].ID)}).Find(&history.Responses[i].Servers)
		history.ResponsesJson[i] = parseOrmToJson(history.Responses[i])
	}
	return history.ResponsesJson
}

func getChangesByDomain(userSessionId string, response ResponseJson) {
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
	db.Where(Response{UserSessionId: userSessionId, Domain: response.Domain}).Find(&history.Responses)
	history.ResponsesJson = make([]ResponseJson, len(history.Responses))
	for i := 0; i < len(history.Responses); i++ {
		createdDate := history.Responses[i].CreatedAt
		difference := start.Sub(createdDate)
		fmt.Printf("difference = %v\n", difference)
		fmt.Println(int(difference.Hours()))
		if int(difference.Hours()) >= 1 {
			db.Where(Server{ResponseId: int(history.Responses[i].ID)}).Find(&history.Responses[i].Servers)
			previousResponse = history.Responses[i]
			break
		}
	}
	if previousResponse.SslGrade != "" {
		response.PreviousSslGrade = previousResponse.SslGrade
	}
	if len(previousResponse.Servers) != 0 && len(previousResponse.Servers) != len(response.Servers) {
		response.ServersChanged = true
	} else {
		response.ServersChanged = hasServerChanged(response, previousResponse)
	}
	fmt.Print(previousResponse)
}
func hasServerChanged(response ResponseJson, previousResponse Response) bool {
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

func parseOrmToJson(response Response) ResponseJson {
	responseJson := ResponseJson{}
	responseJson.Title = response.Title
	responseJson.Logo = response.Logo
	responseJson.SslGrade = response.SslGrade
	responseJson.IsDown = response.IsDown
	responseJson.PreviousSslGrade = response.PreviousSslGrade
	responseJson.ServersChanged = response.ServersChanged
	responseJson.Domain = response.Domain
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
