package setup

import (
	"log"
	"net/http"
	"os"
)

type Config struct {
	WebPort    string
	TargetHost string
}

var Application = Config{
	WebPort:    os.Getenv("WEB_PORT"),
	TargetHost: os.Getenv("TARGET_HOST"),
}

func init() {

	if Application.WebPort == "" {
		log.Fatalln("WebPort not configured!")
	}

	if Application.TargetHost == "" {
		log.Fatalln("TargetHost not configured!")
	}

	http.DefaultTransport.(*http.Transport).MaxIdleConns = 1000
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}
