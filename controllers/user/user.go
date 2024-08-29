package user

import (
	"go-fiber-server/auth"
	"go-fiber-server/models"
	"go-fiber-server/storage"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	storage *storage.UserStorage
}

func NewController(storage *storage.UserStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

func (u *Controller) RegisterRoutes(app fiber.Router) {
	app.Post("/login", u.Login)
	app.Post("/register", u.Register)
}

var validate = validator.New()

func (u *Controller) Register(c *fiber.Ctx) error {
	var userPayload models.RegisterPayload

	if err := c.BodyParser(&userPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(userPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := u.storage.GetUserByEmail(userPayload.Email)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already exists"})
	}

	hashedPassword, err := auth.HashPassword(userPayload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	userPayload.Password = hashedPassword

	user, err := u.storage.Create(&userPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		SameSite: "strict",
		// Secure:   true,
	})
	return c.JSON(userPayload)
}

func (u *Controller) Login(c *fiber.Ctx) error {
	var loginPayload models.LoginPayload

	if err := c.BodyParser(&loginPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(loginPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := u.storage.GetUserByEmail(loginPayload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email; " + err.Error()})
	}

	if !auth.ComparePassword(user.Password, loginPayload.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password"})
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		SameSite: "strict",
		// Secure:   true,
	})
	return c.JSON(loginPayload)
}
