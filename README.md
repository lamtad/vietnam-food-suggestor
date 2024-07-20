# Vietnamese Food Recommendation API

## Overview

This project provides an API that leverages ChatGPT and Google Maps to analyze user input and suggest Vietnamese food options for tourists based on their preferences. The system consists of a backend service developed using Gin (a Go web framework) and a frontend interface built with HTML and JavaScript.

## Purpose

The purpose of this project is to help tourists find Vietnamese food recommendations based on their preferences. Users can describe their taste or food preferences, and the system will suggest nearby Vietnamese food options accordingly.

## Components

1. **Backend (API)**: Built with Go and Gin, this component:
   - Accepts user input through a POST request.
   - Uses ChatGPT to extract keywords and identify food items from the input.
   - Queries Google Maps API to find nearby Vietnamese food locations based on the identified food items.
   - Returns the results to the frontend.

2. **Frontend**: An HTML interface that allows users to submit questions and view responses from the API.

## Getting Started

### Prerequisites

- Go 1.18 or later
- API keys for OpenAI and Google Maps

### Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/your-repo/project.git
   cd project
2. **Create a .env file in the root directory with the following content**

   ```bash
   OPENAI_API_KEY=your_openai_api_key
   GOOGLE_MAPS_API_KEY=your_google_maps_api_key

 
3. **Install Dependencies**

   ```bash
   go mod tidy

4. **Run the Application**

   ```bash
   go run main.go

#The application will start on http://localhost:8888.


   
