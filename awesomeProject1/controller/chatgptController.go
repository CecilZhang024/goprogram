package controller

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func PostToGPT(c *gin.Context) {
	gptApiKey := os.Getenv("CHATGPT_API_KEY")

	url := "https://api.openai.com/v1/chat/completions"
	contentType := "application/json"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"badrequest": "cann't connect to chatgpt"})
		return
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+gptApiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": string(body)})
}
