package models

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	// "os"
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Password  []byte    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Users map[string]User

func NewUser(username string, password string, name string) User {

	// Generate random ID
	id := make([]byte, 16)
	_, err := rand.Read(id)
	if err != nil {
		log.Fatal("[User] Could not generate id.")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("[User] Could not hash password.")
	}

	user := User{
		Id:       hex.EncodeToString(id),
		Username: username,
		Password: hash,
		Name:     name,
	}

	return user
}

func (self *User) Save() {
	self.CreatedAt = time.Now()
	users := LoadUsers("./data/users.json")
	users.AddUser(*self)
	users.Save("./data/users.json")
}

func AllUsers() *Users {
	users := LoadUsers("./data/users.json")
	return users
}

func UserById(userid string) User {
	users := AllUsers()
	return (*users)[userid]
}

func AuthenticateUser(username string, password string) bool {

	return false
}

func (self *Users) Save(filename string) {
	res, err := json.Marshal(self)
	if err != nil {
		log.Fatal("[User] Could not serialize users map")
	}

	err = ioutil.WriteFile(filename, []byte(res), 0644)
	if err != nil {
		log.Fatal("[Users] Could not write userdb.")
	}
}

func (self *Users) AddUser(user User) {
	(*self)[user.Username] = user
}

func LoadUsers(filename string) *Users {
	users := make(Users)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return &users
	}

	err = json.Unmarshal(content, &users)
	if err != nil {
		return &users
	}

	return &users
}
