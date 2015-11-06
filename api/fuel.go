/*
API handlers for managing FUEL entries
*/
package api

import (
	"encoding/json"
	ds "fm-fuel-service/datastore"
	"fm-fuel-service/object"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2/bson"
)

// goji handlers for fuel object
// @todo add routes to get one/many soft deleted entries

// add new fuel entry
func addFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	userid, err := userID(r)
	if isErr(w, r, "userID", err) || userid == "" {
		log(r, "userid missing").Error(err)
		response(w, "userid missing", http.StatusBadRequest)
		return
	}

	// parse post data and decode to fuel object
	fuel, err := decodeFuel(r.Body)
	if isErr(w, r, "decode r.Body", err) {
		return
	}
	// set created time and user
	fuel.Created(userid)
	// save object
	err = ds.Save(&fuel)
	if isErr(w, r, "save fuel", err) {
		return
	}

	response(w, fuel)
}

// modify entry
func modifyFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	userid, err := userID(r)
	if isErr(w, r, "userID", err) || userid == "" {
		log(r, "userid missing").Error(err)
		response(w, "userid missing", http.StatusBadRequest)
		return
	}

	oid := c.URLParams["oid"]
	fuel, err := decodeFuel(r.Body)
	if isErr(w, r, "decode r.Body", err) {
		return
	}

	fuel.Updated(userid)
	err = ds.UpdateById(&fuel, oid)
	if isErr(w, r, "ds.UpdateById", err) {
		return
	}

	response(w, fuel)
}

// delete entry, soft delete used
// object not removed from storage
func delFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	userid, err := userID(r)
	if isErr(w, r, "userID", err) || userid == "" {
		log(r, "userid missing").Error(err)
		response(w, "userid missing", http.StatusBadRequest)
		return
	}
	var fuel object.Fuel
	oid := c.URLParams["oid"]
	err = ds.FindById(&fuel, oid)
	if isErr(w, r, "FindById", err) {
		return
	}
	// set del fields
	fuel.Deleted(userid)
	// set update flds
	fuel.Updated(userid)
	err = ds.UpdateById(&fuel, oid)

	if isErr(w, r, "UpdateById", err) {
		return
	}

	response(w, http.StatusNoContent)
}

// delete entry, soft delete used
// object not removed from storage
func delFuelFromStorage(c web.C, w http.ResponseWriter, r *http.Request) {
	userid, err := userID(r)
	if isErr(w, r, "userID", err) || userid == "" {
		log(r, "userid missing").Error(err)
		response(w, "userid missing", http.StatusBadRequest)
		return
	}
	oid := c.URLParams["oid"]
	err = ds.DelById(&object.Fuel{}, oid)
	if isErr(w, r, "DelById", err) {
		return
	}
	// @todo send event that fuel entry permanantly
	// deleted by userid
	response(w, http.StatusNoContent)
}

// restore soft deleted fuel-entry
func unDelFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	userid, err := userID(r)
	if isErr(w, r, "userID", err) || userid == "" {
		log(r, "userid missing").Error(err)
		response(w, "userid missing", http.StatusBadRequest)
		return
	}
	var fuel object.Fuel
	oid := c.URLParams["oid"]
	// get entry from storage
	// should let find soft deleted entries
	err = ds.FindById(&fuel, oid, true)
	if isErr(w, r, "FindById", err) {
		return
	}
	// unset deleted values
	fuel.DeletedBy = ""
	fuel.DeletedAt = time.Time{}
	// update system fields
	fuel.Updated(userid)
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

	// gen selector for find
	query, err := queryForPeriod(c, r, "vehicle", "uid", "filldate", "sd", "ed")
	if isErr(w, r, "queryForPeriod", err) {
		return
	}

	// use default limit
	err = ds.Find(&object.Fuel{}, query).All(&fuels)
	if isErr(w, r, "Find", err) {
		return
	}

	response(w, fuels)
}

// get entries for provided period for all vehicles in fleet
// @todo add limit param
func getFleetFuelInPeriod(c web.C, w http.ResponseWriter, r *http.Request) {
	var fuels []object.Fuel
	// gen selector for find
	query, err := queryForPeriod(c, r, "fleet", "uid", "filldate", "sd", "ed")
	if isErr(w, r, "queryForPeriod", err) {
		return
	}
	// use default limit
	err = ds.Find(&object.Fuel{}, query).All(&fuels)
	if isErr(w, r, "Find", err) {
		return
	}

	response(w, fuels)
}

// get user id (param key is userid) from url param
// user id must be UUID
func userID(r *http.Request) (string, error) {
	// get all urls params from raw query
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", err
	}

	return vals.Get("userid"), nil
}

// get required params from routing url and url params and
// generate bson.D selector
// @todo add limit param
func queryForPeriod(c web.C, r *http.Request, obj, id, fldname, start, end string) (bson.D, error) {
	// object id
	oid := c.URLParams[id]
	// get all urls params from raw query
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	var sd, ed time.Time
	// @todo add limit param to limit number of results
	err = json.Unmarshal([]byte(vals.Get(start)), &sd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(vals.Get(end)), &ed)
	if err != nil {
		return nil, err
	}
	// find all fuel entries in date interval and fleet
	query := ds.ByDateInterval(fldname, sd, ed)
	// add vehicle id to search
	query = append(query, bson.DocElem{Name: obj, Value: oid})

	return query, err
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
