package config

import "log"

// App handles authentication configuration
type App struct {
	ID       string
	Secret   string
	Public   string
	Origins  []string
	Parent   *App
	Children []*App
}

// Apps contains a map of Apps
type Apps struct {
	configs map[string]*App
}

// Append a new App
func (apps *Apps) Append(app *App) *Apps {
	if _, ok := apps.configs[app.ID]; ok {
		log.Print("[ERROR] App is already created: " + app.ID)
		return nil
	}

	apps.configs[app.ID] = app

	for _, child := range app.Children {
		apps.Append(child)
	}

	return apps
}

func CreateApp(
	id string,
	secret string,
	public string,
	origins []string,
	children []*App,
) *App {
	if children == nil {
		children = []*App{}
	}

	app := &App{id, secret, public, origins, nil, children}

	for _, child := range children {
		child.Parent = app
	}

	return app
}

// Find app from id
func (apps *Apps) Find(id string) *App {
	return apps.configs[id]
}

func (apps *Apps) Length() int {
	return len(apps.configs)
}

func createAppMap() *Apps {
	return &Apps{
		configs: map[string]*App{},
	}
}
