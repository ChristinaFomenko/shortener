package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

func Decompressing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader

		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			reader = gz
			defer func(gz *gzip.Reader) {
				_ = gz.Close()
			}(gz)
		} else {
			reader = r.Body
		}

		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}
