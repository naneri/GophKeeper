package dto

type ResponseRecord struct {
	Id        uint32 `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Data      string `json:"data,omitempty"`
	Path      string `json:"path,omitempty"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
	UpdatedAt string `json:"updated_at"`
}
