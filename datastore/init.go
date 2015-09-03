package datastore

import (
	"encoding/json"
	"errors"
	"fm-fuel-service/object"
	"fmt"
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
	// errors
	ErrObjectId = errors.New("id type wrong")
	ErrNotFound = mgo.ErrNotFound
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
func Initialize(mga MongoAdapter, docs ...Document) (err error) {
	// read all mongo hosts and join to one string
	mgoHost := strings.Join(mga.Hosts(), ",")

	// set mongo db name
	mgoDbName = mga.DB()

	// try to connect, set timeout for request
	msess, err = mgo.DialWithTimeout(mgoHost, time.Second)
	if err != nil {
		return
	}
	// init indexes set for collections
	for _, doc := range docs {
		if err = EnsureIndex(doc); err != nil {
			return
		}
	}

	return
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

// find all documents by query and limit
func Find(doc Document, query bson.D, limit ...int) *mgo.Iter {
	// default limit
	lim := 10000
	sess := msess.Copy()
	defer sess.Close()

	if len(limit) == 1 {
		lim = limit[0]
	}
	return getColl(sess, doc).Find(query).Limit(lim).Iter()
}

// create bson.D document query to find in interval
func ByDateInterval(fld string, start, end time.Time, deletedAlso ...bool) bson.D {
	var doc bson.D
	var elm bson.DocElem

	data, err := json.Marshal(start)
	sd := fmt.Sprintf("ISODate(%s)", string(data))
	data, err = json.Marshal(end)
	ed := fmt.Sprintf("ISODate(%s)", string(data))
	if err != nil {
		return doc
	}
	fmt.Println("sd, ed from ds.init()", sd, ed)
	doc = bson.D{
		{fld, bson.M{"$gte": sd, "$lte": ed}},
	}

	// check deletedAlso if not set search all docs
	if len(deletedAlso) == 1 && deletedAlso[0] {
		return doc
	}
	// by default search not deleted elements
	elm = bson.DocElem{Name: "base.deletedby", Value: ""}
	doc = append(doc, elm)

	return doc
}

// find one document by id, by default will not find
// soft deleted documents
func FindById(doc Document, id string, deletedAlso ...bool) error {
	// before queries check is id fits otherwise
	// it panics
	if !bson.IsObjectIdHex(id) {
		return ErrObjectId
	}

	sess := msess.Copy()
	defer sess.Close()

	// check deletedAlso search set
	if len(deletedAlso) == 1 && deletedAlso[0] {
		return getColl(sess, doc).FindId(bson.ObjectIdHex(id)).One(doc)
	}
	// by default find excludes soft deleted files by field "deletedby"

	return getColl(sess, doc).Find(bson.M{"_id": bson.ObjectIdHex(id), "base.deletedby": ""}).One(doc)
}

// save docuemnt to storage
func Save(doc Document) error {
	sess := msess.Copy()
	defer sess.Close()

	// set ObjectId
	doc.SetID(bson.NewObjectId())
	return getColl(sess, doc).Insert(doc)
}

// update document by id
func UpdateById(doc Document, id string) error {
	sess := msess.Copy()
	defer sess.Close()
	// before queries check is id fits otherwise
	// it panics
	if !bson.IsObjectIdHex(id) {
		return errors.New("id type wrong")
	}

	return getColl(sess, doc).UpdateId(bson.ObjectIdHex(id), doc)
}

// delete document from storage
func DelById(doc Document, id string) error {
	sess := msess.Copy()
	defer sess.Close()

	// before queries check is id fits otherwise
	// it panics
	if !bson.IsObjectIdHex(id) {
		return ErrObjectId
	}
	return getColl(sess, doc).RemoveId(bson.ObjectIdHex(id))
}

// return mongo Collection
func getColl(sess *mgo.Session, doc Document) *mgo.Collection {
	return sess.DB(mgoDbName).C(doc.Collection())
}
