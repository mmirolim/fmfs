// Integration testing suit for fuel api
package tests

import (
	"encoding/json"
	"fm-fuel-service/object"
	"net/http"
	"testing"
	"time"
)

var (
	apiGetFuel    = apiEndpoint{"GET", "/fuel/"}
	apiAddFuel    = apiEndpoint{"POST", "/fuel"}
	apiModifyFuel = apiEndpoint{"POST", "/fuel/"}
	apiDelFuel    = apiEndpoint{"DELETE", "/fuel/"}

	// test user UUID
	dummyUserID = "069c3cc2-41c1-4ae9-8b08-c80cf6ea12e9"
)

// generate new dummy fuel for testing
func newDummyFuel() object.Fuel {
	var fuel object.Fuel
	fuel = object.Fuel{
		Fleet:     "ccec31b5-08f8-4a78-bbee-ad1edea0fd4c",
		Vehicle:   "4dee251a-0f4a-4d98-84d8-658b31c2e670",
		FillDate:  time.Now(),
		FuelUnit:  object.LITRES,
		FuelType:  object.DIESEL,
		Amount:    100,
		TankSize:  100,
		Currency:  "USD",
		TotalCost: 100,
		Info:      "dummy fuel entry for testing purposes",
	}
	return fuel
}
func TestAddFuelApi(t *testing.T) {
	// should return received fuel json object
	// unmarshaled to fuel object
	fuel, err := addFuel()
	if err != nil {
		t.Error(err)
		return
	}
	// if created successfully fuel object should have
	// CreatedBy and CreatedAt properties set
	if fuel.CreatedBy == "" || fuel.CreatedAt.Year() == 1 {
		t.Error("created fuel object should CreatedAt and CreatedBy properties set")
	}
}

func TestModifyFuelApi(t *testing.T) {
	fuel := newDummyFuel()
	fuel.Info = "info-modified"
	// send fuel json object
	api := apiModifyFuel.copy()
	api.suffix(fuel.ID.String())
	// make api request
	body, err := jsonReq(api, fuel)
	if err != nil {
		t.Error(err)
		return
	}

	fr := object.Fuel{}
	if err := json.Unmarshal(body, &fr); err != nil {
		t.Error(err)
		return
	}
	if expectStr(t, fr.Info, fuel.Info, "fuel info") {
		return
	}

}

// test soft delete
func TestDelFuelApi(t *testing.T) {
	// first add new fuel
	fuel := newDummyFuel()
	// do request to working api server to add fuel
	body, err := jsonReq(apiAddFuel, fuel)
	if err != nil {
		t.Error(err)
		return
	}
	// get received object
	fr := object.Fuel{}
	if err := json.Unmarshal(body, &fr); err != nil {
		t.Error(err)
		return
	}
	// now delete it
	// do request to working api server
	api := apiDelFuel.copy()
	api.suffix(fuel.ID.String())
	body, err = jsonReq(api, fuel)
	if err != nil {
		t.Error(err)
		return
	}

	if err := json.Unmarshal(body, &fr); err != nil {
		t.Error(err)
		return
	}

	// try to get it by api
	// it should not get it
	api = apiGetFuel.copy()
	api.suffix(fr.ID.String())
	resp, err := http.Get(api.Url)
	if err != nil {
		t.Error(err)
		return
	}

	if expectInt(t, resp.StatusCode, http.StatusNotFound, "http status") {
		return
	}

}

// request to add fuel entry to api
// with passed fuel entry and return received
// object
// first object will be marshaled and sent
// if no arguments passed dummyFuel entry will be used
func addFuel(fuelObj ...object.Fuel) (object.Fuel, error) {
	fuel := newDummyFuel()
	fuelReceived := object.Fuel{}
	// if fuel provided use it
	if len(fuelObj) == 1 {
		fuel = fuelObj[0]
	}
	// do request to working api server
	body, err := jsonReq(apiAddFuel, fuel)
	if err != nil {
		return fuelReceived, err
	}

	err = json.Unmarshal(body, &fuelReceived)

	return fuelReceived, err
}
