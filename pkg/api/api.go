package api

import (
	"log"
	"net/http"
)

func InitAPI() {
	http.HandleFunc("/trasactions", AddHandler)
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddHandler(w, r)
		log.Println("addhandler")
	}
}
