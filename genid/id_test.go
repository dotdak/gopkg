package genid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	art := assert.New(t)
	id := NewID()
	art.Equal(len(id), 16)
}
