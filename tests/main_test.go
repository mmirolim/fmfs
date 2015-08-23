package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

type apiEndpoint struct {
	Method, Url string
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func jsonReq(ae apiEndpoint, load interface{}, host ...string) ([]byte, error) {
	//@todo set port from configuration?
	apiHost := "http://localhost:4001"
	// if host provided set it instead of default one
	if len(host) == 1 {
		apiHost = host[0]
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(load)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(ae.Method, apiHost+ae.Url, &buf)
	req.Header.Add("Content-Type", "application/json")

	// disable keep alive in client, too make new connection each time
	// if not disabled, there is error EOF
	client := &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return body, err
}
