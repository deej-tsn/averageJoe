package main

import (
	"os"
	"path/filepath"

	"github.com/deej-tsn/averageJoe/model"
	"github.com/deej-tsn/averageJoe/routes"
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

	gameMgr := model.NewGM()

	routesGM := routes.NewGameMgrController(gameMgr, data)

	server := echo.New()

	server.Use(middleware.Logger())

	// DEV ONLY
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	server.GET("/user", routes.GET_newPlayerUUID)

	server.GET("/active-games", routesGM.GET_activeGames)

	server.POST("/connect-to-game", routesGM.POST_connectToGame)

	server.POST("/create-game", routesGM.POST_createGame)

	server.Logger.Fatal(server.Start(":8080"))
}
