package storage

// UsersByUsername enables sorting user accounts by Username A-Z
type UsersByUsername []User

func (u UsersByUsername) Len() int {
	return len(u)
}

func (u UsersByUsername) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsersByUsername) Less(i, j int) bool {
	return u[i].Username < u[j].Username
}

// UsersByFirstName enables sorting user accounts by First Name A-Z
type UsersByFirstName []User

func (u UsersByFirstName) Len() int {
	return len(u)
}

func (u UsersByFirstName) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsersByFirstName) Less(i, j int) bool {
	return u[i].FirstName < u[j].FirstName
}

// UsersByLastName enables sorting user accounts by Last Name A-Z
type UsersByLastName []User

func (u UsersByLastName) Len() int {
	return len(u)
}

func (u UsersByLastName) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsersByLastName) Less(i, j int) bool {
	return u[i].LastName < u[j].LastName
}
