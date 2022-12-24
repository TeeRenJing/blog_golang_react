package repository

import (
	"net/http"

	"strconv"
	"time"

	"github.com/TeeRenJing/blog_golang_react/database/migrations"
	"github.com/TeeRenJing/blog_golang_react/database/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/morkid/paginate"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type ErrorResponse struct {
	FailedField string
	Tag string
	Value string
}

var validate = validator.New()

func ValidateStruct(post models.Post) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(post)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors){
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (r *Repository) GetPosts(context *fiber.Ctx) error {
	db := r.DB
	model := db.Model(&migrations.Post{})
	pg := paginate.New(&paginate.Config{
		DefaultSize: 20,
		CustomParamEnabled: true,
	})

	page := pg.With(model).Request(context.Request()).Response(&[]migrations.Post{})

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"data": page,
	})
	return nil
}

func (r *Repository) CreatePost(context *fiber.Ctx) error {
	post := models.Post{}
	err := context.BodyParser(&post)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message": "Request failed",
			},
		)
		return err
	}
	errors := ValidateStruct(post)
	if errors != nil {
		return context.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if err := r.DB.Create(&post).Error; err != nil {

		return context.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status": "error", 
			"message": "Couldn't create post", 
			"data": err,
		})
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Post has been added", 
		"data": post,
	})

	return nil
}

func (r *Repository) UpdatePost(context *fiber.Ctx) error {
	post := models.Post{}
	err := context.BodyParser(&post)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message": "Request failed",
			},
		)
		return err
	}

	errors := ValidateStruct(post)
	if errors != nil {
		return context.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := r.DB
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	if db.Model(&post).Where("id = ?", id).Updates(&post).RowsAffected == 0 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get post"})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"status": "success", "message": "Post successfully updated"})
	return nil
}

func (r *Repository) DeletePost(context *fiber.Ctx) error {
	postModel := &migrations.Post{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	err := r.DB.Delete(postModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not delete"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"status": "success", "message": "Post successfully deleted"})
	return nil
}

func (r *Repository) GetPostByID(context *fiber.Ctx) error {
	postModel := &migrations.Post{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(postModel).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get the post"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"status": "success", "message": "User profile fetched successfully", "data": postModel})
	return nil
}

const SecretKey = "secret"

func (r *Repository) Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	r.DB.Create(&user)

	return c.JSON(user)
}

func (r *Repository) Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	r.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (r *Repository) User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	r.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func (r *Repository) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
