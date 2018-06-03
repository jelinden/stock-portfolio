package db

import (
	"log"
	"strings"

	"github.com/jelinden/stock-portfolio/app/domain"
)

func SaveUser(user domain.User) bool {
	err := exec(`INSERT INTO user (id, email, username, password, 
		rolename,emailverified,emailverificationstring,
		emailverifieddate, createdate, modifydate) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`,
		user.ID,
		user.Email,
		user.Username,
		user.Password,
		user.RoleName,
		user.EmailVerified,
		user.EmailVerificationString,
		user.EmailVerifiedDate,
		user.CreateDate,
		user.ModifyDate)
	if err != nil {
		log.Printf("failed with '%s'\n", err.Error())
		return false
	}
	return true
}

func UpdateUser(user domain.User) bool {
	err := exec(`UPDATE user SET 
		email=$2, 
		username=$3, 
		password=$4,
		emailverified=$5,
		emailverifieddate=$6,
		modifydate=$7,
		rolename=$8
		where id = $1;`,
		user.ID,
		user.Email,
		user.Username,
		user.Password,
		user.EmailVerified,
		user.EmailVerifiedDate,
		user.ModifyDate,
		user.RoleName)
	if err != nil {
		log.Printf("failed with '%s'\n", err.Error())
		return false
	}
	return true
}

func GetUser(email string) domain.User {
	return execRow(`SELECT id, email, username, rolename, 
		password, createdate, emailverified,
		emailverifieddate, emailverificationstring,
		modifydate FROM user WHERE email LIKE $1;`, strings.Replace(email, "+", `\+`, -1))
}

func GetUserWithVerifyString(verifyString string) domain.User {
	return queryUser(`SELECT id, email, username, rolename, 
		password, createdate, emailverified,
		emailverifieddate, emailverificationstring,
		modifydate FROM user WHERE emailverificationstring LIKE $1;`, verifyString)
}

func GetUsers() domain.UserList {
	return queryAllUsers()
}
