package middlewares

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func AuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie := make([]byte, aes.BlockSize)

		secretKey := []byte("sbHYDYWgdakkHHDS")

		nonce := []byte("YANDEX")

		aesblock, err := aes.NewCipher(secretKey)
		if err != nil {
			log.Infof("Cannot inicialize symmetric encryption interface %v", err)
		}

		if requestUserID, err := r.Cookie("user_id"); err == nil {
			requestUserIDByte, err := hex.DecodeString(requestUserID.Value)
			if err != nil {
				log.Infof("Auth Cookie decoding: %v\n", err)
			}
			aesblock.Decrypt(authCookie, requestUserIDByte)
			if string(authCookie[len(authCookie)-len(nonce):]) == string(nonce) {
				next.ServeHTTP(w, r)
				return //	если ДА, то проверка подлинности пройдена
			}
		}

		userID, _ := generateRandom(10)
		aesblock.Encrypt(authCookie, append(userID, nonce...)) // зашифровываем (UserID + nonce) в переменную authCookie

		cookie := &http.Cookie{
			Name: "userid", Value: hex.EncodeToString(authCookie), Expires: time.Now().AddDate(1, 0, 0),
		}

		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
		next.ServeHTTP(w, r)
	})
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
