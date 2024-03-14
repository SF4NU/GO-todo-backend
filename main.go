package main

import (
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
)

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Deadline    string `json:"deadline"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	dbKey := os.Getenv("DB_KEY")

	dsn := fmt.Sprintf("postgres://ryfljouh:%s@dumbo.db.elephantsql.com/ryfljouh", dbKey)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("db not connected")
	}

	db.AutoMigrate(&Task{})

	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/tasks", func(c echo.Context) error {
		var tasks []Task
		db.Find(&tasks)
		return c.JSON(http.StatusOK, tasks)
	})

	e.POST("/tasks", func(c echo.Context) error {
		var task Task
		if err := c.Bind(&task); err != nil {
			return err
		}

		db.Create(&task)
		return c.JSON(http.StatusCreated, task)
	})

	e.PUT("/tasks/:id", func(c echo.Context) error {
		id := c.Param("id")

		var task Task
		if err := db.First(&task, id).Error; err != nil {
			return c.String(http.StatusNotFound, "task not found")
		}

		var updatedTask Task
		if err := c.Bind(&updatedTask); err != nil {
			return err
		}
		if updatedTask.Description != "" {
			task.Description = updatedTask.Description
		}

		if updatedTask.Completed != task.Completed {
			task.Completed = updatedTask.Completed
		}
		if updatedTask.Deadline != task.Deadline {
			task.Deadline = updatedTask.Deadline
		}

		if err := db.Save(&task).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, task)
	})

	e.DELETE("/tasks/:id", func(c echo.Context) error {
		id := c.Param("id")

		var task Task
		if err := db.First(&task, id).Error; err != nil {
			return c.String(http.StatusNotFound, "task not found")
		}

		db.Delete(&Task{}, id)

		return c.String(http.StatusAccepted, "task deleted")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
