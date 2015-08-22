package main

import (
	"flag"
	"fm-fuel-service/api"
	"fm-fuel-service/log"

	"github.com/zenazn/goji"
)

func init() {
	log.Start()
}

func main() {
	// start http api server
	mux := api.New()
	// set goji server port
	flag.Set("bind", "localhost:4001")
	// set JSON middleware
	goji.Use(api.JSON)
	// register routes
	goji.Handle("/*", mux)
	// start server
	goji.Serve()
}
