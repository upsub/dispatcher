package auth

import (
	"encoding/gob"
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
		conf:  conf,
		auths: map[string]*Auth{},
	}

	store.load()

	if store.Length() > 0 {
		return store
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
		Rules   struct {
			Create bool
			Update bool
			Delete bool
		}
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
		&rules{
			decoded.Rules.Create,
			decoded.Rules.Update,
			decoded.Rules.Delete,
		},
	}
}

// CreateFromMessage creates a new auth and authends it to the Store map
func (store *Store) CreateFromMessage(msg *message.Message, parentID string) bool {
	auth := store.decode(msg)

	if auth == nil {
		return false
	}

	auth.Parent = store.Find(parentID)

	if !auth.Parent.HasChild(auth) {
		auth.Parent.Children = append(auth.Parent.Children, auth)
	}

	store.Append(auth)
	store.save()

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

	store.save()

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

	store.save()

	return true
}

func (store *Store) load() {
	if _, err := os.Stat(store.conf.AuthDataPath); os.IsNotExist(err) {
		return
	}

	file, err := os.Open(store.conf.AuthDataPath)

	if err != nil {
		log.Print("[ERROR] ", err)
		return
	}

	decoded := map[string]serializedAuth{}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&decoded)

	if err != nil {
		log.Print("[ERROR] ", err)
		return
	}

	for id := range decoded {
		store.auths[id] = nil
	}

	for id, serialized := range decoded {
		auth := CreateAuth(
			serialized.ID,
			serialized.Secret,
			serialized.Public,
			serialized.Origins,
			nil,
		)

		store.auths[id] = auth

		auth.SetRules(
			serialized.Rules.Create,
			serialized.Rules.Update,
			serialized.Rules.Delete,
		)
	}

	for id, serialized := range decoded {
		auth := store.auths[id]
		if serialized.Parent != auth.ID {
			auth.Parent = store.Find(serialized.Parent)
		}

		for _, serializedChild := range serialized.Children {
			auth.Children = append(auth.Children, store.Find(serializedChild))
		}
	}

	file.Close()
}

func (store *Store) save() {
	file, err := os.Create(store.conf.AuthDataPath)

	if err != nil {
		log.Print("[ERROR] ", err)
		return
	}

	serialized := map[string]serializedAuth{}

	for id, auth := range store.auths {
		serialized[id] = auth.serialize()
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(serialized)
	file.Close()
}
