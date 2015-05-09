package models

import (
	"gopkg.in/mgo.v2/bson"
)

/**
 * A single note
 */

type Note struct {
	Text        string        `json:"name"`
	PrintedAt   string        `json:"printed_at"`
	PrintedById bson.ObjectId `json:"printed_by_id"`
}
