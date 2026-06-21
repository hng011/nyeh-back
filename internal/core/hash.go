package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(rawToken string) string {
	h := hmac.New(sha256.New, []byte(Settings.TOKEN_DIGEST))
	h.Write([]byte(rawToken))
	return hex.EncodeToString(h.Sum(nil))
}
