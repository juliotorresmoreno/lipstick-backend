package credentials

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/utils"
)

var maxCredentials = 10
var tablename = models.Credential{}.TableName()
var log = logger.SetupLogger()

type CredentialsRouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	documents := &CredentialsRouter{}
	r.GET("", documents.find)
	r.GET("/:id", documents.findOne)
	r.POST("/generate", documents.create)
	r.DELETE("/:id", documents.delete)
}

type Credential struct {
	ID         uint       `json:"id"`
	ApiKey     string     `json:"api_key"`
	ApiSecret  string     `json:"api_secret"`
	LastUsed   *time.Time `json:"last_used"`
	CreationAt time.Time  `json:"creation_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func (h *CredentialsRouter) find(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient

	credentials := make([]Credential, 0)
	tx := conn.Table(tablename).
		Where("user_id = ?", session.ID).
		Where("deleted_at is null").
		Limit(maxCredentials).
		Find(&credentials)
	if tx.Error != nil {
		log.Error("Error getting credentials", tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, credentials)
}

func (h *CredentialsRouter) findOne(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient

	credential := &Credential{}
	tx := conn.Table(tablename).
		Where("id = ?", c.Param("id")).
		Where("user_id = ?", session.ID).
		Where("deleted_at is null").
		First(credential)
	if tx.Error != nil {
		log.Error("Error getting credentials", tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, credential)
}

func (h *CredentialsRouter) create(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient

	count := int64(0)
	tx := conn.Table(tablename).
		Where("user_id = ?", session.ID).
		Where("deleted_at is null").
		Count(&count)
	if tx.Error != nil {
		log.Error("Error getting credentials", tx.Error)
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	if int(count) > maxCredentials {
		c.JSON(400, gin.H{"message": "You can't create more than 5 credentials"})
		return
	}

	apiSecret, _ := utils.GenerateRandomString(64)

	credential := &models.Credential{
		UserId:    session.ID,
		ApiKey:    uuid.New().String(),
		ApiSecret: apiSecret,
	}
	tx = conn.Create(credential)
	if tx.Error != nil {
		log.Error("Error creating credential", tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "create credentials"})
}

func (h *CredentialsRouter) delete(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	conn := db.DefaultClient

	tx := conn.Where("user_id = ?", session.ID).Delete(&models.Credential{}, c.Param("id"))
	if tx.Error != nil {
		log.Error("Error deleting credential", tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "delete credentials"})
}
