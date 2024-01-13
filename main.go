package main

import (
	"log"
	"net/http"
	"risk/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	classification map[int]string
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.MainHandler).Methods("GET")
	r.HandleFunc("/index1", handlers.Index1Handler).Methods("GET", "POST")
	r.HandleFunc("/index2", handlers.Index2Handler).Methods("GET", "POST")

	r.HandleFunc("/assets", handlers.AssetsHandler).Methods("POST")
	r.HandleFunc("/situations", handlers.SituationsHandler).Methods("POST")

	r.HandleFunc("/save_assets", handlers.SaveAssetsHandler).Methods("POST")
	r.HandleFunc("/save_situations", handlers.SaveSituationsHandler).Methods("POST")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
