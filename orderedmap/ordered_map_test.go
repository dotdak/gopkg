package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap(t *testing.T) {
	art := assert.New(t)

	orderedmap := NewOrderedMap[int, int]()
	expected := []int{1, 5, 2, 9}
	for _, i := range expected {
		orderedmap.Add(i, 1)
	}

	art.Equal(expected, orderedmap.Keys())

	for _, i := range expected {
		k, _ := orderedmap.PopFirst()
		art.Equal(k, i)
	}

	art.True(orderedmap.IsEmpty())
}
