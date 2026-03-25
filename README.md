
# web-analyzer

A web analyzer built in Go that processes a webpage URL to extract and evaluate metadata, structure, links, and forms for useful insights.

## Description
The objective of this project is to develop a web-based application in Go that performs automated analysis of a given webpage URL. The system processes user-provided URLs and extracts key structural and metadata information to provide meaningful insights into the composition and accessibility of the webpage.

The application begins by validating the input URL and then retrieves the webpage content for analysis. It identifies the HTML version used (e.g., HTML5, HTML 4.01, XHTML), extracts the page title, and analyzes the document structure by counting heading elements (H1–H6). In addition, it evaluates all hyperlinks on the page, classifying them into internal and external links, and detecting broken or inaccessible links through HTTP checks.

## Getting Started

### Dependencies
Go version 1.20 or higher. Check installed version using command - go version

### Installing
* Clone the repository - https://github.com/ArunGowrish/web-analyzer.git
* Install Dependencies - go mod tidy
* Production Build - go build cmd/server/main.go

### Executing program
go run cmd/server/main.go

### Access the app
http://localhost:8080

### Testing
* go test ./internal/handler -v
* go test ./internal/service -v
* go test ./utils -v
* go test ./internal/client -v

* To check complete test coverage - go test ./... -cover

## Project Structure

```
web-analyzer
├─ cmd
│  └─ server
│     └─ main.go
├─ configs
├─ deployments
├─ go.mod
├─ go.sum
├─ internal
│  ├─ client
│  │  ├─ http_client.go
│  │  └─ http_client_test.go
│  ├─ handler
│  │  ├─ submit_handler.go
│  │  └─ submit_handler_test.go
│  ├─ mocks
│  │  └─ http_client_mock.go
│  ├─ model
│  │  ├─ analysis_result.go
│  │  └─ link.go
│  └─ service
│     ├─ analyzer.go
│     └─ analyzer_test.go
├─ README.md
├─ static
│  └─ css
│     └─ style.css
├─ templates
│  └─ index.html
└─ utils
   ├─ validate_util.go
   └─ validate_util_test.go

```

## Architecture Overview

This project follows a clean layered architecture:

HTTP Request
     ↓
 Handler (controller layer)
     ↓
 Service (business logic)
     ↓
 Utils (helpers like parsing, validation)

Frontend

* Built with HTML + CSS
* Handles and submit the form to analyze the url

Backend

* Build with Golang
* Handles the analyzing part of the given url

## Example Form Submission

http://localhost:8080/submit

* Example Request:

POST /api/brand/logo
{
  "utl": "https://example.com"
}

* Example Response:

{
  "HTML Version": "HTML5"
}

## Authors
B.Gowrikumar (argowrish@gmail.com)
