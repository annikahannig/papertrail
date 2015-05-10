package models

import (
	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func SetDatabase(_db *mgo.Database) {
	db = _db
}
