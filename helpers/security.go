package helpers

import (
	"time"
	"crypto/sha512"
	"encoding/hex"
)


func GenerateTokenRecovery() string {
	today := time.Now()
	k := today.String() + "Presto"
	h := sha512.New()
	h.Write([]byte(k))
	return hex.EncodeToString(h.Sum(nil))
}