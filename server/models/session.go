package models

import (
	"crypto/rand"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

/**
 * Session - Stores the current user and
 * application state.
 *
 * (c) 2015 Matthias Hannig
 */
type Session struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Token     string        `json:"token"`
	UserId    string        `json:"userId"`
	Ttl       time.Duration `json:"ttl"`
	CreatedAt time.Time     `json:"createdAt"`
}

/**
 * Get associated user
 */
func (self *Session) User() User {
	return UserById(self.UserId)
}

/**
 * Create session
 */

func NewSession(userId string) *Session {
	randomToken := make([]byte, 32)
	_, err := rand.Read(randomToken)
	if err != nil {
		log.Fatal("[Session] Could not generate token.")
	}

	session := Session{
		Id:     bson.NewObjectId(),
		Token:  string(randomToken),
		UserId: userId,
		Ttl:    time.Second * 500,
	}

	return &session
}
