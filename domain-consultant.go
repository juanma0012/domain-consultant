package main

import (
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
		r.Get("/{domain}", getInformationByDomain) // GET /information/search
	})
	// RESTy routes for "information" resource
	r.Route("/domain", func(r chi.Router) {
		// r.Post("/", createRecord)                  // POST /information
		r.Get("/", getInformation) // GET /information/search
	})
	http.ListenAndServe(":3333", r)

	/* type Foo struct {
		status, protocol string
	} */

	/* foo := Foo{}
	json.Unmarshal([]byte(string(data)), &foo)
	fmt.Printf(foo.status)

	foo2 := new(Foo)
	// defer response.Body.Close()
	err := json.NewDecoder(response.Body).Decode(foo2)
	if err != nil {
		return
	}
	fmt.Println(foo2) */
	// fmt.Println(string(data))
	/* if data.status == "READY" {
		for i := 0; i < len(data.endpoints); i++ {
			setWhoIsInformation(data.endpoints[i].ipAddress)
			fmt.Println("OrganizationName=", OrganizationName)
			fmt.Println("CountryName=", CountryName)
		}
	} */
}
func getInformation(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=google.com")
	if err != nil {
		w.Write([]byte("The HTTP request failed with error"))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		w.Write([]byte(string(data)))
	}
}
func getInformationByDomain(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain") // from a route like /information/{domain}
	if domain != "" {
		setWhoIsInformation(domain)
	}
	w.Write([]byte(OrganizationName))
}

func setWhoIsInformation(ip string) {
	CountryName = ""
	OrganizationName = ""
	result, err := whois.Whois(ip)
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
		name := strings.TrimSpace(lines[0])
		value := strings.TrimSpace(lines[1])

		if value == "" {
			continue
		} else if name == "Country" {
			CountryName = value
		} else if name == "Organization" {
			OrganizationName = value
		}
		if CountryName != "" && OrganizationName != "" {
			break
		}
	}
}
