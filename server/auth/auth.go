package auth

// Rules describes what a auth os allowed to do
type rules struct {
	create bool
	update bool
	delete bool
}

// Auth handles authentication configuration
type Auth struct {
	ID       string
	Secret   string
	Public   string
	Origins  []string
	Parent   *Auth
	Children []*Auth
	rules    *rules
}

// serializedAuth
type serializedAuth struct {
	ID       string
	Secret   string
	Public   string
	Origins  []string
	Parent   string
	Children []string
	Rules    struct {
		Create bool
		Update bool
		Delete bool
	}
}

// SetRules set the rules of what the auth is allowed to do
func (auth *Auth) SetRules(create bool, update bool, delete bool) {
	auth.rules.create = create
	auth.rules.update = update
	auth.rules.delete = delete
}

// CanCreate can a auth create a child auth
func (auth *Auth) CanCreate() bool {
	return auth.rules.create
}

// CanUpdate can a auth update a child auth
func (auth *Auth) CanUpdate() bool {
	return auth.rules.update
}

// CanDelete can a auth delete a child auth
func (auth *Auth) CanDelete() bool {
	return auth.rules.delete
}

// HasChild
func (parent *Auth) HasChild(child *Auth) bool {
	for _, c := range parent.Children {
		if c.ID == child.ID {
			return true
		}
	}

	return false
}

// RemoveChild remove a child from the auth
func (auth *Auth) RemoveChild(child *Auth) {
	newChildren := []*Auth{}

	for _, c := range auth.Children {
		if c != child {
			newChildren = append(newChildren, c)
		}
	}

	auth.Children = newChildren
}

func (auth *Auth) serialize() serializedAuth {
	serialized := serializedAuth{
		auth.ID,
		auth.Secret,
		auth.Public,
		auth.Origins,
		"",
		[]string{},
		struct {
			Create bool
			Update bool
			Delete bool
		}{
			auth.rules.create,
			auth.rules.update,
			auth.rules.delete,
		},
	}

	if auth.Parent != nil {
		serialized.Parent = auth.Parent.ID
	}

	for _, child := range auth.Children {
		serialized.Children = append(serialized.Children, child.ID)
	}

	return serialized
}

// CreateAuth creates a new app config
func CreateAuth(
	id string,
	secret string,
	public string,
	origins []string,
	children []*Auth,
) *Auth {
	if children == nil {
		children = []*Auth{}
	}

	auth := &Auth{
		id,
		secret,
		public,
		origins,
		nil,
		children,
		&rules{false, false, false},
	}

	for _, child := range children {
		child.Parent = auth
	}

	return auth
}
