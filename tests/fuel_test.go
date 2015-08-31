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
		t.Error("created fuel object should have CreatedAt and CreatedBy properties set")
	}
}

func TestModifyFuelApi(t *testing.T) {
	fuel, err := addFuel()
	if err != nil {
		t.Error(err)
		return
	}
	// make some change to fuel object
	fuel.Info = "info-modified"
	// send fuel json object
	api := apiModifyFuel.copy()
	api.suffix(fuel.ID.Hex())
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
	fuel, err := addFuel()
	if err != nil {
		t.Error(err)
		return
	}

	// now delete it
	// do request to working api server
	api := apiDelFuel.copy()
	api.suffix(fuel.ID.Hex())
	var data struct{ StatusCode int }
	body, err := jsonReq(api, data)
	if err != nil {
		t.Error(err)
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Error(err)
		return
	}
	if expectInt(t, http.StatusNoContent, data.StatusCode, "http status") {
		return
	}
	// try to get it by api
	// it should not get it
	// because it soft deleted
	api = apiGetFuel.copy()
	api.suffix(fuel.ID.Hex())
	resp, err := http.Get(apiHost + api.Url)
	if err != nil {
		t.Error(err)
		return
	}

	if expectInt(t, http.StatusNotFound, resp.StatusCode, "http status") {
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
	fr := object.Fuel{}
	// if fuel provided use it
	if len(fuelObj) == 1 {
		fuel = fuelObj[0]
	}
	// do request to working api server
	body, err := jsonReq(apiAddFuel, fuel)
	if err != nil {
		return fr, err
	}

	err = json.Unmarshal(body, &fr)

	return fr, err
}
