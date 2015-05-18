package models

import (
	"github.com/dchest/uniuri"
	"gopkg.in/mgo.v2/bson"
	"time"
)

/**
 * Session - Stores the current user and
 * application state.
 *
 * (c) 2015 Matthias Hannig
 */
type Session struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Token     string        `json:"token"`
	UserId    string        `json:"userId"`
	Lifetime  time.Duration `json:"lifetime"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	TouchedAt time.Time     `json:"touchedAt"`
}

/**
 * Get associated user
 */
func (self *Session) User() (*User, error) {
	return FindUserById(self.UserId)
}

/**
 * Touch session
 */
func (self *Session) Touch() {
	self.TouchedAt = time.Now()
	self.Save()
}

/**
 * Save / Update Session
 */
func (self *Session) Save() error {
	var err error
	c := db.C("sessions")

	if self.Id.Valid() {
		// Update session
		self.UpdatedAt = time.Now()
		err = c.UpdateId(self.Id, self)
	} else {
		// Insert fresh session
		self.CreatedAt = time.Now()
		self.TouchedAt = time.Now()
		self.Id = bson.NewObjectId()
		err = c.Insert(self)
	}

	return err
}

/**
 * Destroy Session
 */
func (self *Session) Destroy() error {
	c := db.C("sessions")
	err := c.RemoveId(&self.Id)
	return err
}

/**
 * Calculate TTL
 */
func (self *Session) TTL() time.Duration {
	ttl := self.Lifetime - time.Now().Sub(self.TouchedAt)
	if ttl < 0 {
		ttl = 0
	}

	return ttl
}

/**
 * Create session
 */
func NewSession(userId string) *Session {
	randomToken := uniuri.NewLen(48)

	session := Session{
		Token:    randomToken,
		UserId:   userId,
		Lifetime: time.Second * 500,
	}

	return &session
}

/**
 * Find session by token
 */

func FindSessionByToken(sessionToken string) (*Session, error) {
	session := Session{}
	c := db.C("sessions")
	err := c.Find(bson.M{"token": sessionToken}).One(&session)
	return &session, err
}
