package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	GeneratedText string `json:"generated_text"`
}

func GeneratedTextonly(prompt string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	// Ambil konten yang dihasilkan
	var generatedText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				// Debug print part secara keseluruhan
				fmt.Printf("Part: %+v\n", part) // Debug print part untuk melihat semua field

				// Coba akses konten yang tepat berdasarkan struktur part
				// Jika tidak ada field Content, kita harus menyesuaikan cara aksesnya
				if part != nil {
					// Sesuaikan dengan cara kita membaca data, misalnya `part` mungkin memiliki field lain
					generatedText += fmt.Sprintf("%+v\n", part) // Debug print part sebagai string
				}
			}
		}
	}

	// Pastikan generatedText diisi dengan konten yang benar
	fmt.Println("Generated Text:", generatedText)

	return generatedText, nil
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	var request PromptRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate text from prompt
	generatedText, err := GeneratedTextonly(request.Prompt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating content: %v", err), http.StatusInternalServerError)
		return
	}

	// Create response
	response := Response{GeneratedText: generatedText}

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Muat file .env sebelum mengakses variabel lingkungan
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Register the handler for the /generate route
	http.HandleFunc("/generate", generateHandler)

	// Start the HTTP server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
