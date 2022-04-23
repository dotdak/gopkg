package pcontext

import "fmt"

func NewContextKey(key string) contextKey {
	return contextKey(key)
}

type contextKey string

func (c contextKey) String() string {
	return fmt.Sprint("Pcontext ", string(c))
}
