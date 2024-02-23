package conversation

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/server/mmlu"
	"github.com/juliotorresmoreno/tana-api/utils"
)

var log = logger.SetupLogger()

type ConversationRouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	conversation := &ConversationRouter{}
	r.GET("/:id", conversation.findOne)
	r.POST("/:id", conversation.generate)
	r.POST("/:id/attach", conversation.attach)
}

type Message struct {
	Answer   string
	Response string
}

type AttachPayload struct {
	Attachment string `json:"attachment"`
}

func generateRandomFileName(prefix, suffix string) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("Error generating random file name:", err)
	}
	return prefix + fmt.Sprintf("%x", b) + suffix
}

func (h *ConversationRouter) attach(c *gin.Context) {
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

	payload := &AttachPayload{}
	if err := c.BindJSON(payload); err != nil {
		log.Error("Error binding payload", err)
		utils.Response(c, utils.StatusBadRequest)
		return
	}

	connectionID, _ := strconv.Atoi(c.Param("id"))
	connection := &models.Connection{}
	err = mmlu.FindOne(connectionID, connection)
	if err != nil {
		log.Error("Error finding connection", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	if connection.ID == 0 {
		log.Error("connection not found", connection.ID)
		utils.Response(c, utils.StatusNotFound)
		return
	}

	attachment, err := utils.ParseBase64File(payload.Attachment)
	if err != nil {
		log.Error("Error parsing attachment", err)
		utils.Response(c, err)
		return
	}

	// Decodificar el string en base64 a bytes.
	decoded, err := base64.StdEncoding.DecodeString(attachment)
	if err != nil {
		log.Error("Error decoding attachment", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	fileName := generateRandomFileName("attachment_", ".pdf")
	filePath := filepath.Join(os.TempDir(), fileName)

	// Escribir los bytes decodificados en el archivo.
	err = os.WriteFile(filePath, decoded, 0644)
	if err != nil {
		log.Error("Error writing attachment to file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing attachment to file"})
		return
	}

	outputName := generateRandomFileName("attachment_", ".txt")
	outputPath := filepath.Join(os.TempDir(), outputName)

	err = utils.PDFToText(filePath, outputPath)
	if err != nil {
		log.Error("Error converting attachment to text", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting attachment to text"})
		return
	}

	foutput, err := os.Open(outputPath)
	if err != nil {
		log.Error("Error opening attachment", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening attachment"})
		return
	}
	defer foutput.Close()

	output, err := io.ReadAll(foutput)
	if err != nil {
		log.Error("Error reading attachment", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading attachment"})
		return
	}

	body := bytes.NewBufferString("")
	json.NewEncoder(body).Encode(map[string]interface{}{
		"title":         connection.Description,
		"attachment":    string(output),
		"user_id":       session.ID,
		"connection_id": connection.ID,
	})

	var aiUrl = os.Getenv("AI_URL")
	conversation := fmt.Sprintf("conversation-%v-%v", session.ID, connection.ID)
	url := aiUrl + "/conversation/" + conversation + "/attach"
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
	if resp.Header.Get("Content-Type") != "" {
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
	}
	c.Status(resp.StatusCode)
	utils.Copy(c.Writer, resp.Body)
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

	connectionID, _ := strconv.Atoi(c.Param("id"))
	connection := &models.Connection{}
	err = mmlu.FindOne(connectionID, connection)
	if err != nil {
		log.Error("Error finding connection", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	if connection.ID == 0 {
		log.Error("connection not found", connection.ID)
		utils.Response(c, utils.StatusNotFound)
		return
	}

	body := bytes.NewBufferString("")
	json.NewEncoder(body).Encode(map[string]interface{}{
		"title":         connection.Description,
		"prompt":        payload.Prompt,
		"user_id":       session.ID,
		"connection_id": connection.ID,
	})

	var aiUrl = os.Getenv("AI_URL")
	conversation := fmt.Sprintf("conversation-%v-%v", session.ID, connection.ID)
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

	connectionID, _ := strconv.Atoi(c.Param("id"))
	connection := &models.Connection{}
	err = mmlu.FindOne(connectionID, connection)
	if err != nil {
		log.Error("Error finding connection", err)
		utils.Response(c, utils.StatusInternalServerError)
		return
	}

	if connection.ID == 0 {
		log.Error("connection not found", connection.ID)
		utils.Response(c, utils.StatusNotFound)
		return
	}

	var aiUrl = os.Getenv("AI_URL")
	conversation := fmt.Sprintf("conversation-%v-%v", session.ID, connection.ID)
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
