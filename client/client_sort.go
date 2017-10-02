package client

// ByName enables sorting Client applications by the client application Name A-Z
type ByName []Client

func (c ByName) Len() int {
	return len(c)
}

func (c ByName) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByName) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

// ByOwner enables sorting Client applications by the client application Owner A-Z
type ByOwner []Client

func (c ByOwner) Len() int {
	return len(c)
}

func (c ByOwner) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByOwner) Less(i, j int) bool {
	return c[i].Owner < c[j].Owner
}
