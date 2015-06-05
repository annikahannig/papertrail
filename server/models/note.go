package models

import (
	"errors"
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
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
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
 * Create new note
 */
func NewNote(text string, user User) Note {
	note := Note{
		Text:        text,
		CreatedById: user.Id,
	}

	return note
}

func (self *Note) Save() error {
	var err error
	if self.Id.Valid() {
		// Update Note
		err = UpdateNote(self)
	} else {
		// New Note: Insert.
		err = InsertNote(self)
	}
	return err
}

/**
 * Access user via created by id
 */
func (self *Note) CreatedBy() (*User, error) {
	if self.CreatedById == "" {
		return nil, errors.New("user not set for note")
	}
	user, err := FindUserById(self.CreatedById)

	return user, err
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

/**
 * Update note
 */
func UpdateNote(note *Note) error {
	// Get notes from collection
	c := db.C("notes")

	// Assert defaults for some fields
	note.UpdatedAt = time.Now()
	err := c.UpdateId(note.Id, note)
	return err
}
