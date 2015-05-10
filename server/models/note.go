package models

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

/**
 * A single note
 */

type Note struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Text      string        `json:"text"`
	PrintedAt string        `json:"printedAt"`
	Draft     bool          `json:"isDraft"`

	CreatedById string    `json:"createdById"`
	CreatedAt   time.Time `json:"created_at"`
}

/**
 * Get all notes from DB
 */
func AllNotes() []Note {
	var result []Note

	// Get notes from collection
	c := db.C("notes")

	err := c.Find(bson.M{}).Sort("+CreatedAt").All(&result)
	if err != nil {
		log.Println("[DB] Error while fetching notes.")
	}

	return result
}

/**
 * Insert Note
 */
func InsertNote(note *Note) error {
	// Get notes from collection
	c := db.C("notes")

	// Assert defaults for some fields
	note.CreatedAt = time.Now()
	note.Id = bson.NewObjectId()

	err := c.Insert(note)

	return err
}
