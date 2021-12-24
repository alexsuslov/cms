package model

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Invite struct {
	UUID    uuid.UUID `json:"uuid"`
	Email   string    `json:"email"`
	Expired time.Time `json:"expired"`
	Count   int       `json:"count"`
}

func (Invite Invite) IsExpired() bool {
	return Invite.Expired.Before(time.Now())
}

func (Invite Invite) GetUUID() []byte {
	return Invite.UUID.Bytes()
}

func (Invite Invite) GetEmail() string {
	return Invite.Email
}

func (Invite Invite) ToBytes() []byte {
	data, _ := json.Marshal(Invite)
	return data
}
