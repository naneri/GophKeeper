package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hash(key, password string) string {
	h := hmac.New(sha256.New, []byte(key))

	h.Write([]byte(password))

	dst := h.Sum(nil)

	return hex.EncodeToString(dst)
}

func CheckHash(key, stringToCheck, storedString string) (bool, error) {
	h := hmac.New(sha256.New, []byte(key))

	h.Write([]byte(stringToCheck))

	dst := h.Sum(nil)

	storedStringToByte, err := hex.DecodeString(storedString)

	isEqual := hmac.Equal(dst, storedStringToByte)

	return isEqual, err
}
