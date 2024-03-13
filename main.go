package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func getUser(c echo.Context) error {
// 	id := c.Param("id")
// 	return c.String(http.StatusOK, id)
// }

// type User struct {
// 	Name string `json:"name"`
// 	Email string `json:"email"`
// }

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func main() {
	dsn := "host=localhost user=postgres password=123 dbname=todo_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("db not connected")
	}

	db.AutoMigrate(&Task{})

	e := echo.New()

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

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	// e.POST("/users", func(c echo.Context) error {
	// 	u := new(User)
	// 	if err := c.Bind(u); err != nil {
	// 		return err
	// 	}

	// 	fmt.Println(u.Email)
	// 	return c.JSON(http.StatusCreated, u)
	// })

	// e.GET("/users/:id", getUser)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
