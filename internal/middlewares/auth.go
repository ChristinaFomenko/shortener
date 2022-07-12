package middlewares

import (
	"context"
	"fmt"
	"net/http"
)

const (
	authCookieName         = "user-id"
	AuthTokenKey   AuthKey = "user-id-ctx"
)

type AuthKey string

type authService interface {
	SignUp() (string, string, error)
	SignIn(token string) (string, error)
}

type authenticator struct {
	authService authService
}

func NewAuthenticator(authService authService) *authenticator {
	return &authenticator{
		authService: authService,
	}
}

// Auth Авторизация пользователя
func (a *authenticator) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		token, err := a.getAuthToken(r)
		if err != nil {
			userID, token, err = a.authService.SignUp()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			userID, err = a.authService.SignIn(token)
			if err != nil {
				userID, token, err = a.authService.SignUp()

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		a.setUserToken(w, token)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AuthTokenKey, userID)))
	})
}

func (a *authenticator) signUp() {
}

func (a *authenticator) setUserToken(w http.ResponseWriter, userID string) {
	cookie := http.Cookie{
		Name:  authCookieName,
		Value: userID,
		Path:  "/",
	}

	http.SetCookie(w, &cookie)
}

func (a *authenticator) getAuthToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		return "", fmt.Errorf("get user cookie error: %w", err)
	}

	return cookie.Value, nil
}

func (a *authenticator) UserID(ctx context.Context) string {
	return ctx.Value(AuthTokenKey).(string)
}
