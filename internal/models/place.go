package models

type Places struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Phone   string  `json:"phone"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lot"`
}
