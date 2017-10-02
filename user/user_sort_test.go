package user

import (
	// Standard Library Imports
	"sort"
	"testing"

	// External Imports
	"github.com/stretchr/testify/assert"
)

func TestByUsernameAZ(t *testing.T) {
	got := []User{
		{
			Username: "A",
		},
		{
			Username: "B",
		},
		{
			Username: "C",
		},
	}
	expected := []User{
		{
			Username: "A",
		},
		{
			Username: "B",
		},
		{
			Username: "C",
		},
	}
	sort.Sort(ByUsername(got))
	assert.Equal(t, expected, got)
}

func TestByUsernameZA(t *testing.T) {
	got := []User{
		{
			Username: "C",
		},
		{
			Username: "B",
		},
		{
			Username: "A",
		},
	}
	expected := []User{
		{
			Username: "A",
		},
		{
			Username: "B",
		},
		{
			Username: "C",
		},
	}
	sort.Sort(ByUsername(got))
	assert.Equal(t, expected, got)
}

func TestByFirstNameAZ(t *testing.T) {
	got := []User{
		{
			Username: "X",
		},
		{
			Username: "Y",
		},
		{
			Username: "Z",
		},
	}
	expected := []User{
		{
			Username: "X",
		},
		{
			Username: "Y",
		},
		{
			Username: "Z",
		},
	}
	sort.Sort(ByFirstName(got))
	assert.Equal(t, expected, got)
}

func TestByFirstNameZA(t *testing.T) {
	got := []User{
		{
			FirstName: "Z",
		},
		{
			FirstName: "Y",
		},
		{
			FirstName: "X",
		},
	}
	expected := []User{
		{
			FirstName: "X",
		},
		{
			FirstName: "Y",
		},
		{
			FirstName: "Z",
		},
	}
	sort.Sort(ByFirstName(got))
	assert.Equal(t, expected, got)
}

func TestByLastNameAZ(t *testing.T) {
	got := []User{
		{
			LastName: "A",
		},
		{
			LastName: "G",
		},
		{
			LastName: "Z",
		},
	}
	expected := []User{
		{
			LastName: "A",
		},
		{
			LastName: "G",
		},
		{
			LastName: "Z",
		},
	}
	sort.Sort(ByLastName(got))
	assert.Equal(t, expected, got)
}

func TestByLastNameZA(t *testing.T) {
	got := []User{
		{
			LastName: "Z",
		},
		{
			LastName: "G",
		},
		{
			LastName: "A",
		},
	}
	expected := []User{
		{
			LastName: "A",
		},
		{
			LastName: "G",
		},
		{
			LastName: "Z",
		},
	}
	sort.Sort(ByLastName(got))
	assert.Equal(t, expected, got)
}
