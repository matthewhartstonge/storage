package user

// ByUsername enables sorting user accounts by Username A-Z
type ByUsername []User

func (u ByUsername) Len() int {
	return len(u)
}

func (u ByUsername) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u ByUsername) Less(i, j int) bool {
	return u[i].Username < u[j].Username
}

// ByFirstName enables sorting user accounts by First Name A-Z
type ByFirstName []User

func (u ByFirstName) Len() int {
	return len(u)
}

func (u ByFirstName) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u ByFirstName) Less(i, j int) bool {
	return u[i].FirstName < u[j].FirstName
}

// ByLastName enables sorting user accounts by Last Name A-Z
type ByLastName []User

func (u ByLastName) Len() int {
	return len(u)
}

func (u ByLastName) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u ByLastName) Less(i, j int) bool {
	return u[i].LastName < u[j].LastName
}
