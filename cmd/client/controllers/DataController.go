package controllers

import (
	"encoding/json"
	"github.com/naneri/GophKeeper/cmd/client/config"
	"github.com/naneri/GophKeeper/cmd/client/storage"
	"net/http"
)

type DataController struct {
	Config        *config.Config
	RecordStorage *storage.RecordStorage
}

func (c DataController) GetRecordList(w http.ResponseWriter, r *http.Request) {
	records, _ := c.RecordStorage.ListUserRecords()

	err := json.NewEncoder(w).Encode(records)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	return
}
