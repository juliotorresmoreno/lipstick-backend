package conversation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/utils"
)

var log = logger.SetupLogger()

type ConversationRouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	conversation := &ConversationRouter{}
	r.GET("/:id", conversation.findOne)
	r.POST("/:id", conversation.generate)
}

type Message struct {
	Answer   string
	Response string
}

type GeneratePayload struct {
	Prompt string `json:"prompt"`
}

func (h *ConversationRouter) generate(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		utils.Response(c, err)
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		utils.Response(c, err)
		return
	}

	payload := &GeneratePayload{}
	if err := c.BindJSON(payload); err != nil {
		log.Error("Error binding payload", err)
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	mmlu := &models.Mmlu{}
	conn := db.DefaultClient
	if tx := conn.Find(mmlu); tx.Error != nil {
		log.Error("Error finding mmlu", tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	if mmlu.ID == 0 {
		log.Error("Mmlu not found", mmlu.ID)
		utils.Response(c, utils.StatusNotFound)
		return
	}

	body := bytes.NewBufferString("")
	json.NewEncoder(body).Encode(map[string]interface{}{
		"title":  mmlu.Description,
		"prompt": payload.Prompt,
	})

	var aiUrl = os.Getenv("AI_URL")
	conversation := fmt.Sprintf("conversation-%v-%v", session.ID, mmlu.ID)
	url := aiUrl + "/conversation/" + conversation
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error("Error creating request", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Error making request", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	c.Header("Content-Type", "text/plain")
	c.Status(http.StatusOK)
	utils.Copy(c.Writer, resp.Body)
}

func (h *ConversationRouter) findOne(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		utils.Response(c, err)
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		utils.Response(c, err)
		return
	}

	mmlu := &models.Mmlu{}
	conn := db.DefaultClient
	if tx := conn.Find(mmlu); tx.Error != nil {
		log.Error("Error finding mmlu", tx.Error)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	if mmlu.ID == 0 {
		log.Error("Mmlu not found", mmlu.ID)
		utils.Response(c, utils.StatusNotFound)
		return
	}

	var aiUrl = os.Getenv("AI_URL")
	conversation := fmt.Sprintf("conversation-%v-%v", session.ID, mmlu.ID)
	url := aiUrl + "/conversation/" + conversation
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Error creating request", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Error making request", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	_, err = io.Copy(c.Writer, res.Body)
	if err != nil {
		log.Error("Error copying response", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}
}
