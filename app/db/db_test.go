package db

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/util"
)

func TestInsert(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	user := domain.User{
		ID:                      util.GetID(),
		Email:                   "test@test.com",
		Username:                "testing",
		Password:                util.HashPassword([]byte("testing"), []byte("test")),
		RoleName:                "testing",
		EmailVerified:           false,
		EmailVerificationString: util.ShaHashString("testing"),
		EmailVerifiedDate:       now,
		CreateDate:              now,
		ModifyDate:              now,
	}
	ok := SaveUser(user)
	log.Println("save ok", ok)
	u := GetUser("test@test.com")
	assert.Equal(t, "test@test.com", u.Email, "test@test.com isn't equal to "+u.Email)
	log.Println("removing user test@test.com")
	err := exec(`delete from user where email == "test@test.com"`)
	if err != nil {
		log.Println("delete error", err)
	}
	After()
}
