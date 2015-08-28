package tests

import (
	"encoding/json"
	"fm-fuel-service/object"
	"testing"
	"time"
)

var (
	apiAddFuel    = apiEndpoint{"POST", "/fuel"}
	apiModifyFuel = apiEndpoint{"POST", "/fuel/:id"}
)

func TestAddFuelApi(t *testing.T) {
	var fuel object.Fuel

	// set some default data
	fuel.Fleet = "ALKF-DLFK-DFLJ-DLFKJ"
	fuel.Vehicle = "ALKF-DLFK-ASDL-DLFKJ"
	fuel.FillDate = time.Now()
	// do request to working api server
	body, err := jsonReq(apiAddFuel, fuel)
	if err != nil {
		t.Error(err)
		return
	}

	fuelReceived := object.Fuel{}
	if err := json.Unmarshal(body, &fuelReceived); err != nil {
		t.Error(err)
	}
}

func TestModifyFuelApi(t *testing.T) {
	var fuel object.Fuel

	fuel.Fleet = "ALKF-DLFK-DFLJ-DLFKJ"
	fuel.Vehicle = "ALKF-DLFK-ASDL-DLFKJ"
	fuel.FillDate = time.Now()

	// send fuel json object
	body, err := jsonReq(apiModifyFuel, fuel)
	if err != nil {
		t.Error(err)
		return
	}
	fuelReceived := object.Fuel{}

	if err := json.Unmarshal(body, &fuelReceived); err != nil {
		t.Error(err)
	}

}
