package jobs

import (
	"github.com/naneri/GophKeeper/cmd/client/apiClient"
	"github.com/naneri/GophKeeper/cmd/client/config"
	"github.com/naneri/GophKeeper/cmd/client/storage"
	"log"
)

func UpdateRecords(cfg *config.Config, storage *storage.RecordStorage, client *apiClient.ApiClient) error {
	if cfg.PasswordCookie != "" {
		if client == nil {
			client = &apiClient.ApiClient{
				ServerAddress:  cfg.RemoteServerPort,
				PasswordCookie: cfg.PasswordCookie,
			}
		}

		records, recordErr := client.GetRecordsList()
		if recordErr != nil {
			return recordErr
		}

		log.Println("records updated")
		storage.UpdateRecords(records)
	} else {
		log.Println("nothing to update")
	}

	return nil
}
