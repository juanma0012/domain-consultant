package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

// Request the page by the given domain  and parse the Logo and Title from the page.
func parsePageHtml(response *ResponseJson, domainString string) {
	res, err := http.Get(fmt.Sprintf("http://%s", domainString))
	if err != nil {
		return
	} else {
		//create a new tokenizer over the res body
		tokenizer := html.NewTokenizer(res.Body)
		for {
			tokenType := tokenizer.Next()
			if tokenType == html.ErrorToken {
				err := tokenizer.Err()
				if err == io.EOF {
					//end of the file, break out of the loop
					break
				}
				log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			}

			//process the token according to the token type...
			if tokenType == html.StartTagToken {
				token := tokenizer.Token()
				if "title" == token.Data {
					//the next token should be the page title
					tokenType = tokenizer.Next()
					//just make sure it's actually a text token
					if tokenType == html.TextToken {
						response.Title = tokenizer.Token().Data
					}
				}
			}
			if tokenType == html.SelfClosingTagToken {
				token := tokenizer.Token()
				// <link href="image_url" type="image/x-icon"/>
				if "link" == token.Data {
					var (
						logoUrl   = ""
						logoFound = false
					)
					for i := 0; i < len(token.Attr); i++ {
						// get the link
						if token.Attr[i].Key == "href" {
							logoUrl = token.Attr[i].Val
						}
						// validate if the link has the attribute 'image/x-icon'
						if token.Attr[i].Key == "type" && token.Attr[i].Val == "image/x-icon" {
							logoFound = true
						}
					}
					// assign the correct link url
					if logoFound {
						response.Logo = logoUrl
					}
				}
			}
		}
	}
}
