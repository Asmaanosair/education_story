package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"integration/config"
	"integration/models"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var genAIClient *genai.Client

func LoadKey() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyD1SuELh1-nOMsWZbXWs6WlIo-ny8rw9VQ"))
	if err != nil {
		log.Fatalf("Failed to create GenAI client: %v", err)
	}
	genAIClient = client
}

func GenerateStory(c *gin.Context) {
	model := genAIClient.GenerativeModel("gemini-1.5-flash")
	prompt := "Generate an interesting and factual short story about a random topic."
	res, err := model.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate story"})
		return
	}

	storyText := extractResponse(res)
	story := models.Story{Text: storyText}
	config.DB.Create(&story)
	go func(storyID uint) {
		// time.Sleep(5 * time.Minute)
		GenerateQuizForStory(storyID)
	}(story.ID)

	c.JSON(http.StatusOK, gin.H{"story": storyText})
}

func GenerateQuizForStory(storyID uint) {
	var story models.Story
	config.DB.First(&story, storyID)

	model := genAIClient.GenerativeModel("gemini-1.5-flash")
	prompt := "Generate a quiz question based on this story:\n\n" + story.Text
	res, err := model.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		log.Println("Failed to generate quiz question")
		return
	}

	story.Quiz = extractResponse(res)
	config.DB.Save(&story)
	log.Println("Quiz generated for story:", story.ID)
}

func SubmitAnswer(c *gin.Context) {
	var request struct {
		StoryID uint   `json:"story_id"`
		Answer  string `json:"answer"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var story models.Story
	config.DB.First(&story, request.StoryID)

	model := genAIClient.GenerativeModel("gemini-1.5-flash")
	prompt := fmt.Sprintf("Evaluate this answer: \"%s\" for the question: \"%s\".", request.Answer, story.Quiz)
	res, err := model.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate answer"})
		return
	}

	score := extractScore(res)
	story.Answer = request.Answer
	story.Score = score
	config.DB.Save(&story)

	c.JSON(http.StatusOK, gin.H{"score": score})
}

func DownloadScores(c *gin.Context) {
	filePath := "scores.csv"
	file, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Question", "Answer", "Score"})

	var stories []models.Story
	config.DB.Find(&stories)
	for _, story := range stories {
		writer.Write([]string{
			strconv.Itoa(int(story.ID)),
			story.Quiz,
			story.Answer,
			strconv.Itoa(story.Score),
		})
	}

	c.File(filePath)
}

func extractResponse(res *genai.GenerateContentResponse) string {
	if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])
	}
	return "No response"
}

func extractScore(res *genai.GenerateContentResponse) int {
	responseText := extractResponse(res)
	score, err := strconv.Atoi(responseText)
	if err != nil {
		return 0
	}
	return score
}
