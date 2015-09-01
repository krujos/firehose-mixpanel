package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudfoundry-community/cfenv"
)

func root(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Hello World")
}

func setupHTTP(port int) {
	http.HandleFunc("/", root)
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if nil != err {
			log.Fatal("ListenAndServe:", err)
		}
	}()
}

func getToken(appEnv cfenv.Services) string {
	return "this is a dummy token"
}

func main() {
	appEnv, _ := cfenv.Current()
	setupHTTP(appEnv.Port)
	token := getToken(appEnv.Services)

	log.Printf("Using token" + token)
}
