package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/deej-tsn/averageJoe/config"
	"github.com/deej-tsn/averageJoe/model"
	"github.com/deej-tsn/averageJoe/routes"
	"github.com/deej-tsn/averageJoe/util"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	cfg := config.LoadConfig()

	data_path, _ := filepath.Abs("../data/questions.json")

	data_bytes, err := os.ReadFile(data_path)
	if err != nil {
		panic(err)
	}

	data := model.LoadData(data_bytes)

	gameMgr := model.NewGM()

	routesGM := routes.NewGameMgrController(gameMgr, data)
	routesJWT := routes.NewJWTController(cfg)

	server := echo.New()

	server.Use(middleware.Logger())

	// DEV ONLY
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	restricted := server.Group("/games")

	// USER LOGIN

	server.POST("/user", routesJWT.POST_user)

	server.POST("/verify-user", routesJWT.POST_verifyUser)

	restricted_config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(util.JWT_CustomClaim)
		},
		SigningKey: cfg.JWTSecret,
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusForbidden, util.ErrorMessage("Invalid JWT token"))
		},
	}
	restricted.Use(util.ExtractWebsocketToken)
	restricted.Use(echojwt.WithConfig(restricted_config))

	// GAME ROUTES

	restricted.GET("/connect-to-game", routesGM.WS_handler)

	restricted.POST("/create-game", routesGM.POST_createGame)

	restricted.GET("/active-games", routesGM.GET_activeGames)

	server.Logger.Fatal(server.Start(":8080"))
}
