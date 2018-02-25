package config

import (
	"testing"

	"github.com/upsub/dispatcher/server/message"
)

func TestAuthsAppend(t *testing.T) {
	auths := createAuths()
	auths.Append(
		CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
			CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
		}),
	)

	if _, ok := auths.configs["root"]; !ok {
		t.Error("Auths.Append didn't authend root auth")
	}

	if _, ok := auths.configs["child-auth"]; !ok {
		t.Error("Auths.Append didn't handle child auths")
	}

	if auths.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, nil)) != nil {
		t.Error("Auths.Append didn't prevent id colisions")
	}
}

func TestAuthsFind(t *testing.T) {
	auths := createAuths()
	child := CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil)
	root := CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{child})
	auths.Append(root)

	if auths.Find("child-auth") != child {
		t.Error("Auths.Find didn't return the child auth instance")
	}

	if auths.Find("root") != root {
		t.Error("Auths.Find didn't return the root auth instance")
	}
}

func TestAuthsLength(t *testing.T) {
	auths := createAuths()
	auths.Append(
		CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
			CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
		}),
	)

	if auths.Length() != 2 {
		t.Error("Auths.Length didn't return length of the map")
	}
}

func TestIsChildOf(t *testing.T) {
	auths := createAuths()
	auths.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, nil))
	auths.Append(
		CreateAuth("parent", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
			CreateAuth("child", "child", "child", []string{"http://child"}, []*Auth{
				CreateAuth("grand-child", "grand-child", "grand-child", []string{"http://grand-child"}, nil),
			}),
		}),
	)

	root := auths.Find("root")
	parent := auths.Find("parent")
	grandChild := auths.Find("grand-child")

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

func helperCreateAuthFromMessage(auths *Auths) *Auth {
	msg := message.Create("{\"id\":\"upsub\",\"secret\":\"upsub-secret\",\"public\":\"upsub-public\",\"origins\":[\"http://localhost:3000\"],\"rules\":{\"create\":false,\"update\":false,\"delete\":false}}")
	return auths.decode(msg)
}

func TestDecode(t *testing.T) {
	auths := createAuths()
	auths.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:3000"}, nil))
	auth := helperCreateAuthFromMessage(auths)

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
	auths := createAuths()
	auths.Append(CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, nil))
	msg := message.Create("{\"id\"\"upsub\",\"secret\":\"upsub-secret\",\"public\":\"upsub-public\",\"origins\":[\"http://localhost:3000\"],\"parent\":\"root\",\"rules\":{\"create\":false,\"update\":false,\"delete\":false}}")
	auth := auths.decode(msg)

	if auth != nil {
		t.Error("Should return nil when receiving invalid json")
	}
}
