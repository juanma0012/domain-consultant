package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

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
		r.Get("/{domain}", getInformation) // GET /information/{domain}
	})
	// RESTy routes for "history" resource
	r.Route("/history", func(r chi.Router) {
		// r.Post("/", createRecord)                  // POST /information
		r.Get("/", getHistory) // GET /history
	})
	http.ListenAndServe(":3333", r)
}
func getInformation(w http.ResponseWriter, r *http.Request) {

	domainString := chi.URLParam(r, "domain") // from a route like /information/{domain}
	if domainString != "" {
		var domain Domain
		for {
			domain = requestSsllabs(domainString)
			if domain.Status == "ERROR" || domain.Status == "READY" {
				break
			} else {
				time.Sleep(6 * time.Second)
			}
		}
		for i := 0; i < len(domain.Endpoints); i++ {
			setWhoIsInformation(&domain.Endpoints[i])
		}
		var response ResponseJson
		if domain.Status == "ERROR" {
			response.IsDown = true
		} else {
			parseRawDataToResponse(&response, domain)
			parsePageHtml(&response, domainString)
			//getHistoryByUserAndDomain(&response, "test_id3")
			//storeRecord(response)
			//test()
			getChangesByDomain("session_number_5", response)
			storeResponse(response, "session_number_5")
		}
		decodeData, _ := json.Marshal(response)
		w.Write([]byte(decodeData))
	} else {
		decodeData, _ := json.Marshal(ResponseJson{})
		w.Write([]byte(decodeData))
	}
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	var history []ResponseJson
	history = getHistoryByUser("session_number_5")
	decodeData, _ := json.Marshal(history)
	w.Write([]byte(decodeData))
}
