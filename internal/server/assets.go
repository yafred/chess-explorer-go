package server

import (
	"net/http"

	"github.com/yafred/chess-explorer/assets"
)

/*
When assets must be embedded:
go get github.com/pyros2097/go-embed
go-embed -input public/ -output assets/main.go
*/
func assetHandler(res http.ResponseWriter, req *http.Request) {
	data, hash, contentType, err := assets.Asset("www/", req.URL.Path)
	if err != nil {
		data, hash, contentType, err = assets.Asset("www", "/index.html")
		if err != nil {
			data = []byte(err.Error())
		}
	}
	res.Header().Set("Content-Encoding", "gzip")
	res.Header().Set("Content-Type", contentType)
	res.Header().Add("Cache-Control", "public, max-age=31536000")
	res.Header().Add("ETag", hash)
	if req.Header.Get("If-None-Match") == hash {
		res.WriteHeader(http.StatusNotModified)
	} else {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write(data)
		if err != nil {
			panic(err)
		}
	}
}
