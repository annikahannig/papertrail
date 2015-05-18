package middleware

import (
	"github.com/gorilla/context"
	"github.com/mhannig/papertrail/server/models"
	"log"
	"net/http"
)

const (
	SessionKey int = 0
)

func AuthSessionToken(app http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Parse HTTP Authorization Header
		// In case the Authorization SessionToken <token> Header is
		// set try to lookup the session in the database
		log.Println(req.Header)

		app.ServeHTTP(res, req)
	})
}

func SetCurrentSession(req *http.Request, session *models.Session) {
	context.Set(req, SessionKey, session)
}

func CurrentSession(req *http.Request) *models.Session {
	session := context.Get(req, SessionKey)
	if session != nil {
		return session.(*models.Session)
	}
	return nil
}
