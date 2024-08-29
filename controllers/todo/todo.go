package todo

import (
	"go-fiber-server/models"
	"go-fiber-server/storage"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	storage *storage.TodoStorage
}

func NewController(storage *storage.TodoStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

func (tc *Controller) RegisterRoutes(app fiber.Router) {
	app.Get("/todos", tc.GetTodos)
	app.Get("/todos/:id", tc.GetTodo)
	app.Post("/todos", tc.CreateTodo)
	app.Patch("/todos/:id", tc.PatchTodo)
	app.Put("/todos/:id", tc.UpdateTodo)
	app.Delete("/todos/:id", tc.DeleteTodo)
}

func (tc *Controller) GetTodos(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)

	todos, err := tc.storage.GetTodos(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(todos)
}

func (tc *Controller) GetTodo(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)
	todoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	todo, err := tc.storage.GetTodoByID(userID, todoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(todo)
}

func (tc *Controller) CreateTodo(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)

	var todoPayload models.TodoPayload
	if err := c.BodyParser(&todoPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if todoPayload.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}
	dbTodo, err := tc.storage.CreateTodo(&todoPayload, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(dbTodo)
}

func (tc *Controller) PatchTodo(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)
	todoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	todo, err := tc.storage.GetTodoByID(userID, todoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if todo.User_ID != userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request: wrong todo id"})
	}

	var todoPayload models.TodoPayload
	if err := c.BodyParser(&todoPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Patch only the fields that are provided else use the existing values
	if todoPayload.Title == "" {
		todoPayload.Title = todo.Title
	}
	if !todoPayload.Completed {
		todoPayload.Completed = todo.Completed
	}

	if err := tc.storage.UpdateTodo(&todoPayload, todoID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}

func (tc *Controller) UpdateTodo(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)
	todoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var todoPayload models.TodoPayload
	if err := c.BodyParser(&todoPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	todo, err := tc.storage.GetTodoByID(userID, todoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if todo.User_ID != userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request: wrong todo id"})
	}
	if err := tc.storage.UpdateTodo(&todoPayload, todoID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}

func (tc *Controller) DeleteTodo(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)
	todoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	todo, err := tc.storage.GetTodoByID(userID, todoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if todo.User_ID != userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request: wrong todo id"})
	}
	if err := tc.storage.DeleteTodo(todoID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}
