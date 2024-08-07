package proxy

import (
	"bytes"
	"context"
	setup "drone-ci-proxy/app"
	"drone-ci-proxy/cmd/handle"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type DebugTransport struct{}

func (DebugTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(request)
}

func Proxy(writer http.ResponseWriter, request *http.Request) {

	if request.URL.Path == "/hook" {

		isFinished, err := handle.HandleHook(writer, request)

		if err != nil {
			writer.WriteHeader(http.StatusBadGateway)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		if isFinished {
			return
		}
	}

	proxyRequest := request.WithContext(context.TODO())
	reverseProxy, err := createReverseProxy()

	if err != nil {
		writer.WriteHeader(http.StatusBadGateway)
		_, _ = writer.Write([]byte("Reverse proxy not available!"))
		log.Println("createReverseProxy - ", err)
		return
	}

	log.Println("Forwarding request to reverse proxy")

	reverseProxy.ModifyResponse = rewriteBody
	reverseProxy.ServeHTTP(writer, proxyRequest)
}

func rewriteBody(resp *http.Response) (err error) {

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	body := io.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func reverseProxyErrorHandler(writer http.ResponseWriter, request *http.Request, err error) {
	log.Println("ReverseProxyErrorHandler: ")
	log.Println(err.Error())

	b, err := io.ReadAll(request.Body)

	if err != nil {
		log.Println("io.ReadAll - ", err)
		writer.WriteHeader(http.StatusBadGateway)
		_, _ = writer.Write([]byte("Reverse proxy not available!"))
		return
	}

	err = request.Body.Close()

	if err != nil {
		log.Println("request.Body.Close - ", err)
		writer.WriteHeader(http.StatusBadGateway)
		_, _ = writer.Write([]byte("Reverse proxy not available!"))
		return
	}

	body := io.NopCloser(bytes.NewReader(b))
	request.Body = body

	var result map[string]interface{}
	err = json.Unmarshal(b, &result)

	if err != nil {
		log.Println("json.Unmarshal - ", err)
		writer.WriteHeader(http.StatusBadGateway)
		_, _ = writer.Write([]byte("Reverse proxy not available!"))
		return
	}

	log.Println(result)
	writer.WriteHeader(http.StatusBadGateway)
	_, _ = writer.Write([]byte("Reverse proxy not available!"))
}

func createReverseProxy() (*httputil.ReverseProxy, error) {

	director, err := getReverseProxyDirector()

	if err != nil {
		return nil, err
	}

	proxy := &httputil.ReverseProxy{
		Director:     director,
		Transport:    DebugTransport{},
		ErrorHandler: reverseProxyErrorHandler,
	}

	return proxy, err
}

func getReverseProxyDirector() (func(request *http.Request), error) {
	var err error
	var hostAddress *url.URL

	hostAddress, err = url.Parse(setup.Application.TargetHost)

	if err != nil {
		return nil, err
	}

	director := func(request *http.Request) {
		request.Host = hostAddress.Host
		request.URL.Scheme = hostAddress.Scheme
		request.URL.Host = hostAddress.Host
	}

	return director, nil
}
