package record

const typeCredentials = "credentials"
const typeCreditCardData = "credit-card-data"

type TypeCredentials struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type TypeCreditCardData struct {
	FullName string `json:"full_name,omitempty"`
	Number   int64  `json:"number,omitempty"`
	Month    uint8  `json:"month,omitempty"`
	Year     uint16 `json:"year,omitempty"`
}
