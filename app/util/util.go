package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"regexp"

	"github.com/ventu-io/go-shortid"
	"golang.org/x/crypto/pbkdf2"
)

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

var nonce []byte

var key = []byte("AES256Key-32Characters1234567890")

func Encrypt(text string) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	nonce = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(aesgcm.Seal(nil, nonce, []byte(text), nil))
}

func Decrypt(text string) string {
	ciphertext, _ := hex.DecodeString(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err.Error())
	}
	return string(plaintext)
}
