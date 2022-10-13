package record

import (
	"gorm.io/gorm"
)

type DatabaseRepository struct {
	Db *gorm.DB
}

func (d DatabaseRepository) GetRecord(id uint32) (Record, error) {

	panic("implement me")
}

func (d DatabaseRepository) StoreRecord(userId uint32, record Record) (uint32, error) {
	var existingRecord Record

	result := d.Db.Where("user_id = ?", userId).Where("name = ?", record.Name).Find(&existingRecord)

	if result.RowsAffected != 0 {
		return 0, &NameAlreadyUsedError{}
	}

	record.UserId = userId
	result = d.Db.Create(&record)

	if result.Error != nil {
		return 0, result.Error
	}

	return record.ID, nil
}

func (d DatabaseRepository) ListUserRecords(userId uint32) ([]Record, error) {
	var records []Record

	result := d.Db.Where("user_id = ?", userId).Find(&records)

	return records, result.Error
}
