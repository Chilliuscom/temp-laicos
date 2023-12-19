package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/michaeljs1990/sqlitestore"

	"github.com/gofrs/uuid"
)

const cookieExpireTime = 15 //minutes

var store *sqlitestore.SqliteStore

func init() {
	var err error
	store, err = sqlitestore.NewSqliteStore("./social-network.db", "sessions", "/", 3600, []byte("<Super-Secret-Key"))
	if err != nil {
		panic(err)
	}
}

func validateSession(w http.ResponseWriter, r *http.Request) bool {
	cookie, _ := r.Cookie("session-id")
	if cookie == nil {
		return false
	}
	session, _ := store.Get(r, "session-id")
	expires, ok := session.Values["expires"].(time.Time)
	if !ok {
		return false
		// Handle the case where "expires" is not a valid time.Time
		// You can set a default expiration or return an error
	}
	//expires.(time).After(time.Now())
	if !expires.After(time.Now()) {
		http.SetCookie(w, &http.Cookie{ //empties the cookie
			Name:     "session-id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,                  // Ensure it's a Secure (HTTPS) cookie
			SameSite: http.SameSiteNoneMode, // Use the integer value for "None"
		})
		return false

	}

	return true

}

func newSession(w http.ResponseWriter, username string) {
	sessionToken, _ := uuid.NewV4()
	//DELETE FROM Sessions WHERE userID = ?  //delete previous sessions, allowing 1 session per user
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    sessionToken.String(),
		Expires:  time.Now().Add(cookieExpireTime * time.Minute),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	//ADD to db session, userid
	log.Println(username, ": ", sessionToken.String())
}

func endSession(w http.ResponseWriter, r *http.Request) {
	//DELETE FROM Sessions WHERE session-id = r.Cookie("session-id")
	http.SetCookie(w, &http.Cookie{ //empties the cookie
		Name:     "session-id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func getUserBySession(r *http.Request) int {
	session, err := store.Get(r, "session-id")
	if err != nil {
		return 0
	}
	//SELECT username FROM sessions where key = ?, r.Cookie(SessionID).Value
	if session.Values["id"] != nil {
		return (session.Values["id"]).(int)
	}
	return 0
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateSession(w, r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
