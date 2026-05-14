package model

type Tag struct {
	Id    string `gorm:"primaryKey"`
	Title string
}
