package server

import (
	"log"
	"net/http"
	"strconv"
)

// Start ... start a web server
func Start(port int) {
	fs := http.FileServer(http.Dir("../../web/static"))
	http.Handle("/", fs)

	http.HandleFunc("/explore", exploreHandler)

	log.Println("Server is listening on port " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
