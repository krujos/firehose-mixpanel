package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/cloudfoundry-community/cfenv"
	"github.com/krujos/uaaclientcredentials"
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
	uaa, err := GetUserProvidedServiceByName("uaa", appEnv)
	if nil != err {
		log.Fatal("Could not find uaa service, did you bind it?")
		os.Exit(1)
	}

	uaaURL, err := url.Parse(uaa.Credentials["uri"].(string))
	if nil != err {
		log.Fatal("Could not parse uaa URI")
		log.Fatal(err)
		os.Exit(1)
	}

	uaaConnection, err := uaaclientcredentials.New(uaaURL, true, uaa.Credentials["client_id"].(string),
		uaa.Credentials["client_secret"].(string))

	if nil != err {
		log.Fatal("Failed to setup uaa connection")
		log.Fatal(err)
		os.Exit(1)
	}

	token, err := uaaConnection.GetBearerToken()

	if nil != err {
		log.Fatal("Failed to get token from UAA")
		log.Fatal(err)
		os.Exit(1)
	}

}
