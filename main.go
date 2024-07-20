package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var (
	openaiAPIKey     = os.Getenv("OPENAI_API_KEY")
	googleMapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
)

func main() {
	if openaiAPIKey == "" || googleMapsAPIKey == "" {
		log.Fatal("API keys are not set")
	}

	r := gin.Default()

	r.POST("/chat", func(c *gin.Context) {
		var requestBody map[string]string
		if err := c.BindJSON(&requestBody); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		question := requestBody["question"]

		// Step 1: Extract keywords using ChatGPT
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+openaiAPIKey).
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"model": "gpt-3.5-turbo",
				"messages": []map[string]string{
					{"role": "user", "content": "Please ignore all previous instructions. Please respond only in the English language. Do not explain what you are doing. Do not self-reference. You are an expert text analyst and researcher. First, extract the relevant keywords from the provided text. Then, determine a food item that matches the extracted keywords. Present the results in a markdown table with the name of the food item in a list, list as many as possible. And here is the text: \"" + question + "\""},
				},
			}).
			Post("https://api.openai.com/v1/chat/completions")

		if err != nil {
			log.Printf("Error calling OpenAI API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			log.Printf("Error unmarshalling OpenAI response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		choices, ok := result["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			log.Printf("No choices found in OpenAI response")
			c.JSON(http.StatusOK, gin.H{"response": "No response from model"})
			return
		}

		message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		if !ok {
			log.Printf("Unexpected response format from OpenAI")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected response format"})
			return
		}

		// Extract food items from the response
		var foodItems []string
		foodItems = strings.Split(message, "\n")

		// Step 2: Search for places using Google Maps API
		var places []string
		for _, item := range foodItems {
			mapResp, err := client.R().
				SetQueryParams(map[string]string{
					"query": item + " near me",
					"key":   googleMapsAPIKey,
				}).
				Get("https://maps.googleapis.com/maps/api/place/textsearch/json")

			if err != nil {
				log.Printf("Error calling Google Maps API: %v", err)
				continue
			}

			var mapResult map[string]interface{}
			if err := json.Unmarshal(mapResp.Body(), &mapResult); err != nil {
				log.Printf("Error unmarshalling Google Maps response: %v", err)
				continue
			}

			placesResult, ok := mapResult["results"].([]interface{})
			if !ok || len(placesResult) == 0 {
				log.Printf("No places found for %s", item)
				continue
			}

			for _, place := range placesResult {
				placeMap, ok := place.(map[string]interface{})
				if !ok {
					continue
				}
				name, _ := placeMap["name"].(string)
				address, _ := placeMap["formatted_address"].(string)
				places = append(places, name+" - "+address)
			}
		}

		if len(places) == 0 {
			c.JSON(http.StatusOK, gin.H{"response": "No places found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": strings.Join(places, "\n")})
	})

	r.StaticFile("/", "./index.html")

	r.Run(":8888")
}
