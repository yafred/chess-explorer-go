package server

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

// Start ... start a web server
func Start(port int) {

	http.HandleFunc("/explore", exploreHandler)
	http.HandleFunc("/", assetHandler) // we will use embed when go 1.16 is released

	log.Println("Server is listening on port " + strconv.Itoa(port))

	openbrowser("http://localhost:" + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
