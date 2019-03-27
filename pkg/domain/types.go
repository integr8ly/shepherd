package domain

type User struct {
	Admin    bool
	Org      string
	Name     string
	ID       string
	Team     string
	Role     string // admin | member | general
	Timezone string
}

func (u User) IsAdmin() bool {
	return u.Admin
}
