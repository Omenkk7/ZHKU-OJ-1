package dto

// ReqLabel 由于增删改查时，前端都传这四个参数，所以直接命名为ReqLabel
type ReqLabel struct {
	ID      string `json:"_id"`
	Creator string `json:"creator"`
	Status  string `json:"status"`
	Name    string `json:"name"`
}
