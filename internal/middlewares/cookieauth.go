package middlewares

import (
	"context"
	"github.com/ChristinaFomenko/shortener/internal/app/generator"
	"net/http"
)

type Ctxkey struct{}

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			uid         string
			cookieValue string
		)
		idCookie, err := r.Cookie("user_id")
		if err != nil {
			uid, cookieValue = generator.GenerateNewUserCookie()
			cookie := http.Cookie{Name: "user_id", Value: cookieValue}
			http.SetCookie(w, &cookie)
		} else {
			cookieValue = idCookie.Value
			uid, err = generator.GetUserIDFromCookie(cookieValue)
			if err != nil {
				uid, cookieValue = generator.GenerateNewUserCookie()
				cookie := http.Cookie{Name: "uid", Value: cookieValue}
				http.SetCookie(w, &cookie)
			}

		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), Ctxkey{}, uid)))
	})

}
