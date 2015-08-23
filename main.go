package main

import (
	"flag"
	"fm-fuel-service/api"
	"fm-fuel-service/conf"
	"fm-fuel-service/datastore"
	"fm-fuel-service/log"
	"fm-fuel-service/object"

	"github.com/zenazn/goji"
)

func init() {
	// run before main
	log.Start()
	// set configs
	appConf := conf.App{}
	appConf.DS.Mongo.SetHosts("localhost:27017")
	appConf.DS.Mongo.SetDB("fuel")
	// init datastore
	err := datastore.Initialize(appConf.DS.Mongo, &object.Fuel{})
	if err != nil {
		log.Print("main.init", "datasore.Initialize").Fatal(err)
	}
}

func main() {
	// get goji mux from api package
	mux := api.New()
	// set goji server port
	flag.Set("bind", "localhost:4001")
	// set JSON middleware
	goji.Use(api.JSON)
	// register routes
	goji.Handle("/*", mux)
	/// start server
	goji.Serve()
}
