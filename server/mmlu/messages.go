package mmlu

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/utils"
)

type MessageValidationErrors struct {
	Content string `json:"content,omitempty"`
}

type Message struct {
	ID         uint       `json:"id"`
	Content    string     `json:"content" validate:"required"`
	MmluId     uint       `json:"-"`
	Mmlu       Mmlu       `json:"mmlu"`
	Role       string     `json:"role"`
	CreationAt time.Time  `json:"creationAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

func (h *MMLURouter) findMessages(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	mmluId, _ := strconv.Atoi(c.Param("id"))
	messages := &[]Message{}
	conn := db.DefaultClient
	tx := conn.Preload("Mmlu").
		Model(models.Message{}).
		Where("deleted_at is null").
		Where(&models.Message{OwnerId: session.ID, MmluId: uint(mmluId)}).
		Find(messages)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	c.JSON(200, messages)
}

type createMessagePayload struct {
	Content string `json:"content" validate:"required"`
}

func (h *MMLURouter) createMessage(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	mmluId, _ := strconv.Atoi(c.Param("id"))
	payload := &createMessagePayload{}
	if err := c.ShouldBind(payload); err != nil {
		log.Error(err)
		utils.Response(c, utils.StatusBadRequest)
		return
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
		customErrors := MessageValidationErrors{
			Content: errorsMap["content"],
		}

		log.Error("Error validating user input", customErrors)
		c.JSON(http.StatusBadRequest, customErrors)
		return
	}

	message := &models.Message{
		Content: payload.Content,
		MmluId:  uint(mmluId),
		OwnerId: session.ID,
		Role:    "system",
	}
	conn := db.DefaultClient
	tx := conn.Create(message)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "create success"})
}

type AttachPayload struct {
	Attachment string `json:"attachment"`
}

func (h *MMLURouter) attachMessage(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		utils.Response(c, err)
		return
	}

	payload := &AttachPayload{}
	if err := c.BindJSON(payload); err != nil {
		log.Error("Error binding payload", err)
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	mmluID, _ := strconv.Atoi(c.Param("id"))
	mmlu := &models.Mmlu{}
	conn := db.DefaultClient

	tx := conn.Where(&models.Connection{
		OwnerId: session.ID,
	}).First(mmlu, mmluID)

	if tx.Error != nil {
		log.Error("Error finding connection", tx.Error)
		utils.Response(c, utils.StatusNotFound)
		return
	}

	if mmlu.ID == 0 {
		log.Error("mmlu not found", mmlu.ID)
		utils.Response(c, utils.StatusNotFound)
		return
	}

	content, err := utils.ReadPDF(payload.Attachment)
	if err != nil {
		log.Error("Error reading attachment", err)
		utils.Response(c, err)
		return
	}

	message := &models.Message{
		Content: content,
		MmluId:  uint(mmlu.ID),
		OwnerId: session.ID,
		Role:    "system",
	}

	tx = conn.Create(message)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "create success"})
}

type updateMessagePayload struct {
	Content string `json:"content" validate:"required"`
}

func (h *MMLURouter) updateMessage(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	payload := &updateMessagePayload{}
	if err := c.ShouldBind(payload); err != nil {
		log.Error(err)
		utils.Response(c, utils.StatusBadRequest)
		return
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
		customErrors := MessageValidationErrors{
			Content: errorsMap["content"],
		}

		log.Error("Error validating user input", customErrors)
		c.JSON(http.StatusBadRequest, customErrors)
		return
	}

	messageId, _ := strconv.Atoi(c.Param("messageId"))
	message := &models.Message{
		Content: payload.Content,
	}
	conn := db.DefaultClient
	tx := conn.Model(message).
		Where(&models.Message{OwnerId: session.ID}).
		Where("id = ?", messageId).
		Updates(message)
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "update success"})
}

func (h *MMLURouter) deleteMessage(c *gin.Context) {
	session, err := utils.ValidateSession(c)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	messageId, _ := strconv.Atoi(c.Param("messageId"))
	conn := db.DefaultClient
	tx := conn.Model(models.Message{}).
		Where(&models.Message{OwnerId: session.ID}).
		Where("id = ?", messageId).
		Update("deleted_at", time.Now())
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}
