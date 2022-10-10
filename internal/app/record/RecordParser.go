package record

import (
	"encoding/json"
	"errors"
)

func ParseRecord(name, recordType string, data json.RawMessage) (Record, error) {
	var recordData string
	var validationError error

	if recordType == typeCredentials {
		recordData, validationError = validateCredentialsType(data)
	} else if recordType == typeCreditCardData {
		recordData, validationError = validateCreditCardData(data)
	} else {
		return Record{}, errors.New("wrong data type")
	}

	if validationError != nil {
		return Record{}, errors.New("error validating data: " + validationError.Error())
	}

	return Record{
		Name: name,
		Type: recordType,
		Data: recordData,
	}, nil
}

func validateCredentialsType(data json.RawMessage) (string, error) {
	var record TypeCredentials

	err := json.Unmarshal(data, &record)

	return string(data), err
}

func validateCreditCardData(data json.RawMessage) (string, error) {
	var record TypeCreditCardData

	err := json.Unmarshal(data, &record)

	return string(data), err
}
