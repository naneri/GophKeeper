package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/naneri/GophKeeper/cmd/server/config"
	"log"
	"net/http"
)

func SetUserIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			data   []byte
			err    error
			idSign []byte
		)

		// parse cookie
		cookie, err := r.Cookie("user")

		// can I make this code prettier?
		if err != nil {
			log.Println(err.Error())
			next.ServeHTTP(w, r)
			return
		}
		log.Println(cookie.Value)
		data, err = hex.DecodeString(cookie.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		cfg := r.Context().Value(AppConfig(AppConfigKey)).(*config.Config)
		userID := binary.BigEndian.Uint32(data[:4])
		h := hmac.New(sha256.New, []byte(cfg.AppKey))
		h.Write(data[:4])
		idSign = h.Sum(nil)

		// if parse correctly, add the cookie to context
		if !hmac.Equal(idSign, data[4:]) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, UserID(UserIDContextKey), userID))
		*r = *req
		// else grant user the signed cookie with Unique identifier
		next.ServeHTTP(w, r)
	})
}
