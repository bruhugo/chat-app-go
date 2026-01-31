package services

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HashService interface {
	Hash(str string) string
}

type ShaH256Service struct {
	hash hash.Hash
}

func NewShaH256Service() *ShaH256Service {
	return &ShaH256Service{
		hash: sha256.New(),
	}
}

func (hs *ShaH256Service) Hash(str string) string {
	hs.hash.Write([]byte(str))
	bs := hs.hash.Sum(nil)

	return base64.RawStdEncoding.EncodeToString(bs)
}
