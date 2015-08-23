package tests

import (
	"bytes"
	"encoding/json"
	"fm-fuel-service/object"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAddFuelApi(t *testing.T) {
	var fuel object.Fuel

	fuel.Fleet = "ALKF-DLFK-DFLJ-DLFKJ"
	fuel.Vehicle = "ALKF-DLFK-ASDL-DLFKJ"

	var load bytes.Buffer
	err := json.NewEncoder(&load).Encode(fuel)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:4001/fuel", &load)
	req.Header.Add("Content-Type", "application/json")

	// disable keep alive in client, too make new connection each time
	// if not disabled, there is error EOF
	client := &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Error(err)
	}
	fuelReceived := object.Fuel{}

	if err := json.Unmarshal(body, &fuelReceived); err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", fuelReceived)

}
