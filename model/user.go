package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"time"
)

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Roles    []string  `json:"roles"`
	Token    []byte    `json:"token"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

func (User *User) GetName() string {
	//TODO implement me
	return User.Username
}

func (User *User) SetCreated() {
	User.Created = time.Now()
}

func (User *User) SetUpdated() {
	User.Updated = time.Now()
}

func NewAdminUser() *User {
	username := Env("ADMIN_USER", "root")
	pass := Env("ADMIN_USER_PASS", "123456")

	return (&User{
		Username: username,
		Roles:    []string{"admin"},
	}).SetToken(pass)
}

func (User *User) GenToken(pass string) []byte {
	token := sha256.Sum256([]byte(User.ID + ":" + pass))
	token1 := sha256.Sum256(token[:])
	return token1[:]
}

func (User *User) SetToken(pass string) *User {
	User.Token = User.GenToken(pass)
	return User
}

func (User User) PassEQ(token []byte) bool {
	return bytes.Equal(User.Token, token)
}

func (User User) IsAllow(roles []string) (allow bool) {
	for _, role := range roles {
		for _, userRole := range User.Roles {
			if role == userRole {
				return true
			}
		}
	}
	return
}

func (User *User) GetRoles() []string {
	return User.Roles
}

func (User User) ToBytes() []byte {
	data, _ := json.Marshal(User)
	return data
}
