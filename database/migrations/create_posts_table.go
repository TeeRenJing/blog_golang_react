package migrations

import (
	// "time"

	// pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title string `json:"title"`
	Content string `json:"content"`
	Category string `json:"category"`
	Upvotes int32 `json:"upvotes"`
	UserID int
  	User   User
}

type Comment struct {
	gorm.Model
	Comment string `json:"comment"`
	PostID int `json:"postID"`
	Post  Post
	UserID int `json:"userID"`
	User  User
}

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}

func MigratePosts(db *gorm.DB) error {
	err := db.AutoMigrate(&Post{}, &Comment{}, &User{})
	return err
}