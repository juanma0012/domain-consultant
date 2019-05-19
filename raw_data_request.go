package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/likexian/whois-go"
)

func requestSsllabs(domainString string) Domain {
	var domain Domain
	response, err := http.Get(fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", domainString))
	if err != nil {
		return Domain{Status: "ERROR"}
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &domain)
		domain.Domain = domainString
		return domain
	}
}

func setWhoIsInformation(endpoint *Endpoint) {
	result, err := whois.Whois(endpoint.IpAddress)
	if err != nil {
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
