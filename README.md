# Domain Consultant (Backend)
This repository provides  two endpoints to get the information of a domain. Also, the information is stored and available using the history endpoint
## API Docs

**Get the information in accordance to the domain provided**
> http://endpoint:3333/information/{domain}

Return an object with this information:
- Domain
- Servers
- Servers Changed
- Ssl Grade
- Previous Ssl Grade
- Logo
- Title
- Is Down
- Creation date

Response:
<pre>
 {
  "domain": "truora.com",
  "servers": [
    {
      "address": "34.193.69.252",
      "ssl_grade": "A",
      "country": "US",
      "owner": "Amazon Technologies Inc. (AT-88-Z)"
    },
    {
      "address": "34.193.204.92",
      "ssl_grade": "A",
      "country": "US",
      "owner": "Amazon Technologies Inc. (AT-88-Z)"
    }
  ],
  "servers_changed": false,
  "ssl_grade": "A",
  "previous_ssl_grade": "",
  "logo": "https://uploads-ssl.webflow.com/5b559a554de48fbcb01fd277/5b97f0ac932c3291fa40d053_icon32.png",
  "title": "Test",
  "is_down": false,
  "created_at": ""
}
</pre>


**Get the records of the previous consults**

> http://endpoint:3333/history

Returns a list with the records that the user requested previously, with the same json structure.
Response:
> [
{record},
{record}
]



#### Implemented technologies
- Go (Language)
- Chi (Routing)
- Gorm (Database ORM)
- Cockroach (Database SQL)

#### Running database
1. Download the database [here](https://www.cockroachlabs.com/docs/stable/install-cockroachdb-mac.html "here")
2. Start the node `$ cockroach start --insecure`
3. Setting the user and database `$ cockroach sql --insecure --execute="$(cat db_setting/inital_setting.sql)"`

#### Running the service
1. Install the dependencies (See Gopkg.tolm file) using  [dep](https://golang.github.io/dep "dep")
2. Build the package `$ go build`
3. Execute the package `$ ./domain-consultant)"`

#### Folder Structure
On this repository  will  find those files:

	├── Gopkg.lock
	├── Gopkg.toml
	├── README.md
	├── db_setting
	│   ├── commands
	│   └── initial_setting.sql
	├── domain-consultant
	├── domain-consultant.go
	├── html_tokenizer.go
	├── orm_functionality.go
	├── raw_data_request.go
	├── struct_types.go
	├── utility.go
    
### Database diagram

![](https://raw.githubusercontent.com/juanma0012/domain-consultant/master/db_setting/db_diagram.png)
