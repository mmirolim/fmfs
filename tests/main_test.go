package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	// @todo get from configuration
	apiHost = "http://localhost:4001"
)

type apiEndpoint struct {
	Method, Url string
}

// copy api endpoint, do not change default
func (ae apiEndpoint) copy() apiEndpoint {
	return ae
}

// add suffix for Url param
func (ae *apiEndpoint) suffix(str string) {
	ae.Url = ae.Url + str
}

// add url params
func (ae *apiEndpoint) params(m map[string]interface{}) {
	var params string
	// create url params string with &key=value from provided map
	for k, v := range m {
		params += fmt.Sprintf("&%v=%v", k, v)
	}
	// remove first &
	params = strings.Replace(params, "&", "", 1)
	ae.Url = ae.Url + "?" + params

}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func jsonReq(ae apiEndpoint, load interface{}, host ...string) ([]byte, error) {
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
	// return error if response status code is Internal error
	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New(res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return body, err
}

// expect func to compare two arguments
// if equal then it returns false
// otherwise true
/*
if expectInt(t, want, got, what) {
return
}
*/
func expectInt(t *testing.T, want, got int, what string) bool {
	if want == got {
		return false
	}
	t.Errorf("expected %s %d, got %d", what, want, got)
	return true
}

// expect func to compare and stop execution of test
func expectStr(t *testing.T, want, got string, what string) bool {
	if want == got {
		return false
	}
	t.Errorf("expected %s %s, got %s", what, want, got)
	return true
}
