package genid

import "crypto/rand"

func NewID() string {
	const size = 16
	return NewIDWithSize(size)
}

func NewIDWithSize(size int) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" // base58

	var id = make([]byte, size)
	if _, err := rand.Read(id); err != nil {
		panic(err)
	}
	for i, p := range id {
		id[i] = alphabet[int(p)%len(alphabet)] // discard everything but the least significant bits
	}
	return string(id)
}
