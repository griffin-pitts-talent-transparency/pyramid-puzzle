package seed

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
)

func HMACSeed(key []byte, message string) int64 {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sum := h.Sum(nil)
	return int64(binary.BigEndian.Uint64(sum[:8]))
}
