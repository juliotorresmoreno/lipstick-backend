package connections

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/utils"
)

var log = logger.SetupLogger()
var tablename = models.Connection{}.TableName()

type ConnectionsRouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	connections := &ConnectionsRouter{}
	r.GET("", connections.find)
	r.GET("/:id", connections.findOne)
	r.POST("", connections.create)
	r.PATCH("/:id", connections.update)
	r.DELETE("/:id", connections.delete)
}

type Mmlu struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	PhotoURL    string     `json:"photo_url"`
	CreationAt  time.Time  `json:"creation_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Connection struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name" validate:"required,max=100"`
	Description string     `json:"description" validate:"max=256"`
	PhotoURL    string     `json:"photo_url" validate:"url,max=1000"`
	MmluId      uint       `json:"mmlu_id"`
	Mmlu        Mmlu       `json:"mmlu"`
	CreationAt  time.Time  `json:"creation_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func (h *ConnectionsRouter) findOne(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient
	connection := &Connection{}

	tx := conn.Table(tablename).
		Where("id = ?", c.Param("id")).
		Where(&models.Connection{
			OwnerId: session.ID,
		}).First(connection)
	if tx.Error != nil {
		log.Error(tx.Error)
		c.JSON(404, gin.H{"message": "Not found"})
		return
	}

	c.JSON(200, connection)
}

func (h *ConnectionsRouter) find(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient
	connections := &[]Connection{}

	tx := conn.Table(tablename).
		Where(&models.Connection{
			OwnerId: session.ID,
		}).Find(connections)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, connections)
}

type ConnectionValidationErrors struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
}

func (h *ConnectionsRouter) create(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient
	payload := &Connection{}

	err = c.ShouldBind(payload)
	if err != nil {
		log.Error("Error binding payload", err)
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	connection := &models.Connection{
		Name:        payload.Name,
		Description: payload.Description,
		PhotoURL:    payload.PhotoURL,
		OwnerId:     session.ID,
		MmluId:      payload.MmluId,
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		log.Error("Error validating user input", err)
		errorsMap := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()

			switch tag {
			case "required":
				errorsMap[field] = "This field is required!"
			case "email":
				errorsMap[field] = "Invalid email format!"
			case "phone":
				errorsMap[field] = "Invalid phone number!"
			default:
				errorsMap[field] = "Invalid field!"
			}
		}
		customErrors := ConnectionValidationErrors{
			Name:        errorsMap["Name"],
			PhotoURL:    errorsMap["PhotoURL"],
			Description: errorsMap["Description"],
		}

		log.Error("Error validating user input", customErrors)
		c.JSON(http.StatusBadRequest, customErrors)
		return
	}

	tx := conn.Create(connection)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "create connection"})
}

func (h *ConnectionsRouter) update(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient
	payload := &Connection{}

	err = c.ShouldBind(payload)
	if err != nil {
		log.Error("Error binding payload", err)
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	connection := &models.Connection{
		Name:        payload.Name,
		Description: payload.Description,
		PhotoURL:    payload.PhotoURL,
		OwnerId:     session.ID,
		MmluId:      payload.MmluId,
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		log.Error("Error validating user input", err)
		errorsMap := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()

			switch tag {
			case "required":
				errorsMap[field] = "This field is required!"
			case "email":
				errorsMap[field] = "Invalid email format!"
			case "phone":
				errorsMap[field] = "Invalid phone number!"
			default:
				errorsMap[field] = "Invalid field!"
			}
		}
		customErrors := ConnectionValidationErrors{
			Name:        errorsMap["Name"],
			PhotoURL:    errorsMap["PhotoURL"],
			Description: errorsMap["Description"],
		}

		log.Error("Error validating user input", customErrors)
		c.JSON(http.StatusBadRequest, customErrors)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	tx := conn.Where(&models.Connection{
		ID:      uint(id),
		OwnerId: session.ID,
	}).Updates(connection)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "update connection"})
}

func (h *ConnectionsRouter) delete(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient
	id, _ := strconv.Atoi(c.Param("id"))
	tx := conn.Where(&models.Connection{
		OwnerId: session.ID,
	}).Delete(&models.Connection{}, uint(id))
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "delete connection"})
}
