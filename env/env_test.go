package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoalesce(t *testing.T) {
	art := assert.New(t)
	art.Equal("abc", Coalesce("", "abc"))
	art.Equal("abc", Coalesce("abc", ""))
	art.Equal(1, Coalesce(1, 0))
	art.Equal(1, Coalesce(0, 1))
	art.Equal(0.1, Coalesce(0.0, 0.0, 0.1))
	art.Equal(0.1, Coalesce(0.1, 0.0, 0.2))
	art.Equal(uint(1), Coalesce(uint(0), uint(1)))
}
