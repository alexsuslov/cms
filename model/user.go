package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
)

type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Token    []byte   `json:"token"`
}

func (User *User) GenToken(pass string) []byte {
	token := sha256.Sum256([]byte(User.ID + ":" + pass))
	token1 := sha256.Sum256(token[:])
	return token1[:]
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
