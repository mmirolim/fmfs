package api

import (
	"bufio"
	"encoding/json"
	"fm-fuel-service/object"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zenazn/goji/web"
)

// test handler performance
func BenchmarkAddFuelApi(b *testing.B) {
	var fuel object.Fuel
	var c web.C
	c.Env = make(map[interface{}]interface{})

	// set some default data
	fuel.Fleet = "ALKF-DLFK-DFLJ-DLFKJ"
	fuel.Vehicle = "ALKF-DLFK-ASDL-DLFKJ"
	fuel.FillDate = time.Now()

	data, err := json.Marshal(fuel)
	if err != nil {
		b.Error(err)
		return
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rw := httptest.NewRecorder()
		r := bufio.NewReader(strings.NewReader(string(data)))
		req, err := http.NewRequest("POST", "http://localhost:4040", r)
		if err != nil {
			b.Error(err)
			return
		}
		addFuel(c, rw, req)
	}
}
