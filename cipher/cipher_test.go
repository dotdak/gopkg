package cipher

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	art := assert.New(t)
	key := make([]byte, 32)
	_, e := rand.Read(key)
	if e != nil {
		t.Error(e)
	}
	data := []byte("data_test")
	enc, err := Encrypt(key, data)
	art.Nil(err)
	art.Greater(len(enc), 0)
	dec, err := Decrypt(key, enc)
	art.Nil(err)
	art.Equal(data, dec)
}
