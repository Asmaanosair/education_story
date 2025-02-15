package routes

import (
	"integration/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/generate-story", controllers.GenerateStory)
	router.POST("/submit-answer", controllers.SubmitAnswer)
	router.GET("/download-scores", controllers.DownloadScores)

	return router
}
