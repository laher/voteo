package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/", fs)
	http.HandleFunc("/yt/data", ytDataHandler)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func ytDataHandler(r *http.Request, w http.ResponseWriter) {

}
