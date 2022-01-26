package routes

import (
	"net/http"
	"regexp"
	"time"

	"github.com/jelinden/stock-portfolio/app/config"
	"github.com/jelinden/stock-portfolio/app/email"
	"github.com/jelinden/stock-portfolio/app/util"

	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/julienschmidt/httprouter"
)

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	checkEmail := checkInput(email)
	checkPassword := checkInput(password)
	if !checkEmail && !checkPassword {
		sessionKey := util.ShaHashString(email)
		user := db.GetUser(email)
		if user.ID != "" && user.Password == util.HashPassword([]byte(password), []byte(email)) {
			db.PutSession(sessionKey, email)
			http.SetCookie(w, &http.Cookie{Name: "login", Value: sessionKey, MaxAge: 2592000})
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
			w.Write(nil)
			return
		}
	}
	w.Header().Add("Location", "/login?login=failed")
	w.WriteHeader(302)
	w.Write(nil)
}

func Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	expiration := time.Now().Add(-365 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{Name: "login", Value: "", Expires: expiration})
	w.Header().Add("Location", "/")
	w.WriteHeader(302)
	w.Write(nil)
	return
}

func Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username := r.FormValue("username")
	emailParam := r.FormValue("email")
	password := r.FormValue("password")
	checkUsername := checkInput(username)
	checkEmail := checkInput(emailParam)
	checkPassword := checkInput(password)
	if checkEmail || !util.ValidateEmail(emailParam) {
		validationFailed(w, "email")
		return
	}
	if checkPassword || len(password) < 8 {
		validationFailed(w, "credentials")
		return
	}
	if checkUsername {
		validationFailed(w, "credentials")
		return
	}
	if db.GetUser(emailParam).ID == "" {
		user := domain.User{
			Username:                username,
			Password:                util.HashPassword([]byte(password), []byte(emailParam)),
			Email:                   emailParam,
			CreateDate:              time.Now().UTC().Format(time.RFC3339),
			ModifyDate:              time.Now().UTC().Format(time.RFC3339),
			EmailVerificationString: util.ShaHashString(emailParam),
			ID:                      util.GetID(),
			RoleName:                domain.Normal,
		}
		db.SaveUser(user)
		email.SendVerificationEmail(emailParam, user.EmailVerificationString, config.Config.FromEmail)
		w.Header().Add("Location", "/verify")
	} else {
		w.Header().Add("Location", "/signup?emailused=true")
	}
	w.WriteHeader(302)
	w.Write(nil)
	return
}

func validationFailed(w http.ResponseWriter, msg string) {
	w.Header().Add("Location", "/signup?validation="+msg)
	w.WriteHeader(302)
	w.Write(nil)
}

func checkInput(input string) bool {
	re, _ := regexp.Compile(`([\'\";%\n\t\r\0\x08\x1a]+)`)
	return re.MatchString(input)
}

func Verify(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := db.GetUserWithVerifyString(p.ByName("verifystring"))
	if user.ID != "" {
		user.EmailVerified = true
		user.EmailVerifiedDate = time.Now().UTC().Format(time.RFC3339)
		user.ModifyDate = time.Now().UTC().Format(time.RFC3339)
		if user.Email == config.Config.AdminUser {
			user.RoleName = domain.Admin
		}
		ok := db.UpdateUser(user)
		if ok {
			w.Header().Add("Location", "/login?verified=true")
			w.WriteHeader(302)
			w.Write(nil)
			return
		}
	}
	w.Header().Add("Location", "/login?verified=false")
	w.WriteHeader(302)
	w.Write(nil)
}
