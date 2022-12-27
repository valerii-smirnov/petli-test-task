package hasher

import (
	"crypto/md5"
	"encoding/hex"
)

type MD5 struct {
	salt string
}

func NewMD5(salt string) *MD5 {
	return &MD5{
		salt: salt,
	}
}

func (h MD5) StringHash(inp string) string {
	hash := md5.Sum([]byte(inp + h.salt))
	return hex.EncodeToString(hash[:])
}
