package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-chi/chi/v5"
)

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil) // simulate request
	rr := httptest.NewRecorder()                         // simulate response
	rootHandler(rr, req)                                 // call the root handler
	if status := rr.Code; status != http.StatusOK {      // check status code
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}
	expected := "Hello World"
	if rr.Body.String() != expected { // check the response body
		t.Errorf("expected response body %q, got %q", expected, rr.Body.String())
	}
}

func newFileUploadRequest(uri, paramName, fileName string, fileContent []byte) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, err
	}
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestUploadNonPOST(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/upload", uploadHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 Method Not Allowed, got %d", rr.Code)
	}
}

func TestUploadNonMCAP(t *testing.T) {
	req, err := newFileUploadRequest("/upload", "file", "test.txt", []byte("hello"))
	require.Nil(t, err)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/upload", uploadHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected 415 Unsupported Media Type, got %d", rr.Code)
	}

	var resp McapError
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Err != "Unsupported file type" {
		t.Errorf("unexpected error: %s", resp.Err)
	}
}

func TestUploadMCAP(t *testing.T) {
	fileContent := []byte("dummy MCAP content")
	req, err := newFileUploadRequest("/upload", "file", "test.mcap", fileContent)
	require.Nil(t, err)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/upload", uploadHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var resp McapSuccess

	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != "MCAP uploaded successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
	if resp.File.Name != "test.mcap" {
		t.Errorf("unexpected file name: %s", resp.File.Name)
	}
	expected := int32(len(fileContent))
	if resp.File.Size != expected {
		t.Errorf("unexpected file size: %d", resp.File.Size)
	}
}
