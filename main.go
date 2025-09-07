// echo, for learning purpose.

package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ---------- Types ----------

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

type CreateUserReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{}
var nextID = 1

// ---------- Handlers ----------

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func hello(c echo.Context) error {
	name := c.QueryParam("name")
	if strings.TrimSpace(name) == "" {
		name = "world"
	}
	return c.JSON(http.StatusOK, echo.Map{"hello": name})
}

func createUser(c echo.Context) error {
	var req CreateUserReq
	if err := c.Bind(&req); err != nil {
		return &AppError{Code: http.StatusBadRequest, Message: "invalid json"}
	}
	if strings.TrimSpace(req.Name) == "" || req.Age <= 0 {
		return &AppError{Code: http.StatusBadRequest, Message: "name and age are required"}
	}
	u := User{ID: nextID, Name: req.Name, Age: req.Age}
	nextID++
	users = append(users, u)
	return c.JSON(http.StatusCreated, u)
}

func listUsers(c echo.Context) error {
	minAge := 0
	if s := c.QueryParam("min_age"); s != "" {
		var err error
		minAge, err = strconv.Atoi(s)
		if err != nil || minAge < 0 {
			return &AppError{Code: http.StatusBadRequest, Message: "invalid min_age"}
		}
	}
	if minAge == 0 {
		return c.JSON(http.StatusOK, users)
	}
	out := make([]User, 0, len(users))
	for _, u := range users {
		if u.Age >= minAge {
			out = append(out, u)
		}
	}
	return c.JSON(http.StatusOK, out)
}

// ---------- main ----------

func main() {
	e := echo.New()

	// Centralized error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var app *AppError
		if errors.As(err, &app) {
			_ = c.JSON(app.Code, echo.Map{"error": app.Message})
			return
		}
		c.Logger().Error(err)
		_ = c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	// Routes
	e.GET("/healthz", health)
	e.GET("/hello", hello)
	e.GET("/users", listUsers)
	e.POST("/users", createUser)

	// Server config + graceful shutdown
	go func() {
		e.Logger.Print("Echo server listening on :8080")
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal(err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

