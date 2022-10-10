package middleware

import (
	"context"
	"github.com/naneri/GophKeeper/cmd/server/config"
	"net/http"
)

const AppConfigKey = "AppConfig"

type AppConfig string

type ConfigMiddlewareStruct struct {
	Config *config.Config
}

func (c *ConfigMiddlewareStruct) SetConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(context.WithValue(r.Context(), AppConfig(AppConfigKey), c.Config))

		*r = *req

		next.ServeHTTP(w, r)
	})
}
