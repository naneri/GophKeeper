package apiClient

import (
	"encoding/json"
	"errors"
	"github.com/naneri/GophKeeper/cmd/client/responses"
	"io"
	"log"
	"net/http"
)

type ApiClient struct {
	ServerAddress  string
	PasswordCookie string
}

func (client *ApiClient) TestLogin() error {
	request, err := http.NewRequest(http.MethodGet, client.ServerAddress+"/test-login", nil)
	cookie := http.Cookie{
		Name:   "user",
		Value:  client.PasswordCookie,
		MaxAge: 300,
	}
	request.AddCookie(&cookie)

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		log.Println(res.StatusCode)
		return errors.New("wrong status code returned")
	}

	return nil
}

func (client *ApiClient) GetRecordsList() ([]responses.ResponseRecord, error) {
	var records []responses.ResponseRecord

	request, err := http.NewRequest(http.MethodGet, client.ServerAddress+"/records/list", nil)

	cookie := http.Cookie{
		Name:   "user",
		Value:  client.PasswordCookie,
		MaxAge: 300,
	}
	request.AddCookie(&cookie)

	if err != nil {
		return records, err
	}

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		return records, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return records, err
	}

	err = json.Unmarshal(body, &records)

	return records, err
}
