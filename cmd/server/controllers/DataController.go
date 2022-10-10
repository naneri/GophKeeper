package controllers

import (
	"encoding/json"
	"errors"
	"github.com/naneri/GophKeeper/cmd/server/config"
	"github.com/naneri/GophKeeper/cmd/server/dto"
	"github.com/naneri/GophKeeper/cmd/server/middleware"
	"github.com/naneri/GophKeeper/internal/app/record"
	"io"
	"log"
	"net/http"
)

type DataController struct {
	RecordRepo record.DatabaseRepository
	Config     *config.Config
}

func (c DataController) Store(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserID(middleware.UserIDContextKey))

	userIdUint, ok := userId.(uint32)

	if !ok {
		log.Println("User id is not uint32 check the correct cookie setting")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	var RecordParams dto.RecordParams

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error parsing request body: " + err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &RecordParams)

	if err != nil {
		log.Println("Error unmarshalling request body: " + err.Error())
		http.Error(w, "Incorrect data provided", http.StatusBadRequest)
		return
	}

	dbRecord, err := record.ParseRecord(RecordParams.Name, RecordParams.Type, RecordParams.Data)

	if err != nil {
		http.Error(w, "Incorrect data provided", http.StatusBadRequest)
		return
	}

	recordId, err := c.RecordRepo.StoreRecord(userIdUint, dbRecord)

	if err != nil {
		if errors.Is(err, &record.NameAlreadyUsedError{}) {
			http.Error(w, "Name already used", http.StatusBadRequest)
			return
		}

		http.Error(w, "Incorrect data provided", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"record_id": recordId,
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	return
}

func (c DataController) List(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserID(middleware.UserIDContextKey))

	userIdUint, ok := userId.(uint32)

	if !ok {
		log.Println("User id is not uint32 check the correct cookie setting")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	userRecords, err := c.RecordRepo.ListUserRecords(userIdUint)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var response []dto.ResponseRecord

	for _, dbRecord := range userRecords {
		var deletedAt string
		if dbRecord.DeletedAt.Time.IsZero() {
			deletedAt = ""
		} else {
			deletedAt = dbRecord.DeletedAt.Time.String()
		}
		newResponseRecord := dto.ResponseRecord{
			Id:        dbRecord.ID,
			Name:      dbRecord.Name,
			Type:      dbRecord.Type,
			Data:      dbRecord.Data,
			Path:      "",
			CreatedAt: dbRecord.CreatedAt.String(),
			DeletedAt: deletedAt,
			UpdatedAt: dbRecord.UpdatedAt.String(),
		}

		response = append(response, newResponseRecord)
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	return
}
