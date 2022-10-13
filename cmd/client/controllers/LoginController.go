package controllers

import (
	"encoding/json"
	"github.com/naneri/GophKeeper/cmd/client/apiClient"
	"github.com/naneri/GophKeeper/cmd/client/config"
	"log"
	"net/http"
)

type LoginController struct {
	Config *config.Config
}

func (c LoginController) Register(w http.ResponseWriter, r *http.Request) {

}

func (c LoginController) Login(w http.ResponseWriter, r *http.Request) {
	if c.Config.PasswordCookie != "" {
		response := map[string]interface{}{
			"result": "User already logged in",
		}

		err := json.NewEncoder(w).Encode(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		return
	}

	request, err := http.NewRequest(http.MethodPost, c.Config.RemoteServerPort+"/login", r.Body)

	if err != nil {
		log.Println("error generating login request: " + err.Error())
		http.Error(w, "wrong body input", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("error sending request: " + err.Error())
		http.Error(w, "error sending request", http.StatusInternalServerError)
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "user" {
			c.Config.PasswordCookie = cookie.Value
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "error logging in", http.StatusInternalServerError)
	return
}

func (c LoginController) TestLogin(w http.ResponseWriter, r *http.Request) {
	serverClient := apiClient.ApiClient{
		ServerAddress:  c.Config.RemoteServerPort,
		PasswordCookie: c.Config.PasswordCookie,
	}

	err := serverClient.TestLogin()

	if err != nil {
		http.Error(w, "error logging in: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
