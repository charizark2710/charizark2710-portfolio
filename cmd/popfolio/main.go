package main

import (
	"log"
	"path/filepath"

	"popfolio/internal/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Use correct relative path for templates from cmd/popfolio
	templatePath := filepath.Join("..", "..", "templates", "*.html")
	router.LoadHTMLGlob(templatePath)
	router.Static("/static", filepath.Join("..", "..", "static"))

	routers.SetupRoutes(router)

	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
