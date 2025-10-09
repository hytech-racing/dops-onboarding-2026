package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Test with actual file
	err := upload("http://localhost:3000/upload", "text.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err_1 := upload("http://localhost:3000/upload", "mcap.mcap")
	if err_1 != nil {
		fmt.Printf("Error: %v\n", err_1)
		return
	}

	err_2 := upload("http://localhost:3000/upload", "blank.txt")
	if err_2 != nil {
		fmt.Printf("Error: %v\n", err_2)
		return
	}

}

func upload(url, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	// Create a buffer for the multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create form file field
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}

	// Copy file content to the form
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("copy file content: %w", err)
	}

	// Add additional fields
	writer.WriteField("description", "Uploaded file")

	// Close the writer
	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))

	return nil
}
