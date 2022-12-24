package models

// import (
// 	pq "github.com/lib/pq"
// )

// type TagIDs struct {
// 	tagIDs []int64
// }


type Post struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Upvotes int32 `json:"upvotes"`
	UserID int `json:"userID"`
}

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}