package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/naneri/GophKeeper/cmd/client/config"
	"github.com/naneri/GophKeeper/cmd/client/controllers"
	"github.com/naneri/GophKeeper/cmd/client/storage"
)

type Router struct {
}

func (router *Router) GetHandler(config *config.Config, storage *storage.RecordStorage) *chi.Mux {
	r := chi.NewRouter()

	loginController := controllers.LoginController{Config: config}
	dataController := controllers.DataController{Config: config, RecordStorage: storage}

	r.Post("/login", loginController.Login)

	// routes that need cookie set
	r.Group(func(r chi.Router) {
		r.Get("/test-login", loginController.TestLogin)
		r.Get("/records/list", dataController.GetRecordList)
	})
	return r
}
