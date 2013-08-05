package mgo

import (
	"fmt"
	"labix.org/v2/mgo"
)

// DbHelperInitOptions contain initialization options for DbHelper constructor.
type DbHelperInitOptions struct {
	Safe *mgo.Safe // Safe mode set for DbHelper session.
}

// DbHelper is a convenience helper for systems that work with a single mongo DB.
// It creates a 'master' session, given the mgo url/db and init options and then
// uses its copies for db operations.
type DbHelper struct {
	dbName        string // used mongodb database name
	masterSession *mgo.Session
}

func NewDbHelper(mgoUrl string, dbName string, opts *DbHelperInitOptions) (d *DbHelper, err error) {
	d = new(DbHelper)
	d.masterSession, err = mgo.Dial(mgoUrl)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to mongo: '%s'", err)
	}
	d.masterSession.SetSafe(opts.Safe)
	d.dbName = dbName
	return d, nil
}

// c is a shortcut for getting a collection by name inside current database
func (db *DbHelper) c(collectionName string) (*mgo.Collection, *mgo.Session) {
	session := db.masterSession.New()
	return session.DB(db.dbName).C(collectionName), session
}

// s is a shortcut for creating a session for current database
func (db *DbHelper) s() *mgo.Session {
	return db.masterSession.New()
}

// cs is a shortcut for creating a collection using an existing session
func (db *DbHelper) cs(collectionName string, session *mgo.Session) *mgo.Collection {
	return session.DB(db.dbName).C(collectionName)
}