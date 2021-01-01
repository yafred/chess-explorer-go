package server

import (
	"log"
	"net/http"
	"strconv"
)

// Start ... start a web server on port 8080
func Start(port int) {
	fs := http.FileServer(http.Dir("../../web/static"))
	http.Handle("/", fs)

	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/explore", exploreHandler)
	http.HandleFunc("/find", findHandler)

	log.Println("Server is listening on port " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
