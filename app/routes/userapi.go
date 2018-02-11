package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/util"
	"github.com/julienschmidt/httprouter"
)

func UserJSON(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var marshalled []byte
	var err error
	if user := getUser(r); user.ID != "" {
		marshalled, err = json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
	}
	ok(w, marshalled)
}

func AllUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	users := db.GetUsers()
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(usersJSON)
}

func getUser(r *http.Request) domain.User {
	loginCookie, err := r.Cookie("login")
	var user domain.User
	if err == nil && loginCookie != nil {
		session := db.GetSession(loginCookie.Value)
		if session != "" {
			user = db.GetUser(util.Decrypt(session))
		}
	}
	return user
}
