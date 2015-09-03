package api

import (
	"encoding/json"
	ds "fm-fuel-service/datastore"
	"fm-fuel-service/object"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2/bson"
)

var (
	//@todo get uid from jwt
	uid = "QOIO-EOIL-EIRU-JLKL"
)

// goji handlers for fuel object
// @todo add routes to get one/many soft deleted entries

// add new fuel entry
func addFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	// parse post data and decode to fuel object
	fuel, err := decodeFuel(r.Body)
	if isErr(w, r, "decode r.Body", err) {
		return
	}
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
	err := ds.FindById(&fuel, oid)
	if isErr(w, r, "FindById", err) {
		return
	}
	// set del fields
	fuel.Deleted(uid)
	err = ds.UpdateById(&fuel, oid)

	if isErr(w, r, "UpdateById", err) {
		return
	}

	response(w, http.StatusNoContent)
}

// delete entry, soft delete used
// object not removed from storage
func delFuelFromStorage(c web.C, w http.ResponseWriter, r *http.Request) {
	oid := c.URLParams["oid"]
	err := ds.DelById(&object.Fuel{}, oid)
	if isErr(w, r, "DelById", err) {
		return
	}

	response(w, http.StatusNoContent)
}

// restore soft deleted fuel-entry
func unDelFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	var fuel object.Fuel
	oid := c.URLParams["oid"]
	// get entry from storage
	// should let find soft deleted entries
	err := ds.FindById(&fuel, oid, true)
	if isErr(w, r, "FindById", err) {
		return
	}
	// unset deleted values
	fuel.DeletedBy = ""
	fuel.DeletedAt = time.Time{}
	// update system fields
	fuel.Updated(uid)
	// update object
	err = ds.UpdateById(&fuel, oid)
	if isErr(w, r, "UpdateById", err) {
		return
	}

	response(w, fuel)
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
	var fuels []object.Fuel
	vehicleId := c.URLParams["oid"]
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if isErr(w, r, "ParseQuery", err) {
		return
	}

	var sd, ed time.Time
	err = json.Unmarshal([]byte(vals.Get("sd")), &sd)
	err = json.Unmarshal([]byte(vals.Get("ed")), &ed)
	if isErr(w, r, "unmarshal", err) {
		return
	}
	// find all fuel entries in date interval and fleet
	query := ds.ByDateInterval("filldate", sd, ed)
	// add vehicle id to search
	query = append(query, bson.DocElem{Name: "vehicle", Value: bson.M{"$eq": vehicleId}})
	err = ds.Find(&object.Fuel{}, query).All(&fuels)
	if isErr(w, r, "Find", err) {
		return
	}
	fmt.Println("FUEL from API", fuels)
	response(w, fuels)
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
