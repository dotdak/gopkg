package env

import (
	"os"
	"strconv"
)

func Coalesce[K comparable](vals ...K) K {
	var zero K
	for _, v := range vals {
		if v != zero {
			return v
		}
	}

	return zero
}

func EnvString(key, defaultValue string) string {
	return Coalesce(os.Getenv(key), defaultValue)
}

func EnvInt(key string, defaultValue int) int {
	raw := os.Getenv(key)
	val, _ := strconv.Atoi(raw)
	return Coalesce(val, defaultValue)
}

func Has(key string) bool {
	return os.Getenv(key) == "yes"
}
