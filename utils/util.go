package utils

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	mrand "math/rand"
	"time"
)

func init() {
	mrand.Seed(time.Now().UTC().UnixNano())
}

//IndexOfStr returns index of string in slice if exists else -1
func IndexOfStr(slice []string, needle string) int {
	for p, v := range slice {
		if v == needle {
			return p
		}
	}
	return -1
}

//Hashing password bcrypt algorithm with default cost 10

func HashPass(in string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	return string(pass), err
}

func CutStr(s string, n int) string {
	if i := len(s); i >= n {
		return s[:i-n]
	}
	return ""
}

func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func RandInt(min, max int) int {
	return min + mrand.Intn(max-min)
}
