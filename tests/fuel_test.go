// Integration testing suit for fuel api
package tests

import (
	"encoding/json"
	"fm-fuel-service/object"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var (
	apiGetFuel        = apiEndpoint{"GET", "/fuel/"}
	apiAddFuel        = apiEndpoint{"POST", "/fuel"}
	apiModifyFuel     = apiEndpoint{"POST", "/fuel/"}
	apiDelFuel        = apiEndpoint{"DELETE", "/fuel/"}
	apiUnDelFuel      = apiEndpoint{"POST", "/fuel-entries/"}
	apiDelFromStorage = apiEndpoint{"DELETE", "/fuel-entries/"}
	apiVehicle        = apiEndpoint{"GET", "/vehicle/"}
	apiFleet          = apiEndpoint{"GET", "/fleet/"}

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
	if fuel.CreatedBy == "" || fuel.CreatedAt.IsZero() {
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
	if !expectStr(t, fr.Info, fuel.Info, "fuel info") {
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
	if !expectInt(t, http.StatusNoContent, data.StatusCode, "http status") {
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

	if !expectInt(t, http.StatusNotFound, resp.StatusCode, "http status") {
		return
	}

}

// test undel soft deleted entry
func TestUnDelFuel(t *testing.T) {
	// create new entry
	fuel, err := addFuel()
	if err != nil {
		t.Error(err)
		return
	}
	// make api request to soft del it
	api := apiDelFuel.copy()
	api.suffix(fuel.ID.Hex())
	body, err := jsonReq(api, fuel)
	if err != nil {
		t.Error(err)
		return
	}
	var data struct{ StatusCode int }
	if err := json.Unmarshal(body, &data); err != nil {
		t.Error(err)
		return
	}
	// check for correct response
	if !expectInt(t, http.StatusNoContent, data.StatusCode, "http status") {
		return
	}
	// now restore soft deleted item
	api = apiUnDelFuel.copy()
	api.suffix(fuel.ID.Hex())
	var fr object.Fuel
	body, err = jsonReq(api, fr)
	if err != nil {
		t.Error(err)
		return
	}
	if err := json.Unmarshal(body, &fr); err != nil {
		t.Error(err)
		return
	}
	// check object
	if fr.UpdatedBy == "" || fr.DeletedBy != "" {
		t.Errorf("restored object should have deleted properties unset, obj received %+v\n", fr)
		return
	}

}

// completely remove object from storage
// test undel soft deleted entry
func TestDelFuelFromStorage(t *testing.T) {
	// create new entry
	fuel, err := addFuel()
	if err != nil {
		t.Error(err)
		return
	}
	// make api request to soft del it
	api := apiDelFromStorage.copy()
	api.suffix(fuel.ID.Hex())
	body, err := jsonReq(api, fuel)
	if err != nil {
		t.Error(err)
		return
	}

	var data struct{ StatusCode int }
	if err := json.Unmarshal(body, &data); err != nil {
		t.Error(err)
		return
	}
	// check for correct response
	if !expectInt(t, http.StatusNoContent, data.StatusCode, "http status") {
		return
	}
}

// test get all fuel entries for a vehicle by filldate
func TestGetVehicleFuelInPeriod(t *testing.T) {
	quantity := 5
	fillInterval := 24 * time.Hour
	// load test data with different filldate
	fuelEntries := make([]object.Fuel, quantity)
	// create 5 fuel object
	// with one day difference of filldate
	layout := "02 Jan 06 15:04 MST"
	startDate, err := time.Parse(layout, "01 Jan 09 15:04 MST")
	if err != nil {
		t.Error(err)
		return
	}
	filldate := startDate
	for i, _ := range fuelEntries {
		fuelEntries[i] = newDummyFuel()
		fuelEntries[i].FillDate = filldate
		// incr by one day
		filldate = filldate.Add(fillInterval)
	}
	// load data to api
	for _, v := range fuelEntries {
		f, err := addFuel(v)
		if err != nil {
			t.Error(err)
			return
		}
		if f.CreatedAt.IsZero() {
			t.Error("load test data failed")
			return
		}
	}

	// now get fuel entries by vehicle in period
	// format urlapi + params
	// ?sd=js.Date.toJson&ed=js.Date.toJson
	// after start date we have 5 entries with one day difference
	// let's get 3 of them
	// exclude startdate and enddate
	data, err := json.Marshal(startDate.Add(70 * time.Hour))
	if err != nil {
		t.Error(err)
		return
	}
	ed := string(data) // end date
	data, err = json.Marshal(startDate.Add(5 * time.Hour))
	if err != nil {
		t.Error(err)
		return
	}
	sd := string(data) // start date
	// make api request with params
	api := apiVehicle.copy()
	// get fuel object with appropriate vehicle set
	fuel := newDummyFuel()
	api.suffix(fuel.Vehicle)
	api.params(map[string]interface{}{
		"sd": sd,
		"ed": ed,
	})
	resp, err := http.Get(apiHost + api.Url)
	if err != nil {
		t.Error(err)
		return
	}
	// first check status
	if !expectInt(t, http.StatusOK, resp.StatusCode, "http status") {
		return
	}
	// then check that we have correct result
	// array of fuel entries
	data, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("all fule vehicle data", string(data))
	var fuels []object.Fuel
	err = json.Unmarshal(data, &fuels)
	fmt.Println("rec fuels", fuels)
	if err != nil {
		t.Error(err)
		return
	}
	// check that received object is in correct interval
	if len(fuels) == 0 {
		t.Error("resp from api is empty, expected fuel entries by vehicle")
		return
	}
	for _, v := range fuels {
		if !(v.FillDate.After(startDate) && v.FillDate.Before(filldate)) {
			t.Error("time interval of received fuel entries for vehicle is wrong")
			return
		}
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
