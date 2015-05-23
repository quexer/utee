package utee

import (
	"gopkg.in/mgo.v2"
	"log"
)

func MongoConnect(db string) *mgo.Session {
	session, err := mgo.Dial(db)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func MongoNotFound(err error) bool {
	return mgo.ErrNotFound == err
}
