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
	Category string `json:"category"`
	Upvotes int32 `json:"upvotes"`
	UserID int `json:"userID"`
}

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}

type Comment struct{
	PostID int `json:"postID"`
	UserID int `json:"userID"`
	Comment string `json:"comment"`
}