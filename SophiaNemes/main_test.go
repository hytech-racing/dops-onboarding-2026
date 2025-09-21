package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
func TestUploadNonPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rr := httptest.NewRecorder()
	uploadHandler(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestUploadNonMCAP(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	fileWriter.Write([]byte("hello"))
	writer.Close()
	
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	uploadHandler(rr, req)
	
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("expected status code %d, got %d", http.StatusUnsupportedMediaType, status)
	}
	
	if !strings.Contains(rr.Body.String(), "Unsupported file type") {
		t.Errorf("expected response to contain 'Unsupported file type'")
	}
}

func TestUploadMCAP(t *testing.T) {
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    fileWriter, err := writer.CreateFormFile("file", "test.mcap")
    if err != nil {
        t.Fatal(err)
    }
    fileWriter.Write([]byte("fake mcap data"))
    writer.Close()
    
    req := httptest.NewRequest(http.MethodPost, "/upload", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    rr := httptest.NewRecorder()
    uploadHandler(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("expected status code %d, got %d", http.StatusOK, status)
    }
    
    if !strings.Contains(rr.Body.String(), "test.mcap") {
        t.Errorf("expected response to contain filename 'test.mcap'")
    }
}
