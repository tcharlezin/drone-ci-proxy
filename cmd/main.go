package main

import (
	setup "drone-ci-proxy/app"
	"drone-ci-proxy/cmd/proxy"
	"fmt"
	"log"
	"net/http"
)

func main() {

	mutex := http.NewServeMux()
	mutex.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK!"))
	})
	mutex.HandleFunc("/", proxy.Proxy)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", setup.Application.WebPort),
		Handler: mutex,
	}

	log.Println(fmt.Sprintf("Starting server on :%s", setup.Application.WebPort))
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}
