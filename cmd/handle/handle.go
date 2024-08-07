package handle

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func HandleHook(writer http.ResponseWriter, request *http.Request) (bool, error) {

	b, err := io.ReadAll(request.Body)

	if err != nil {
		log.Println("io.ReadAll - ", err)
		return false, errors.New("can't process")
	}

	err = request.Body.Close()

	if err != nil {
		log.Println("request.Body.Close - ", err)
		return false, errors.New("can't process")
	}

	body := io.NopCloser(bytes.NewReader(b))
	request.Body = body

	var result map[string]interface{}
	err = json.Unmarshal(b, &result)

	if err != nil {
		log.Println("json.Unmarshal - ", err)
		return false, errors.New("can't process")
	}

	if result["action"] != "closed" {
		return false, nil
	}

	if result["pull_request"] == nil {
		return false, nil
	}

	if result["number"] == nil {
		return false, nil
	}

	//
	// TODO: Handle the closed action
	// Here you have all content of the hook to deal with it.

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("OK!"))

	return true, nil
}
