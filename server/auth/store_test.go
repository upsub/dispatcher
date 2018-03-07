package auth

import (
	"testing"

	"github.com/upsub/dispatcher/server/message"
)

func TestStoreAppend(t *testing.T) {
	store := NewStore(nil)
	store.Append(
		CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
			CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
		}),
	)

	if _, ok := store.auths["root"]; !ok {
		t.Error("store.Append didn't authend root auth")
	}

	if _, ok := store.auths["child-auth"]; !ok {
		t.Error("store.Append didn't handle child Store")
	}

	if store.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, nil)) != nil {
		t.Error("store.Append didn't prevent id colisions")
	}
}

func TestStoreFind(t *testing.T) {
	store := NewStore(nil)
	child := CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil)
	root := CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{child})
	store.Append(root)

	if store.Find("child-auth") != child {
		t.Error("store.Find didn't return the child auth instance")
	}

	if store.Find("root") != root {
		t.Error("store.Find didn't return the root auth instance")
	}
}

func TestStoreLength(t *testing.T) {
	store := NewStore(nil)
	store.Append(
		CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
			CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
		}),
	)

	if store.Length() != 2 {
		t.Error("store.Length didn't return length of the map", store.Length())
	}
}

func TestIsChildOf(t *testing.T) {
	store := NewStore(nil)
	store.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, nil))
	store.Append(
		CreateAuth("parent", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
			CreateAuth("child", "child", "child", []string{"http://child"}, []*Auth{
				CreateAuth("grand-child", "grand-child", "grand-child", []string{"http://grand-child"}, nil),
			}),
		}),
	)

	root := store.Find("root")
	parent := store.Find("parent")
	grandChild := store.Find("grand-child")

	if grandChild.ChildOf(root) != false {
		t.Error("Auth.ChildOf Shouldn't be child of root")
	}

	if root.ChildOf(parent) != false {
		t.Error("Auth.ChildOf shouldn't have any children")
	}

	if grandChild.ChildOf(parent) != true {
		t.Error("Auth.ChildOf Should be child of parent")
	}

}

func helperCreateAuthFromMessage(store *Store) *Auth {
	msg := message.Create("{\"id\":\"upsub\",\"secret\":\"upsub-secret\",\"public\":\"upsub-public\",\"origins\":[\"http://localhost:3000\"],\"rules\":{\"create\":false,\"update\":false,\"delete\":false}}")
	return store.decode(msg)
}

func TestDecode(t *testing.T) {
	store := NewStore(nil)
	store.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:3000"}, nil))
	auth := helperCreateAuthFromMessage(store)

	if auth.ID != "upsub" {
		t.Error("Auth.ID wasn't correctly")
	}

	if auth.Secret != "upsub-secret" {
		t.Error("Auth.Secret wasn't correctly")
	}

	if auth.Public != "upsub-public" {
		t.Error("Auth.Public wasn't correctly")
	}

	if auth.Origins[0] != "http://localhost:3000" {
		t.Error("Auth.Public wasn't correctly")
	}
}

func TestDecodeWithInvalidJSON(t *testing.T) {
	store := NewStore(nil)
	store.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, nil))
	msg := message.Create("{\"id\"\"upsub\",\"secret\":\"upsub-secret\",\"public\":\"upsub-public\",\"origins\":[\"http://localhost:3000\"],\"parent\":\"root\",\"rules\":{\"create\":false,\"update\":false,\"delete\":false}}")
	auth := store.decode(msg)

	if auth != nil {
		t.Error("Should return nil when receiving invalid json")
	}
}
