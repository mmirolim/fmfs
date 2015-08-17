package datastore

import (
	"fm-fuel-service/conf"
	"fm-fuel-service/objects"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

// variables
var (
	mgoHosts  []string // list of host to dial to
	mgoDbName string
	msess     *mgo.Session // mongo session
)

//Define Document interface for Document based nosql db
type Document interface {
	Collection() string                 // get collection of document
	Index() []mgo.Index                 // get slice of all indexes required by document
	Created(by string, at ...time.Time) // who created and when
	Updated(by string, at ...time.Time) // who updated and when
	Deleted(by string, at ...time.Time) // who deleted and when
}

type DocumentWithLocation interface {
	Document
	Location() object.Geo // get object should have Geo properties
}

// Initialize datastore package
func Initialize(app conf.App) error {
	// read all mongo hosts and join to one string
	mgoHosts = strings.Join(app.DS.MongoHosts, ",")

	// set mongo db name
	mgoDbName = app.DS.MongoDbName

	return nil
}
