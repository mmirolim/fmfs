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
	// complete remove object from storage
	m.Delete("/fuel/:oid", delFuel)
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
		data = struct{ Msg string }{v}
	case error:
		data = struct{ Err string }{v.Error()}
	}
	// json encode data
	b, err := json.Marshal(data)
	// check for errors
	if err != nil {
		res = err.Error()
	} else {
		res = string(b)
	}
	// check if status code set
	if len(status) == 1 {
		w.WriteHeader(status[0])
	}
	// respond
	fmt.Fprintf(w, res)
}

// on error prepare response and return true
func isErr(w http.ResponseWriter, scope, action string, err error, status ...int) bool {
	if err == nil {
		return false
	}
	code := http.StatusInternalServerError
	// set status code according to
	// err type/msg
	switch {
	case err == ds.ErrNotFound:
		code = http.StatusNotFound
	case len(status) == 1:
		// use provided code status
		code = status[0]
	}
	response(w, err, code)

	log.Print(scope, action).Error(err)
	return true
}
