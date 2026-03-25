package routers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"popfolio/internal/model"
	"popfolio/internal/storage"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	portfolioData, err := storage.LoadCSVData()
	if err != nil {
		log.Printf("Error loading portfolio data: %v", err)
		portfolioData = &model.PortfolioData{
			Name:  "Portfolio",
			Title: "Developer",
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"data": portfolioData,
		})
	})

	router.GET("/api/work-experience", func(c *gin.Context) {
		// Convert each Summary to template.HTML so the HTML is rendered, not escaped
		type SafeWorkExp struct {
			Position string
			Company  string
			Duration string
			Summary  template.HTML
			GitHub   string
		}

		safeWorkExp := make([]SafeWorkExp, len(portfolioData.WorkExp.Details))
		for i, w := range portfolioData.WorkExp.Details {
			safeWorkExp[i] = SafeWorkExp{
				Position: w.Position,
				Company:  w.Company,
				Duration: w.Duration,
				Summary:  template.HTML(w.Summary),
				GitHub:   w.GitHub,
			}
		}

		c.HTML(http.StatusOK, "work-experience.html", gin.H{
			"workExp": safeWorkExp,
		})
	})

	router.GET("/api/education", func(c *gin.Context) {
		c.HTML(http.StatusOK, "education.html", gin.H{
			"education": portfolioData.Education.Details,
		})
	})

	router.GET("/api/work-experience/:id", func(c *gin.Context) {
		id := c.Param("id")
		index := 0
		_, err := fmt.Sscanf(id, "%d", &index)
		if err != nil || index < 0 || index >= len(portfolioData.WorkExp.Details) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Work experience not found",
			})
			return
		}

		w := portfolioData.WorkExp.Details[index]
		c.HTML(http.StatusOK, "work-detail.html", gin.H{
			"work": gin.H{
				"Position":    w.Position,
				"Company":     w.Company,
				"Duration":    w.Duration,
				"Description": template.HTML(w.Description),
				"PreviewFile": w.PreviewFile,
				"GitHub":      w.GitHub,
			},
		})
	})

	router.GET("/api/education/:id", func(c *gin.Context) {
		id := c.Param("id")
		index := 0
		_, err := fmt.Sscanf(id, "%d", &index)
		if err != nil || index < 0 || index >= len(portfolioData.Education.Details) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Education not found",
			})
			return
		}

		edu := portfolioData.Education.Details[index]
		c.HTML(http.StatusOK, "education-detail.html", gin.H{
			"education": edu,
		})
	})

	router.GET("/preview/:filename", func(c *gin.Context) {
		filename := filepath.Clean(c.Param("filename"))
		if strings.Contains(filename, "..") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid filename",
			})
			return
		}

		fp := filepath.Join("previews", filename)
		c.File(fp)
	})
}
