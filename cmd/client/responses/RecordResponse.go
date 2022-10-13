package responses

type ResponseRecord struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Data      string `json:"data"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
	UpdatedAt string `json:"updated_at"`
}
