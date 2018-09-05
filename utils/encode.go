package utils

import (
	"crypto/md5"
	"fmt"
	"encoding/hex"
	"io"
	"encoding/base64"
	"crypto/rand"
)

func GetMD5(content string) string {
	hashContent := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hashContent)
}

//生成Guid字串
func GetUUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}

//生成32位md5字串
func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
