package api

import (
	"encoding/json"
	ds "fm-fuel-service/datastore"
	"fm-fuel-service/object"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/zenazn/goji/web"
)

// goji handlers for fuel object

// add new fuel entry
func addFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	scope := "api.addFuel"
	// parse post data and decode to fuel object
	fuel, err := decodeFuel(r.Body)
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
	scope := "api.modifyFuel"
	oid := c.URLParams["oid"]
	fuel, err := decodeFuel(r.Body)
	if isErr(w, scope, "decode r.Body", err, 400) {
		return
	}
	//@todo get uid from jwt
	uid := "QOIO-EOIL-EIRU-JLKL"
	fuel.Updated(uid)
	err = ds.UpdateById(&fuel, oid)
	if isErr(w, scope, "ds.UpdateById", err) {
		return
	}

	response(w, fuel)
}

// delete entry
func delFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	response(w, "fuel entry removed", 204)
}

// get entry
func getFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	scope := "api.getFuel"
	fuel := object.Fuel{}
	oid := c.URLParams["oid"]
	err := ds.FindById(&fuel, oid)
	if isErr(w, scope, "FindById", err, 404) {
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

// decode incoming fuel object json
func decodeFuel(rc io.ReadCloser) (object.Fuel, error) {
	var fuel object.Fuel
	data, err := ioutil.ReadAll(rc)
	// release resource
	// discard rest of input on err
	defer rc.Close()
	defer io.Copy(ioutil.Discard, rc)

	if err != nil {
		return fuel, err
	}
	err = json.Unmarshal(data, &fuel)

	return fuel, err
}
