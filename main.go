package main

import (
	"github.com/TeeRenJing/blog_golang_react/bootstrap"
	"github.com/TeeRenJing/blog_golang_react/repository"
	"github.com/gofiber/fiber/v2"
)

type Repository repository.Repository

func main() {
	app := fiber.New()
	bootstrap.InitializeApp(app)


}