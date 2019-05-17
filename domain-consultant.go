package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	TextReplacer = regexp.MustCompile(`\n\[(.+?)\][\ ]+(.+?)`)

	sslGrades = map[string]int{"A+": 1, "A": 2, "B": 3, "C": 4, "D": 5, "E": 6, "F": 7}
)

type Server struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}
type Response struct {
	Domain           string   `json:"domain"`
	Servers          []Server `json:"servers"`
	ServersChanged   bool     `json:"servers_changed"`
	SslGrade         string   `json:"ssl_grade"`
	PreviousSslGrade string   `json:"previous_ssl_grade"`
	Logo             string   `json:"logo"`
	Title            string   `json:"title"`
	IsDown           bool     `json:"is_down"`
}

func main() {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go to /information/{domain}"))
	})

	// RESTy routes for "information" resource
	r.Route("/information", func(r chi.Router) {
		// r.Post("/", createRecord)                  // POST /information
		r.Get("/{domain}", getInformation) // GET /information/search
	})
	http.ListenAndServe(":3333", r)
}
func getInformation(w http.ResponseWriter, r *http.Request) {

	domainString := chi.URLParam(r, "domain") // from a route like /information/{domain}
	if domainString != "" {
		var domain Domain
		var attempt = 0
		for {
			domain = requestSsllabs(domainString)
			if domain.Status == "ERROR" || domain.Status == "READY" || attempt >= 3 {
				break
			} else {
				attempt += 1
				time.Sleep(3 * time.Second)
			}
		}
		for i := 0; i < len(domain.Endpoints); i++ {
			setWhoIsInformation(&domain.Endpoints[i])
		}
		var response Response
		if domain.Status == "ERROR" {
			response.IsDown = true
		} else {
			parseRawDataToResponse(&response, domain)
			parsePageHtml(&response, domainString)
			addResponse(response)
		}
		decodeData, _ := json.Marshal(response)
		w.Write([]byte(decodeData))
	} else {
		decodeData, _ := json.Marshal(Response{})
		w.Write([]byte(decodeData))
	}
}

/* func getPreviousData(userSessionId string, response *Response) {
	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:26257/bank?sslmode=disable")
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
		var id, balance int
		if err := rows.Scan(&id, &balance); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d %d\n", id, balance)
	}
}
*/
