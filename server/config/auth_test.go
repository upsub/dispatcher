package config

import (
	"testing"
)

func TestCreateAuth(t *testing.T) {
	auth := CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
		CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
	})

	if auth.ID != "root" {
		t.Error("Auth.ID wasn't set correctly")
	}

	if auth.Secret != "secret" {
		t.Error("Auth.Secret wasn't set correctly")
	}

	if auth.Public != "public" {
		t.Error("Auth.Public wasn't set correctly")
	}

	if auth.Origins[0] != "http://localhost:4400" {
		t.Error("Auth.Origins wasn't set correctly")
	}

	if auth.Children[0].ID != "child-auth" {
		t.Error("Auth.Children wasn't set correctly")
	}

	if auth.Children[0].Parent != auth {
		t.Error("Auth.Parent wasn't set correctly")
	}

	if auth.rules.create != false || auth.rules.update != false || auth.rules.delete != false {
		t.Error("Auth.rules wans't set correctly")
	}
}

func TestSetRules(t *testing.T) {
	auth := CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
		CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
	})

	auth.SetRules(true, true, true)

	if auth.rules.create != true {
		t.Error("Auth.rules.create wasn't set correctly")
	}

	if auth.rules.update != true {
		t.Error("Auth.rules.update wasn't set correctly")
	}

	if auth.rules.delete != true {
		t.Error("Auth.rules.delete wasn't set correctly")
	}
}

func TestCanCreate(t *testing.T) {
	auth := CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
		CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
	})

	auth.SetRules(true, false, false)

	if auth.CanCreate() != true {
		t.Error("auth.CanCreate didn't return the expected value")
	}
}

func TestCanUpdate(t *testing.T) {
	auth := CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
		CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
	})

	auth.SetRules(false, true, false)

	if auth.CanUpdate() != true {
		t.Error("auth.CanUpdate didn't return the expected value")
	}
}

func TestCanDelete(t *testing.T) {
	auth := CreateAuth("root", "secret", "public", []string{"http://localhost:4400"}, []*Auth{
		CreateAuth("child-auth", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
	})

	auth.SetRules(false, false, true)

	if auth.CanDelete() != true {
		t.Error("auth.CanDelete didn't return the expected value")
	}
}
