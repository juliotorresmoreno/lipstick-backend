package auth

import (
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/utils"
)

type AuthRouter struct {
}

func SetupAUTHRoutes(r *gin.RouterGroup) {
	auth := &AuthRouter{}

	r.GET("", auth.Ping)
	r.GET("/session", auth.Session)
	r.POST("/sign-in", auth.SignIn)
	r.POST("/sign-up", auth.SignUp)
}

type SignUpPayload struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

var signUpValidator = NewSignUpValidator()

func (auth *AuthRouter) SignUp(c *gin.Context) {
	payload := &SignUpPayload{}

	err := c.Bind(payload)
	if err != nil {
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	validation, err := signUpValidator.ValidateSignUp(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	payload.Password, err = utils.HashPassword(payload.Password)
	if err != nil {
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	conn := db.DefaultClient

	user := &models.User{
		Name:     payload.Name,
		LastName: payload.LastName,
		Phone:    payload.Phone,
		Email:    payload.Email,
		Password: payload.Password,
	}
	tx := conn.Save(user)
	if tx.Error != nil {
		log.Info(tx.Error)

		if strings.Contains(tx.Error.Error(), "duplicate key") {
			c.JSON(http.StatusBadRequest, gin.H{
				"email": payload.Email + " already exists",
			})
			return
		}
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	session, err := utils.MakeSession(user)
	if err != nil {
		utils.Response(c, err)
	}
	c.JSON(200, session)
}

type SignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (auth *AuthRouter) SignIn(c *gin.Context) {
	payload := &SignInPayload{}

	err := c.Bind(payload)
	if err != nil {
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	conn := db.DefaultClient
	user := &models.User{}

	tx := conn.Select(utils.SessionFields, "password").First(
		user, "email = ?", payload.Email,
	)
	if tx.Error != nil {
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	ok, err := utils.ComparePassword(payload.Password, user.Password)
	if !ok || err != nil {
		utils.Response(c, utils.StatusUnauthorized)
		return
	}
	user.Password = ""

	session, err := utils.MakeSession(user)
	if err != nil {
		utils.Response(c, err)
	}
	c.JSON(200, session)
}

func (auth *AuthRouter) Session(c *gin.Context) {
	token, _ := c.Cookie("token")
	if token == "" {
		token = c.Request.URL.Query().Get("token")
	}

	if token == "" {
		token = c.Request.Header.Get("authorization")
	}
	if len(token) > 7 && strings.ToLower(token[:6]) == "bearer" {
		token = token[7:]
	}

	if token == "" {
		log.Error("utils.StatusUnauthorized")
		utils.Response(c, utils.StatusUnauthorized)
		return
	}

	session, err := utils.ValidateSession(token)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, session)
}

func (auth *AuthRouter) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
