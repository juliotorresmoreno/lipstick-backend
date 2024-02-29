package mmlu

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
var tablename = models.Mmlu{}.TableName()

type MMLURouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	ai := &MMLURouter{}
	r.GET("", ai.find)
	r.GET("/:id", ai.findOne)
	r.POST("", ai.create)
	r.PATCH("/:id", ai.update)
	r.DELETE("/:id", ai.delete)
}

type Mmlu struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name" validate:"required,max=100"`
	Description string     `json:"description" validate:"max=256"`
	PhotoURL    string     `json:"photo_url" validate:"url,max=1000"`
	Model       string     `json:"model" validate:"required,max=100"`
	Provider    string     `json:"provider" validate:"required,oneof=ollama"`
	CreationAt  time.Time  `json:"creation_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MmluValidationErrors struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
	Model       string `json:"model,omitempty"`
	Provider    string `json:"provider,omitempty"`
}

func (h *MMLURouter) create(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	payload := &Mmlu{}
	err = c.ShouldBind(payload)
	if err != nil {
		log.Error("Error binding payload", err)
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	conn := db.DefaultClient

	mmlu := &models.Mmlu{
		Name:        payload.Name,
		Description: payload.Description,
		PhotoURL:    payload.PhotoURL,
		Provider:    payload.Provider,
		Model:       payload.Model,
		OwnerId:     session.ID,
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
		customErrors := MmluValidationErrors{
			Name:        errorsMap["Name"],
			PhotoURL:    errorsMap["PhotoURL"],
			Description: errorsMap["Description"],
			Model:       errorsMap["Model"],
			Provider:    errorsMap["Provider"],
		}

		log.Error("Error validating user input", customErrors)
		c.JSON(http.StatusBadRequest, customErrors)
		return
	}

	tx := conn.Create(mmlu)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	result := Mmlu{}
	tx = conn.Table(tablename).Where("id = ?", mmlu.ID).First(&result)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	c.JSON(200, result)
}

func (h *MMLURouter) update(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	payload := &Mmlu{}
	err = c.ShouldBind(payload)
	if err != nil {
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	conn := db.DefaultClient

	mmlu := &models.Mmlu{
		Name:        payload.Name,
		Description: payload.Description,
		PhotoURL:    payload.PhotoURL,
		OwnerId:     session.ID,
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
		customErrors := MmluValidationErrors{
			Name:        errorsMap["Name"],
			PhotoURL:    errorsMap["PhotoURL"],
			Description: errorsMap["Description"],
		}
		c.JSON(http.StatusBadRequest, customErrors)
		return
	}

	id := c.Param("id")
	tx := conn.Where("id = ?", id).Updates(mmlu)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	result := Mmlu{}
	tx = conn.Table(tablename).Where("id = ?", id).First(&result)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	c.JSON(200, result)
}

func (h *MMLURouter) delete(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient
	id := c.Param("id")
	tx := conn.Where("id = ?", id).
		Where(&models.Mmlu{
			OwnerId: session.ID,
		}).
		Delete(&models.Mmlu{})
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})
}

func (h *MMLURouter) find(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	mmlus := &[]Mmlu{}
	conn := db.DefaultClient
	tx := conn.Where("deleted_at is null").Where(&models.Mmlu{OwnerId: session.ID}).
		Find(mmlus)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	c.JSON(200, mmlus)
}

func (h *MMLURouter) findOne(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	mmlu := &Mmlu{}
	conn := db.DefaultClient
	tx := conn.Where("deleted_at is null").
		Where("id = ?", id).
		Where(&models.Mmlu{
			OwnerId: session.ID,
		}).First(mmlu)
	if tx.Error != nil {
		log.Error(tx.Error)
		c.JSON(404, gin.H{"message": "Not found"})
		return
	}
	c.JSON(200, mmlu)
}
