package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/likexian/whois-go"
)

var (
	CountryName      = ""
	OrganizationName = ""
	TextReplacer     = regexp.MustCompile(`\n\[(.+?)\][\ ]+(.+?)`)
)

type Endpoint struct {
	IpAddress    string `json:"ipAddress"`
	Grade        string `json:"grade"`
	Country      string `json:"country"`
	Organization string `json:"organization"`
}
type Domain struct {
	Endpoints []Endpoint `json:"endpoints"`
	Status    string     `json:"status"`
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
	var domain Domain
	domainString := chi.URLParam(r, "domain") // from a route like /information/{domain}
	if domainString != "" {
		var attempt = 0
		for {
			domain = requestSsllabs(domainString)
			if domain.Status == "ERROR" || domain.Status == "READY" || attempt >= 4 {
				break
			} else {
				attempt += 1
				time.Sleep(3 * time.Second)
				fmt.Println("timeout ")
			}
		}
		decodeData, _ := json.Marshal(domain)
		w.Write([]byte(decodeData))
	} else {
		decodeData, _ := json.Marshal(Domain{})
		w.Write([]byte(decodeData))
	}
}
func requestSsllabs(domainString string) Domain {
	var domain Domain
	response, err := http.Get(fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", domainString))
	if err != nil {
		return Domain{}
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &domain)
		for i := 0; i < len(domain.Endpoints); i++ {
			setWhoIsInformation(&domain.Endpoints[i])
		}
		return domain
	}
}

/* func getInformationByDomain(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain") // from a route like /information/{domain}
	if domain != "" {
		setWhoIsInformation(domain)
	}
	w.Write([]byte(OrganizationName))
} */

func setWhoIsInformation(endpoint *Endpoint) {
	result, err := whois.Whois(endpoint.IpAddress)
	if err != nil {
		fmt.Println(result)
		return
	}
	whoisText := strings.Replace(result, "\r", "", -1)
	whoisText = TextReplacer.ReplaceAllString(whoisText, "\n$1: $2")

	whoisLines := strings.Split(whoisText, "\n")
	for i := 0; i < len(whoisLines); i++ {
		line := strings.TrimSpace(whoisLines[i])
		if len(line) < 5 || !strings.Contains(line, ":") {
			continue
		}

		fChar := line[:1]
		if fChar == ">" || fChar == "%" || fChar == "*" {
			continue
		}

		if line[len(line)-1:] == ":" {
			i += 1
			for ; i < len(whoisLines); i++ {
				thisLine := strings.TrimSpace(whoisLines[i])
				if strings.Contains(thisLine, ":") {
					break
				}
				line += thisLine + ","
			}
			line = strings.Trim(line, ",")
			i -= 1
		}

		lines := strings.SplitN(line, ":", 2)
		name := strings.ToLower(strings.TrimSpace(lines[0]))
		value := strings.TrimSpace(lines[1])

		if value == "" {
			continue
		} else if name == "country" {
			endpoint.Country = value
		} else if name == "organization" {
			endpoint.Organization = value
		}
		if endpoint.Country != "" && endpoint.Organization != "" {
			break
		}
	}
}
