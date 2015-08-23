package datastore

import (
	"errors"
	"fm-fuel-service/object"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// variables
var (
	mgoDbName string
	// @todo use singleton or sync.Once
	msess *mgo.Session // mongo session
)

type MongoAdapter interface {
	Hosts() []string
	DB() string
}

//Define Document interface for Document based nosql db
type Document interface {
	SetID(bson.ObjectId)
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
func Initialize(mga MongoAdapter, docs ...Document) error {
	var err error
	// read all mongo hosts and join to one string
	mgoHost := strings.Join(mga.Hosts(), ",")

	// set mongo db name
	mgoDbName = mga.DB()

	// try to connect, set timeout for request
	msess, err = mgo.DialWithTimeout(mgoHost, time.Second)
	if err != nil {
		return err
	}
	// init indexes set for collections
	for _, doc := range docs {
		if err := EnsureIndex(doc); err != nil {
			return err
		}
	}

	return err
}

// ensures index in collection created
func EnsureIndex(doc Document) error {
	var err error
	sess := msess.Copy()
	defer sess.Close()

	for _, v := range doc.Index() {
		err = getColl(sess, doc).EnsureIndex(v)
		if err != nil {
			break
		}
	}
	return err
}

// find one document by id
func FindById(doc Document, id string) error {
	var err error

	sess := msess.Copy()
	defer sess.Close()
	// before queries check is id fits otherwise
	// it panics
	if !bson.IsObjectIdHex(id) {
		return errors.New("id type wrong")
	}

	err = getColl(sess, doc).FindId(bson.ObjectIdHex(id)).One(doc)
	return err
}

// save docuemnt to storage
func Save(doc Document) error {
	var err error

	sess := msess.Copy()
	defer sess.Close()

	// set ObjectId
	doc.SetID(bson.NewObjectId())
	err = getColl(sess, doc).Insert(doc)

	return err
}

// update document by id
func UpdateById(doc Document, id string) error {
	var err error

	sess := msess.Copy()
	defer sess.Close()
	// before queries check is id fits otherwise
	// it panics
	if !bson.IsObjectIdHex(id) {
		return errors.New("id type wrong")
	}

	err = getColl(sess, doc).UpdateId(bson.ObjectIdHex(id), doc)
	return err
}

// return mongo Collection
func getColl(sess *mgo.Session, doc Document) *mgo.Collection {
	return sess.DB(mgoDbName).C(doc.Collection())
}
