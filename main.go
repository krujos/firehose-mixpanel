package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func root(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Hello World")
}

func setupHTTP() {
	http.HandleFunc("/", root)
	go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if nil != err {
			log.Fatal("ListenAndServe:", err)
		}
	}()
}

func main() {
	setupHTTP()
}
