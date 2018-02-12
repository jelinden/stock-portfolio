package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/ventu-io/go-shortid"
	"golang.org/x/crypto/pbkdf2"
)

var nonce []byte

var key = []byte("AES256Key-32Characters1234567890")

func init() {
	n, err := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	if err != nil {
		fmt.Println(err.Error())
	}
	nonce = n
}

func HashPassword(password, salt []byte) string {
	return base64.URLEncoding.EncodeToString(pbkdf2.Key(password, salt, 4096, sha512.Size, sha512.New))
}

func ShaHashString(hashable string) string {
	h := sha1.New()
	h.Write([]byte(hashable))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(".+@.+\\..+")
	return re.Match([]byte(email))
}

func GetID() string {
	id, err := shortid.Generate()
	if err != nil {
		log.Println(err)
		return ""
	}
	return id
}

func Encrypt(text string) string {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(text), nil))
}

func Decrypt(text string) string {
	dText, _ := hex.DecodeString(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("aes.NewCipher", err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("cipher.NewGCM", err.Error())
	}

	nonceSize := gcm.NonceSize()
	if len(dText) < nonceSize {
		log.Println("ciphertext too short")
	}

	nonce, ciphertext := dText[:nonceSize], dText[nonceSize:]
	decryptedText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err.Error())
	}
	return string(decryptedText)
}
