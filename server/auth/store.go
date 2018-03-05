package auth

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/message"
)

// Store contains a map of Store
type Store struct {
	conf  *config.Config
	auths map[string]*Auth
}

func NewStore(conf *config.Config) *Store {
	if conf == nil {
		conf = config.Create()
	}

	store := &Store{
		auths: map[string]*Auth{},
	}

	root := CreateAuth(
		os.Getenv("AUTH_APP_ID"),
		os.Getenv("AUTH_SECRET"),
		os.Getenv("AUTH_PUBLIC"),
		strings.Split(os.Getenv("AUTH_ORIGINS"), ","),
		nil,
	)

	if root.ID != "" || root.Secret != "" || root.Public != "" {
		store.Append(root)
	}

	if root.ID != "" || root.Secret != "" {
		root.SetRules(true, true, true)
	}

	return store
}

// Append a new Auth
func (store *Store) Append(auth *Auth) *Store {
	if _, ok := store.auths[auth.ID]; ok {
		log.Print("[ERROR] Auth is already created: " + auth.ID)
		return nil
	}

	store.auths[auth.ID] = auth

	for _, child := range auth.Children {
		store.Append(child)
	}

	return store
}

// Find auth from id
func (store *Store) Find(id string) *Auth {
	return store.auths[id]
}

// Length of the Store map
func (store *Store) Length() int {
	return len(store.auths)
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

func (store *Store) decode(msg *message.Message) *Auth {
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
		store.Find(decoded.Parent),
		nil,
		decoded.rules,
	}
}

// CreateFromMessage creates a new auth and authends it to the Store map
func (store *Store) CreateFromMessage(msg *message.Message, parentID string) bool {
	auth := store.decode(msg)

	if auth == nil {
		return false
	}

	store.Append(auth)

	auth.Parent = store.Find(parentID)
	auth.Parent.Children = append(auth.Parent.Children, auth)

	return true
}

func (store *Store) UpdateFromMessage(msg *message.Message) bool {
	newAuth := store.decode(msg)

	if newAuth == nil {
		return false
	}

	if _, ok := store.auths[newAuth.ID]; !ok {
		log.Print("[ERROR] Couldn't update auth " + newAuth.ID + " because it doesn't exist.")
		return false
	}

	oldAuth := store.Find(newAuth.ID)
	oldAuth.Secret = newAuth.Secret
	oldAuth.Public = newAuth.Public
	oldAuth.Origins = newAuth.Origins
	oldAuth.rules = newAuth.rules

	return true
}

func (store *Store) DeleteFromMessage(msg *message.Message) bool {
	id := strings.Replace(msg.Payload, "\"", "", 2)
	return store.Remove(id)
}

func (store *Store) Remove(id string) bool {
	auth := store.Find(id)

	if auth == nil {
		log.Print("[ERROR] Couldn't delete the " + id + " because it doesn't exist.")
		return false
	}

	for _, child := range auth.Children {
		store.Remove(child.ID)
	}

	if auth.Parent != nil {
		auth.Parent.RemoveChild(auth)
	}

	delete(store.auths, auth.ID)

	return true
}
