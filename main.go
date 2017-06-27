package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	backendIP := os.Getenv("SERVER_IP")
	//backendDB := os.Getenv("DB_URL")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello peer")
	})

	router.GET("/pass/:port/:word", func(c *gin.Context) {
		if backendIP == "" {
			c.String(http.StatusNotFound, "No backend configured")
			return
		}
		resp, err := http.Get("http://" + backendIP + ":" + c.Param("port") + "/backend/" + c.Param("word"))
		if err != nil {
			c.String(http.StatusInternalServerError, "Backend error: "+err.Error())
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "problem reading response from backend: "+err.Error())
			return
		}
		c.String(resp.StatusCode, "Backend said:\n"+string(body)+"\n")
	})

	// This endpoint is used when acting as the backend server
	router.GET("/backend/:word", func(c *gin.Context) {
		c.String(http.StatusOK, "You say "+c.Param("word")+"\nI say hello\n")
	})

	router.Run(":" + port)
}
