package api

import (
	"encoding/json"
	"fm-fuel-service/log"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

func New() *web.Mux {
	log.Print("api.Start", "some action").Debug("logging something")
	// WARNING more specific routes should be first then more general
	m := web.New()

	// set all routes for fuel object
	m.Post("/fuel", addFuel)
	// @todo maybe to use Patch, but should be tested first
	m.Post("/fuel/:oid", modifyFuel)
	m.Delete("/fuel/:oid", delFuel)
	m.Get("/fuel/:oid", getFuel)
	m.Get("/vehicle/:uid", getVehicleFuelInPeriod)
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
	if err != nil {
		code := 500
		log.Print(scope, action).Error(err)
		if len(status) == 1 {
			code = status[0]
		}
		response(w, err, code)
		return true
	}
	return false
}
