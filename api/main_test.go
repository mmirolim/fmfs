package api

import (
	"fm-fuel-service/conf"
	"fm-fuel-service/datastore"
	"fm-fuel-service/log"
	"fm-fuel-service/object"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
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

	os.Exit(m.Run())
}
