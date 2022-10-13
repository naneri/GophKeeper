package dto

import "encoding/json"

type RecordParams struct {
	Name string          `json:"name,omitempty"`
	Type string          `json:"type,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}
