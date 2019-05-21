package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Applying the middleware setting
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//CORS setting
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "User-Session-Id"},
		ExposedHeaders:   []string{"User-Session-Id"},
		AllowCredentials: true,
	})
	r.Use(cors.Handler)
	r.Use(middleware.Timeout(60 * time.Second))

	// Default URI
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go to /information/{domain}"))
	})

	// REST route for "information" resource
	r.Route("/information", func(r chi.Router) {
		r.Get("/{domain}", getInformation) // GET /information/{domain}
	})
	// REST route for "history" resource
	r.Route("/history", func(r chi.Router) {
		r.Get("/", getHistory) // GET /history
	})
	http.ListenAndServe(":3333", r)
}
func getInformation(w http.ResponseWriter, r *http.Request) {
	// Getting the domain variable from the route /information/{domain}
	domainString := chi.URLParam(r, "domain")
	// Getting the use session id variable from the header request
	userSessionId := r.Header.Get("User-Session-Id")
	if domainString != "" {
		var domain Domain
		for {
			// Calling the ssllabs endpoint
			domain = requestSsllabs(domainString)
			// The endpoint return the status ERROR, READY, IN-PROGRESS, DNS
			if domain.Status == "ERROR" || domain.Status == "READY" {
				break
			} else {
				// if the status is either IN-PROGRESS or DNS, we wait
				// 6 seconds to request again the information with the status READY
				time.Sleep(6 * time.Second)
			}
		}
		// Calling for each server the command whois ip
		for i := 0; i < len(domain.Endpoints); i++ {
			setWhoIsInformation(&domain.Endpoints[i])
		}
		var response ResponseJson
		if domain.Status == "ERROR" {
			response.IsDown = true
		} else {
			// Collect all the information in one response object
			parseRawDataToResponse(&response, domain)
			parsePageHtml(&response, domainString)
			getChangesByDomain(userSessionId, &response)
			storeResponse(response, userSessionId)
		}
		decodeData, _ := json.Marshal(response)
		w.Write([]byte(decodeData))
	} else {
		decodeData, _ := json.Marshal(ResponseJson{})
		w.Write([]byte(decodeData))
	}
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	userSessionId := r.Header.Get("User-Session-Id")
	var history []ResponseJson
	// Getting the history of the  requests done previously
	history = getHistoryByUser(userSessionId)
	decodeData, _ := json.Marshal(history)
	w.Write([]byte(decodeData))
}
