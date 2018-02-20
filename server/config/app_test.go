package config

import (
	"testing"
)

func TestCreateApp(t *testing.T) {
	app := CreateApp("root", "secret", "public", []string{"http://localhost:4400"}, []*App{
		CreateApp("child-app", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
	})

	if app.ID != "root" {
		t.Error("App.ID wasn't set correctly")
	}

	if app.Secret != "secret" {
		t.Error("App.Secret wasn't set correctly")
	}

	if app.Public != "public" {
		t.Error("App.Public wasn't set correctly")
	}

	if app.Origins[0] != "http://localhost:4400" {
		t.Error("App.Origins wasn't set correctly")
	}

	if app.Children[0].ID != "child-app" {
		t.Error("App.Children wasn't set correctly")
	}

	if app.Children[0].Parent != app {
		t.Error("App.Parent wasn't set correctly")
	}
}

func TestAppsAppend(t *testing.T) {
	apps := createAppMap()
	apps.Append(
		CreateApp("root", "secret", "public", []string{"http://localhost:4400"}, []*App{
			CreateApp("child-app", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
		}),
	)

	if _, ok := apps.configs["root"]; !ok {
		t.Error("Apps.Append didn't append root app")
	}

	if _, ok := apps.configs["child-app"]; !ok {
		t.Error("Apps.Append didn't handle child apps")
	}

	if apps.Append(CreateApp("root", "secret", "public", []string{"http://localhost:4400"}, nil)) != nil {
		t.Error("Apps.Append didn't prevent id colisions")
	}
}

func TestAppsFind(t *testing.T) {
	apps := createAppMap()
	child := CreateApp("child-app", "child-secret", "child-public", []string{"http://child-localhost"}, nil)
	root := CreateApp("root", "secret", "public", []string{"http://localhost:4400"}, []*App{child})
	apps.Append(root)

	if apps.Find("child-app") != child {
		t.Error("Apps.Find didn't return the child app instance")
	}

	if apps.Find("root") != root {
		t.Error("Apps.Find didn't return the root app instance")
	}
}

func TestAppsLength(t *testing.T) {
	apps := createAppMap()
	apps.Append(
		CreateApp("root", "secret", "public", []string{"http://localhost:4400"}, []*App{
			CreateApp("child-app", "child-secret", "child-public", []string{"http://child-localhost"}, nil),
		}),
	)

	if apps.Length() != 2 {
		t.Error("Apps.Length didn't return length of the map")
	}
}

func TestIsChildOf(t *testing.T) {
	apps := createAppMap()
	apps.Append(CreateApp("root", "secret", "public", []string{"http://localhost:4400"}, nil))
	apps.Append(
		CreateApp("parent", "secret", "public", []string{"http://localhost:4400"}, []*App{
			CreateApp("child", "child", "child", []string{"http://child"}, []*App{
				CreateApp("grand-child", "grand-child", "grand-child", []string{"http://grand-child"}, nil),
			}),
		}),
	)

	root := apps.Find("root")
	parent := apps.Find("parent")
	grandChild := apps.Find("grand-child")

	if grandChild.ChildOf(root) != false {
		t.Error("App.ChildOf Shouldn't be child of root")
	}

	if root.ChildOf(parent) != false {
		t.Error("App.ChildOf shouldn't have any children")
	}

	if grandChild.ChildOf(parent) != true {
		t.Error("App.ChildOf Should be child of parent")
	}

}
