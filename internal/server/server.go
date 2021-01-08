package server

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/viper"
)

// Start ... start a web server
func Start() {

	fs := http.FileServer(http.Dir("./docs/v1"))
	http.Handle("/", fs)

	http.HandleFunc("/nextmove", nextMoveHandler)
	http.HandleFunc("/games", gamesHandler)

	port := viper.GetInt("server-port")
	if port == 0 {
		log.Fatal("server-port does not have a valid integer value")
	}
	log.Println("Server is listening on port " + strconv.Itoa(port))

	browser := viper.GetBool("start-browser")
	if browser {
		openbrowser("http://localhost:" + strconv.Itoa(port))
	}
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
