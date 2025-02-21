package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	NumGPU      int       `json:"num_gpu_layers,omitempty"`
	NumThreads  int       `json:"num_threads,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Message            Message   `json:"message"`
	Done               bool      `json:"done"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

const defaultOllamaURL = "http://localhost:11434/api/chat"
const conversationFile = "conversation.txt"

func main() {
	fmt.Println("Chat with Ollama AI (type 'exit' to quit)")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nEnter your message: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		userInput = strings.TrimSpace(userInput) // Remove newline

		if userInput == "exit" {
			fmt.Println("Exiting chat. Conversation saved in", conversationFile)
			break
		}

		msg := Message{
			Role:    "user",
			Content: userInput,
		}

		req := Request{
			Model:       "llama3.2-vision",
			Stream:      false, // Tunggu sampai selesai (non-streaming)
			Messages:    []Message{msg},
			NumGPU:      50,  // Gunakan 50 layer di GPU untuk inference
			NumThreads:  8,   // Gunakan 8 thread CPU untuk inference
			Temperature: 0.7, // Kontrol kreativitas output (0.0 - 1.0)
			TopP:        0.9, // Filtering token
		}

		fmt.Println("Processing request...")

		done := make(chan bool)
		go showLoading(done)

		start := time.Now()
		resp, err := talkToOllama(defaultOllamaURL, req)
		elapsed := time.Since(start)

		done <- true
		fmt.Println()

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		responseText := resp.Message.Content
		fmt.Println("Response:", responseText)
		fmt.Printf("Completed in %v\n", elapsed)

		saveConversation(userInput, responseText)
	}
}

func talkToOllama(url string, ollamaReq Request) (*Response, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("error response from server: %s - %s", httpResp.Status, string(body))
	}

	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ollamaResp, nil
}

func showLoading(done chan bool) {
	chars := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-done:
			return
		default:
			fmt.Printf("\rProcessing... %s", chars[i%len(chars)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func saveConversation(input, output string) {
	file, err := os.OpenFile(conversationFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	logEntry := fmt.Sprintf("User: %s\nOllama: %s\n\n", input, output)
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Println("Failed to write to file:", err)
	}
}
