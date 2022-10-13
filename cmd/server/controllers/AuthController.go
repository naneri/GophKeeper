package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/naneri/GophKeeper/cmd/server/config"
	"github.com/naneri/GophKeeper/cmd/server/dto"
	"github.com/naneri/GophKeeper/cmd/server/middleware"
	"github.com/naneri/GophKeeper/internal/app/security"
	"github.com/naneri/GophKeeper/internal/app/user"
	"io"
	"log"
	"net/http"
	"time"
)

type AuthController struct {
	UserRepo user.RepositoryInterface
	Config   *config.Config
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.RegisterParams

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error parsing request body: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &requestBody)

	if err != nil {
		log.Println("Error unmarshalling request body: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPass := security.Hash(c.Config.AppKey, requestBody.Password)

	userId, err := c.UserRepo.Store(requestBody.Login, hashedPass)

	if err != nil {
		log.Println("Error storing user: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userCookie := generateUserCookie(userId, c.Config.AppKey)
	http.SetCookie(w, &userCookie)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserID(middleware.UserIDContextKey))

	if userId != nil {
		response := map[string]interface{}{
			"result": "User already logged in",
			"userId": userId,
		}

		err := json.NewEncoder(w).Encode(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		return
	}

	var loginParams dto.LoginParams
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error parsing request body: " + err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &loginParams)

	if err != nil {
		log.Println("Error unmarshalling request body: " + err.Error())
		http.Error(w, "Incorrect data provided", http.StatusInternalServerError)
		return
	}

	userRecord, err := c.UserRepo.GetByLogin(loginParams.Login)

	if err != nil {
		log.Println("User with this login not found: " + err.Error())
		http.Error(w, "User with this login not found", http.StatusUnauthorized)
		return
	}

	check, err := security.CheckHash(c.Config.AppKey, loginParams.Password, userRecord.Password)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if !check {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	userCookie := generateUserCookie(userRecord.Id, c.Config.AppKey)
	http.SetCookie(w, &userCookie)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) TestLoggedIn(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.UserID(middleware.UserIDContextKey))

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func generateUserCookie(userId uint32, secretKey string) http.Cookie {
	uint32userIDBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(uint32userIDBuf[0:], userId)

	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write(uint32userIDBuf)
	sign := hash.Sum(uint32userIDBuf)
	userCookie := hex.EncodeToString(sign)

	expire := time.Now().Add(10 * time.Minute)
	httpCookie := http.Cookie{Name: "user", Value: userCookie, Path: "/", Expires: expire, MaxAge: 90000}

	return httpCookie
}
