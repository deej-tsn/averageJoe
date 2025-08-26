package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/deej-tsn/averageJoe/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	data_path, _ := filepath.Abs("../data/questions.json")

	data_bytes, err := os.ReadFile(data_path)
	if err != nil {
		panic(err)
	}

	data := model.LoadData(data_bytes)

	game := model.NewGame(data.GetRandomRound())

	server := echo.New()

	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} body=${custom}\n",
	}))

	// DEV ONLY
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	server.GET("/", func(c echo.Context) error {
		round := game.CurrentRound
		return c.JSON(http.StatusAccepted, map[string]any{"question": round.Question, "options": round.Options})
	})

	server.GET("/next-round", func(c echo.Context) error {
		game.CurrentRound = data.GetRandomRound()
		round := game.CurrentRound
		return c.JSON(http.StatusAccepted, map[string]any{"question": round.Question, "options": round.Options})
	})

	server.POST("/", func(c echo.Context) error {
		choice := c.FormValue("choice")
		index, err := strconv.Atoi(choice)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid choice"})
		}
		game.CurrentRound.Votes[index] += 1
		output := make(map[string]int)
		for index, option := range game.CurrentRound.Options {
			output[option] = game.CurrentRound.Votes[index]
		}
		return c.JSON(http.StatusAccepted, output)
	})

	server.Logger.Fatal(server.Start(":8080"))
}
