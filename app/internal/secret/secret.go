package secret

import (
	"errors"
	"os"
)

func LoadPuzzleSecret() (string, error) {
	s := os.Getenv("PYRAMID_SECRET")
	if s == "" {
		return "", errors.New("PYRAMID_SECRET is not set")
	}
	return s, nil
}
