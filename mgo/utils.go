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
	masterSession *mgo.Session
}

// Dial connects to the db using specified options and returns helper or error if dial fails.
func Dial(dialInfo *mgo.DialInfo, opts *DbHelperInitOptions) (d *DbHelper, err error) {
	d = new(DbHelper)
	d.masterSession, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to mongo: '%s'", err)
	}
	d.masterSession.SetSafe(opts.Safe)
	return d, nil
}

// C is a shortcut for getting a collection+session (master session copy) by name inside current database
func (db *DbHelper) C(collectionName string) (*mgo.Collection, *mgo.Session) {
	session := db.masterSession.New()
	return session.DB("").C(collectionName), session
}

// S is a shortcut for creating a session for current database
func (db *DbHelper) S() *mgo.Session {
	return db.masterSession.New()
}

// Cs is a shortcut for creating a collection using an existing session
func (db *DbHelper) Cs(collectionName string, session *mgo.Session) *mgo.Collection {
	return session.DB("").C(collectionName)
}