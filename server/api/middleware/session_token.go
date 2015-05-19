package middleware

import (
	"errors"
	"github.com/gorilla/context"
	"github.com/mhannig/papertrail/server/models"
	"log"
	"net/http"
	"strings"
)

const (
	SessionKey  int    = 0
	SessionAuth string = "Session"
)

func AuthSessionToken(app http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Parse HTTP Authorization Header
		// In case the Authorization SessionToken <token> Header is
		// set try to lookup the session in the database
		authorizations := req.Header["Authorization"]
		if authorizations == nil {
			goto next
		}

		for _, authorization := range authorizations {
			tokens := strings.SplitN(authorization, " ", 2)
			if len(tokens) > 1 {
				if tokens[0] != SessionAuth {
					continue // Skip this.
				}

				// We have a session token.
				sessionToken := tokens[1]

				// Try to get the corresponding session from the DB
				session, err := models.FindSessionByToken(sessionToken)
				if err != nil {
					log.Println("[DB] Error while getting session by token:", err)
					continue
				}

				// Check if session is expired.
				if session.TTL() == 0 {
					continue
				}

				// Everythings fine? Make sure the session is up to date.
				session.Touch()

				// Looks like we have some valid session. Hooray.
				SetCurrentSession(req, session)
			}
		}

	next:
		app.ServeHTTP(res, req)
	})
}

func SetCurrentSession(req *http.Request, session *models.Session) {
	context.Set(req, SessionKey, session)
}

func CurrentSession(req *http.Request) (*models.Session, error) {
	session := context.Get(req, SessionKey)
	if session != nil {
		return session.(*models.Session), nil
	}
	return nil, errors.New("Session not authorized.")
}

func AssertAuthenticated(req *http.Request) error {
	_, err := CurrentSession(req)
	return err
}
