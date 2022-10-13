package record

type RepositoryInterface interface {
	GetRecord(id uint32) (Record, error)
	StoreRecord(record Record) (uint32, error)
	ListUserRecords(userId uint32) ([]Record, error)
}
