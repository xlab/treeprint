package treeprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type nameStruct struct {
	One   string `json:"one" tree:"one"`
	Two   int    `tree:"two"`
	Three struct {
		SubOne   []string
		SubTwo   []interface{}
		SubThree struct {
			InnerOne   *float64  `tree:"inner_one,omitempty"`
			InnerTwo   *struct{} `tree:",omitempty"`
			InnerThree *float64  `tree:"inner_three"`
		}
	}
}

func TestFromStructName(t *testing.T) {
	assert := assert.New(t)

	tree, err := FromStruct(nameStruct{}, StructNameTree)
	assert.NoError(err)

	actual := tree.String()
	expected := `.
├── one
├── two
└── Three
    ├── SubOne
    ├── SubTwo
    └── SubThree
        └── inner_three
`
	assert.Equal(expected, actual)
}

type valueStruct struct {
	Name string
	Bio  struct {
		Age  int
		City string
		Meta interface{}
	}
}

func TestFromStructValue(t *testing.T) {
	assert := assert.New(t)

	val := valueStruct{
		Name: "Max",
	}
	val.Bio.Age = 100
	val.Bio.City = "NYC"
	val.Bio.Meta = []byte("hello")
	tree, err := FromStruct(val, StructValueTree)
	assert.NoError(err)

	actual := tree.String()
	expected := `.
├── [Max]  Name
└── Bio
    ├── [100]  Age
    ├── [NYC]  City
    └── [[104 101 108 108 111]]  Meta
`
	assert.Equal(expected, actual)
}
