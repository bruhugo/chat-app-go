package services

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type hashService struct {
	hash hash.Hash
}

func NewHashService() *hashService {
	return &hashService{
		hash: sha256.New(),
	}
}

func (hs *hashService) Hash(str string) string {
	hs.hash.Write([]byte(str))
	bs := hs.hash.Sum(nil)

	return base64.RawStdEncoding.EncodeToString(bs)
}
