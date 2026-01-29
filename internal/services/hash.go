package services

import "crypto/sha256"

func Hash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	bs := hash.Sum(nil)

	return string(bs)
}
