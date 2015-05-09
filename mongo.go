package utee

import (
	"labix.org/v2/mgo"
	"log"
)

func ConnectMongo(db string) *mgo.Session {
	session, err := mgo.Dial(db)
	if err != nil {
		log.Fatal(err)
	}
	return session
}
