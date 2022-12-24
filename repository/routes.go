package repository

import (
	"github.com/gofiber/fiber/v2"
)

func (repo *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Get("/posts", repo.GetPosts)
	api.Post("/posts", repo.CreatePost)
	api.Patch("/posts/:id", repo.UpdatePost)
	api.Delete("/posts/:id", repo.DeletePost)
	api.Get("/posts/:id", repo.GetPostByID)
	api.Post("/register", repo.Register)
	api.Post("/login", repo.Login)
	api.Get("/user", repo.User)
	api.Post("/logout", repo.Logout)

}