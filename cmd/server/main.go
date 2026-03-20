package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	handler "github.com/ArunGowrish/web-analyzer/internal/handler"
)

func main() {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}
	h := &handler.Handler{Tmpl: tmpl}

	// serve/load static folder
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/submit", h.SubmitHandler)
	fmt.Println("Server started and listening to 8080")
	http.ListenAndServe(":8080", nil)
}
