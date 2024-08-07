package handle

import (
	"bytes"
	"drone-ci-proxy/app"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// HandleHook intercept the request and deal with the handle
// -   If you will end the process here, please, write the StatusCode of Header (writer.WriteHeader) like below
// bool of return signify if the process should end and not forward to TARGET_HOST
// error of return signify if an error happened here
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

	var hookData map[string]interface{}
	err = json.Unmarshal(b, &hookData)

	if err != nil {
		log.Println("json.Unmarshal - ", err)
		return false, errors.New("can't process")
	}

	if hookData["action"] != "closed" {
		return false, nil
	}

	if hookData["pull_request"] == nil {
		return false, nil
	}

	if hookData["number"] == nil {
		return false, nil
	}

	//
	// TODO: Handle the closed action
	// Here you have all content of the hook to deal with it.
	app.Application.Log.Info("Hook received:", "hook", struct {
		Action      string
		PullRequest any
		Number      float64
	}{
		Action:      hookData["action"].(string),
		PullRequest: nil,
		Number:      hookData["number"].(float64),
	})

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("OK!"))

	return true, nil
}
