package api

import (
	"encoding/json"
	ds "fm-fuel-service/datastore"
	"fm-fuel-service/object"
	"fmt"
	"io"
	"net/http"

	"github.com/zenazn/goji/web"
)

// goji handlers for fuel object

// add new fuel entry
func addFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	// parse post data and decode to fuel object
	fuel, err := decodeFuel(r.Body)
	if isErr(w, r, "decode r.Body", err) {
		return
	}
	//@todo get uid from jwt
	uid := "QOIO-EOIL-EIRU-JLKL"
	// set created time and user
	fuel.Created(uid)
	// save object
	err = ds.Save(&fuel)
	if isErr(w, r, "save fuel", err) {
		return
	}

	response(w, fuel)
}

// modify entry
func modifyFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	oid := c.URLParams["oid"]
	fuel, err := decodeFuel(r.Body)
	if isErr(w, r, "decode r.Body", err) {
		return
	}
	//@todo get uid from jwt
	uid := "QOIO-EOIL-EIRU-JLKL"
	fuel.Updated(uid)
	err = ds.UpdateById(&fuel, oid)
	if isErr(w, r, "ds.UpdateById", err) {
		return
	}

	response(w, fuel)
}

// delete entry, soft delete used
// object not removed from storage
func delFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	fuel := object.Fuel{}
	oid := c.URLParams["oid"]
	//@todo get uid from jwt
	uid := "QOIO-EOIL-EIRU-JLKL"
	fuel.Deleted(uid)
	err := ds.UpdateById(&fuel, oid)
	if isErr(w, r, "UpdateById", err) {
		return
	}

	response(w, http.StatusNoContent)
}

// get entry
func getFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	fuel := object.Fuel{}
	oid := c.URLParams["oid"]
	err := ds.FindById(&fuel, oid)
	if isErr(w, r, "FindById", err) {
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
	err := json.NewDecoder(rc).Decode(&fuel)
	// release resource
	rc.Close()
	return fuel, err
}
