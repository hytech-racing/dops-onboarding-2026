package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper to create a file upload request
func newFileUploadRequest(method, target, filename string, content []byte) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", filename)
	part.Write(content)
	writer.Close()

	req := httptest.NewRequest(method, target, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUploadHandler_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestUploadHandler_InvalidFileType(t *testing.T) {
	req := newFileUploadRequest(http.MethodPost, "/upload", "file.png", []byte("fake content"))
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("expected 415, got %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Unsupported file type") {
		t.Errorf("expected error message, got %q", body)
	}
}

func TestUploadHandler_ValidMCAP(t *testing.T) {
	req := newFileUploadRequest(http.MethodPost, "/upload", "test.mcap", []byte("fake content"))
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "MCAP uploaded successfully") {
		t.Errorf("expected success message, got %q", body)
	}
	if !strings.Contains(body, "test.mcap") {
		t.Errorf("expected filename in response, got %q", body)
	}
}
