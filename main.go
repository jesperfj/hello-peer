package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	backendIP := os.Getenv("SERVER_IP")
	backendDB := os.Getenv("DB_URL")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	//router.LoadHTMLGlob("templates/*.tmpl.html")
	//router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello peer")
	})

	router.GET("/pass/:port/:word", func(c *gin.Context) {
		if backendIP == "" {
			c.String(http.NotFound, "No backend configured")
			return
		}
		resp, err := http.Get("http://" + backendIP + ":" + c.Param("port") + "/backend/" + c.Param("word"))
		if err != nil {
			c.String(http.StatusInternalServerError, "Backend error: "+err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "problem reading response from backend: "+err)
			return
		}
		c.String(resp.StatusCode, "Backend said: "+body)
	})

	// This endpoint is used when acting as the backend server
	router.GET("/backend/:word", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from the backend")
	})

	router.Run(":" + port)
}
