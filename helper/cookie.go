package helper

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"net/http"
	"os"
)

var hashKey = []byte(os.Getenv("COOKIE_SECRET"))
var s = securecookie.New(hashKey, nil)

func SetCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string, cookieValue string) {
	encoded, err := s.Encode(cookieName, cookieValue)
	if err == nil {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		var value string
		if err = s.Decode(cookieName, cookie.Value, &value); err == nil {
			fmt.Fprintln(w, value)
		}
	}
}
