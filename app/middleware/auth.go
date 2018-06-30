package middleware

import (
	"net/http"

	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/julienschmidt/httprouter"
)

func Auth(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		loginCookie, err := r.Cookie("login")
		if err == nil {
			if loginCookie != nil {
				session := db.GetSession(loginCookie.Value)
				if session != "" {
					fn(w, r, ps)
					return
				}
			}
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(403)
		w.Write([]byte(`{"error":"no rights"}`))
	}
}

func AdminAuth(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		loginCookie, err := r.Cookie("login")
		if err == nil {
			if loginCookie != nil {
				session := db.GetSession(loginCookie.Value)
				if session != "" {
					if db.GetUser(session).RoleName == domain.Admin {
						fn(w, r, ps)
						return
					}
				}
			}
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(403)
		w.Write([]byte(`{"error":"no rights"}`))
	}
}
