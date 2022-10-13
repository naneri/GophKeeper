package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/GophKeeper/cmd/server/config"
	"github.com/naneri/GophKeeper/cmd/server/controllers"
	"github.com/naneri/GophKeeper/cmd/server/middleware"
	"github.com/naneri/GophKeeper/internal/app/record"
	"github.com/naneri/GophKeeper/internal/app/user"
	"gorm.io/gorm"
	"net/http"
)

type Router struct {
}

func (router *Router) GetHandler(db *gorm.DB, config *config.Config) *chi.Mux {
	userRepo := user.DatabaseRepository{Db: db}
	recordRepo := record.DatabaseRepository{Db: db}
	r := chi.NewRouter()
	authController := controllers.AuthController{
		UserRepo: userRepo,
		Config:   config,
	}

	dataController := controllers.DataController{
		RecordRepo: recordRepo,
		Config:     config,
	}

	configMiddleware := middleware.ConfigMiddlewareStruct{Config: config}

	r.Use(configMiddleware.SetConfig)
	r.Use(middleware.SetUserIdMiddleware)

	r.Post("/register", authController.Register)
	r.Post("/login", authController.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.CheckAuthMiddleware)
		r.Get("/test-login", authController.TestLoggedIn)
		r.Post("/records/store", dataController.Store)
		r.Get("/records/list", dataController.List)
	})

	r.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		_, _ = writer.Write([]byte(fmt.Sprintf("test")))
	})

	return r
}
