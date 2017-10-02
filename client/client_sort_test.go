package client

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByNameAZ(t *testing.T) {
	got := []Client{
		{
			Name:  "A",
			Owner: "C",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "C",
			Owner: "A",
		},
	}
	expected := []Client{
		{
			Name:  "A",
			Owner: "C",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "C",
			Owner: "A",
		},
	}
	sort.Sort(ByName(got))
	assert.Equal(t, expected, got)
}

func TestByNameZA(t *testing.T) {
	got := []Client{
		{
			Name:  "C",
			Owner: "A",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "A",
			Owner: "C",
		},
	}
	expected := []Client{
		{
			Name:  "A",
			Owner: "C",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "C",
			Owner: "A",
		},
	}
	sort.Sort(ByName(got))
	assert.Equal(t, expected, got)
}

func TestByOwnerAZ(t *testing.T) {
	got := []Client{
		{
			Name:  "C",
			Owner: "A",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "A",
			Owner: "C",
		},
	}
	expected := []Client{
		{
			Name:  "C",
			Owner: "A",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "A",
			Owner: "C",
		},
	}
	sort.Sort(ByOwner(got))
	assert.Equal(t, expected, got)
}

func TestByOwnerZA(t *testing.T) {
	got := []Client{
		{
			Name:  "A",
			Owner: "C",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "C",
			Owner: "A",
		},
	}
	expected := []Client{
		{
			Name:  "C",
			Owner: "A",
		},
		{
			Name:  "B",
			Owner: "B",
		},
		{
			Name:  "A",
			Owner: "C",
		},
	}
	sort.Sort(ByOwner(got))
	assert.Equal(t, expected, got)
}
