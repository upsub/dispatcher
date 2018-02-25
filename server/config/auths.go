package config

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/upsub/dispatcher/server/message"
)

// Auths contains a map of Auths
type Auths struct {
	configs map[string]*Auth
}

// Append a new Auth
func (auths *Auths) Append(auth *Auth) *Auths {
	if _, ok := auths.configs[auth.ID]; ok {
		log.Print("[ERROR] Auth is already created: " + auth.ID)
		return nil
	}

	auths.configs[auth.ID] = auth

	for _, child := range auth.Children {
		auths.Append(child)
	}

	return auths
}

// Find auth from id
func (auths *Auths) Find(id string) *Auth {
	return auths.configs[id]
}

// Length of the Auths map
func (auths *Auths) Length() int {
	return len(auths.configs)
}

// ChildOf checks if an auth is a child of another auth
func (child *Auth) ChildOf(parent *Auth) bool {
	if child.Parent == parent {
		return true
	}

	if child.Parent == nil {
		return false
	}

	return child.Parent.ChildOf(parent)
}

func createAuths() *Auths {
	return &Auths{
		configs: map[string]*Auth{},
	}
}

func (auths *Auths) decode(msg *message.Message) *Auth {
	var decoded struct {
		ID      string
		Secret  string
		Public  string
		Origins []string
		Parent  string
		rules   *rules
	}

	err := json.Unmarshal([]byte(msg.Payload), &decoded)

	if err != nil {
		log.Print("[ERROR] ", err)
		return nil
	}

	return &Auth{
		decoded.ID,
		decoded.Secret,
		decoded.Public,
		decoded.Origins,
		auths.Find(decoded.Parent),
		nil,
		decoded.rules,
	}
}

// CreateFromMessage creates a new auth and authends it to the auths map
func (auths *Auths) CreateFromMessage(msg *message.Message, parentID string) bool {
	auth := auths.decode(msg)

	if auth == nil {
		return false
	}

	auths.Append(auth)

	auth.Parent = auths.Find(parentID)
	auth.Parent.Children = append(auth.Parent.Children, auth)

	return true
}

func (auths *Auths) UpdateFromMessage(msg *message.Message) bool {
	newAuth := auths.decode(msg)

	if newAuth == nil {
		return false
	}

	if _, ok := auths.configs[newAuth.ID]; !ok {
		log.Print("[ERROR] Couldn't update auth " + newAuth.ID + " because it doesn't exist.")
		return false
	}

	oldAuth := auths.Find(newAuth.ID)
	oldAuth.Secret = newAuth.Secret
	oldAuth.Public = newAuth.Public
	oldAuth.Origins = newAuth.Origins
	oldAuth.rules = newAuth.rules

	return true
}

func (auths *Auths) DeleteFromMessage(msg *message.Message) bool {
	id := strings.Replace(msg.Payload, "\"", "", 2)
	return auths.Remove(id)
}

func (auths *Auths) Remove(id string) bool {
	auth := auths.Find(id)

	if auth == nil {
		log.Print("[ERROR] Couldn't delete the " + id + " because it doesn't exist.")
		return false
	}

	for _, child := range auth.Children {
		auths.Remove(child.ID)
	}

	if auth.Parent != nil {
		auth.Parent.RemoveChild(auth)
	}

	delete(auths.configs, auth.ID)

	return true
}
