package api

import (
	"encoding/json"
	ds "fm-fuel-service/datastore"
	"fm-fuel-service/log"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

func New() *web.Mux {
	// WARNING more specific routes should be first then more general
	m := web.New()

	// set all routes for fuel object
	m.Post("/fuel", addFuel)
	// @todo maybe to use Patch, but should be tested first
	// modify also used to soft delete object
	m.Post("/fuel/:oid", modifyFuel)
	// soft delete object
	m.Delete("/fuel/:oid", delFuel)
	// delete fuel entry from storage
	m.Delete("fuel-entries/:oid", delFuelFromStorage)
	// restore soft deleted object
	m.Post("/fuel-entries/:oid", unDelFuel)
	// get on fuel entry
	m.Get("/fuel/:oid", getFuel)
	// get all fuel entries for particular vehicle
	// for provided period in url params
	m.Get("/vehicle/:uid", getVehicleFuelInPeriod)
	// get all fuel entries for particualr fleet
	// for provided period in url params
	m.Get("/fleet/:uid", getFleetFuelInPeriod)

	return m
}

// set response to json format
func JSON(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// reply in json format with status code
func response(w http.ResponseWriter, data interface{}, status ...int) {
	var res string
	// create struct for json marshaling
	// struct depends on concrete type of data
	switch v := data.(type) {
	case string:
		// if custom string passed
		data = struct{ Msg string }{v}
	case error:
		// if err passed
		data = struct{ Err string }{v.Error()}
	case int:
		// if http status passed as data for response
		data = struct{ Msg int }{v}
	default:
		// if type we do not expect show as it is
		data = struct{ Data interface{} }{v}
	}
	// json encode data
	b, err := json.Marshal(data)
	// check for errors
	if err != nil {
		res = err.Error()
	} else {
		res = string(b)
	}
	// check if status code explicitly provided
	if len(status) == 1 {
		w.WriteHeader(status[0])
	}
	// respond
	fmt.Fprintf(w, res)
}

// on error prepare response and return true
func isErr(w http.ResponseWriter, r *http.Request, action string, err error, respStatus ...int) bool {
	if err == nil {
		return false
	}

	// get request method and url escaped path as err scope
	scope := r.Method + " " + r.URL.EscapedPath()
	// set default status
	status := http.StatusInternalServerError
	// set status code according to
	// err type/msg
	switch {
	case err == ds.ErrNotFound:
		status = http.StatusNotFound
	case len(respStatus) == 1:
		// use provided code status
		status = respStatus[0]
	}
	response(w, err, status)

	log.Print(scope, action).Error(err)
	return true
}
