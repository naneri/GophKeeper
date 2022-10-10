package recordKeeper

import (
	"github.com/naneri/GophKeeper/internal/app/record"
	"os"
)

type KeeperInterface interface {
	Fetch(record record.Record) (os.File, error)
	Store(data []byte) (string, error)
	Delete(record record.Record) error
}
