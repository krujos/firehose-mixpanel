package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/cloudfoundry-community/cfenv"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/krujos/uaaclientcredentials"
)

var skipSSLVerify = strings.ToLower(os.Getenv("SSL_VERIFY")) == "false"
var subscriptionID = "FirehoseToMixPanel"

func dieIfError(msg string, err error) {
	if nil != err {
		log.Fatal(msg)
		log.Fatal(err)
		os.Exit(1)
	}
}

func root(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Hello World")
}

func setupHTTP(port int) {
	http.HandleFunc("/", root)
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		dieIfError("Could not listen for http", err)
	}()
}

func getToken(appEnv *cfenv.App) string {
	uaa, err := GetUserProvidedServiceByName("uaa", appEnv)
	dieIfError("Could not find uaa service, did you bind it?", err)

	uaaURL, err := url.Parse(uaa.Credentials["uri"].(string))
	dieIfError("Could not parse uaa URI", err)

	uaaConnection, err := uaaclientcredentials.New(uaaURL, skipSSLVerify, uaa.Credentials["client_id"].(string),
		uaa.Credentials["client_secret"].(string))

	dieIfError("Failed to setup uaa connection", err)

	token, err := uaaConnection.GetBearerToken()
	dieIfError("Failed to get token from UAA", err)
	return token
}

func connectToFirehose(appEnv *cfenv.App, token string) {
	doppler, err := GetUserProvidedServiceByName("doppler", appEnv)
	dieIfError("Failed to get doppler service", err)
	consumer := noaa.NewConsumer(doppler.Credentials["uri"].(string), &tls.Config{InsecureSkipVerify: skipSSLVerify}, nil)
	msgChan := make(chan *events.Envelope)
	go func() {
		defer close(msgChan)
		errorChan := make(chan error)
		if nil != err {
			panic(err)
		}
		go consumer.Firehose(subscriptionID, token, msgChan, errorChan)

		for err := range errorChan {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}
	}()
}

func main() {
	appEnv, _ := cfenv.Current()
	setupHTTP(appEnv.Port)
	token := getToken(appEnv)
	connectToFirehose(appEnv, token)
}
