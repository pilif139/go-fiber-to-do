package main

import (
	"go-fiber-server/controllers/todo"
	"go-fiber-server/controllers/user"
	"go-fiber-server/db"
	"go-fiber-server/storage"
	"log"

	"go-fiber-server/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := fiber.New()

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Use("/todos", AuthenticateMiddleware)

	userStorage := storage.NewUserStorage(db)
	userController := user.NewController(userStorage)
	userController.RegisterRoutes(app)

	todoStorage := storage.NewTodoStorage(db)
	todoController := todo.NewController(todoStorage)
	todoController.RegisterRoutes(app)

	app.Listen(":3000")
}

func AuthenticateMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	if tokenString == "" {
		return c.Redirect("/login", fiber.StatusUnauthorized)
	}
	token, err := auth.VerifyToken(tokenString)
	if err != nil {
		return c.Redirect("/login", fiber.StatusUnauthorized)
	}
	userId, err := auth.GetUserId(token)
	if err != nil {
		return c.Redirect("/login", fiber.StatusUnauthorized)
	}
	c.Locals("userId", userId)
	return c.Next()
}
