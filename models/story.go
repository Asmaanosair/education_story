package models

type Story struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Text   string `json:"text"`
	Quiz   string `json:"quiz"`
	Answer string `json:"answer"`
	Score  int    `json:"score"`
}
