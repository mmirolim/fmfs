package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

// goji handlers for fuel object

// add new fuel entry
func addFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "add fuel")
}

// modify entry
func modifyFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "modify fuel")
}

// delete entry
func delFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "del fuel")
}

// get entry
func getFuel(c web.C, w http.ResponseWriter, r *http.Request) {
	response(w, errors.New("something wrong"), 200)
}

// get entries for provided period for particular vehicle
func getVehicleFuelInPeriod(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get all fuel entries for period by vehicle")
}

// get entries for provided period for all vehicles in fleet
func getFleetFuelInPeriod(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get all fuel entries for period by fleet")
}
