package storage

import (
	"github.com/naneri/GophKeeper/cmd/client/responses"
	"sync"
)

type RecordStorage struct {
	AccessLocker *sync.Mutex
	UserRecords  []responses.ResponseRecord
}

func (storage *RecordStorage) ListUserRecords() ([]responses.ResponseRecord, error) {
	return storage.UserRecords, nil
}

func (storage *RecordStorage) UpdateRecords(records []responses.ResponseRecord) {
	storage.UserRecords = records
}
