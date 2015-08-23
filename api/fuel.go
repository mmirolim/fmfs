package api

import (
	"encoding/json"
	ds "fm-fuel-service/datastore"
	"fm-fuel-service/log"
	"fm-fuel-service/object"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

// goji handlers for fuel object

// add new fuel entry
func addFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	var fuel object.Fuel
	scope := "api.addFuel"
	// parse post data and decode to fuel object
	log.Print(scope, "register decoder").Debug(r.Body)
	err := json.NewDecoder(r.Body).Decode(&fuel)
	if isErr(w, scope, "decode r.Body", err, 400) {
		return
	}

	// @todo get user from jwt
	uid := "ALKJ-LDKFJ-DLFKJ-DLKFJ"
	// set created time and user
	fuel.Created(uid)
	// save object
	err = ds.Save(&fuel)
	if isErr(w, scope, "save fuel", err) {
		return
	}

	response(w, fuel)
}

// modify entry
func modifyFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "modify fuel")
}

// delete entry
func delFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	response(w, "fuel entry removed", 204)
}

// get entry
func getFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	fuel := object.Fuel{}
	id := c.URLParams["id"]
	err := ds.FindById(&fuel, id)
	if err != nil {
		response(w, err, http.StatusNotFound)
		return
	}

	response(w, fuel)
}

// get entries for provided period for particular vehicle
func getVehicleFuelInPeriod(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get all fuel entries for period by vehicle")
}

// get entries for provided period for all vehicles in fleet
func getFleetFuelInPeriod(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get all fuel entries for period by fleet")
}
