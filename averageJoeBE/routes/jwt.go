package routes

import (
	"net/http"

	"github.com/deej-tsn/averageJoe/config"
	"github.com/deej-tsn/averageJoe/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type JWTController struct {
	cfg *config.Config
}

type userReq struct {
	Username string `json:"username"`
}

type verifyReq struct {
	Token string `json:"token"`
}

func NewJWTController(cfg *config.Config) *JWTController {
	return &JWTController{
		cfg: cfg,
	}
}

// POST (/user)
func (jc *JWTController) POST_user(c echo.Context) error {
	var user userReq
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorMessage("no username provided"))
	}
	token, err := util.CreateToken(user.Username, jc.cfg.JWTSecret)
	if err != nil {
		return c.JSON(http.StatusRequestTimeout, util.ErrorMessage("token could not be produced"))
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (jc *JWTController) POST_verifyUser(c echo.Context) error {
	var verReq verifyReq
	if err := c.Bind(&verReq); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorMessage("No Token Provided"))
	}
	user, err := util.VerifyToken(verReq.Token, jc.cfg.JWTSecret)
	if err != nil {
		log.Printf("error from verify token: %v", err)
		return c.JSON(http.StatusBadRequest, util.ErrorMessage("Invalid Token Provided"))
	}
	return c.JSON(http.StatusOK, user)

}
