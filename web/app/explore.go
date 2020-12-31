package server

import (
	"encoding/json"
	"net/http"
)

func exploreHandler(w http.ResponseWriter, r *http.Request) {
	type Sample struct {
		Title   string `json:"title"`
		Desc    string `json:"desc"`
		Content string `json:"content"`
	}
	json.NewEncoder(w).Encode(Sample{Title: "hello", Desc: "world", Content: "yeah"})
}
