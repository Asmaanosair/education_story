 Go Quiz App - Docker & Postman Guide

This project uses Docker Compose to manage services and Postman to test the API.

🛠 Running the Project with Docker Compose

Ensure Docker and Docker Compose are installed

Build and start the containers

docker-compose up --build


 API Testing with Postman

Generate a Story

POST http://localhost:8080/generate-story

Body:

{ "story": "In ancient Greece..." }

Submit an Answer

POST http://localhost:8080/submit-answer

Body:

{ "story_id": 1, "answer": "The Greeks believed..." }

Download Scores

GET http://localhost:8080/download-scores
