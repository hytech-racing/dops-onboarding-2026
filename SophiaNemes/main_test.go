package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil) // simulate request
	rr := httptest.NewRecorder()                         // simulate response
	rootHandler(rr, req)                                 // call the root handler
	require.Equal(t, http.StatusOK, rr.Code)            // use require instead of manual if check
	expected := "Hello World"
	require.Equal(t, expected, rr.Body.String()) // use require for cleaner assertions
}

func TestUploadNonPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUploadNonMCAP(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	require.Nil(t, err)

	testFileContent := []byte("hello")
	_, err = fileWriter.Write(testFileContent)
	require.Nil(t, err)

	err = writer.Close()
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)

	 
	resp := McapError{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.Nil(t, err)
	require.Equal(t, "Unsupported file type", resp.Err)
}

func TestUploadMCAP(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileWriter, err := writer.CreateFormFile("file", "test.mcap")
	require.Nil(t, err)

	fileContent := []byte("fake mcap data")
	_, err = fileWriter.Write(fileContent)
	require.Nil(t, err)

	err = writer.Close()
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	
	resp := McapSuccess{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.Nil(t, err)
	require.Equal(t, "test.mcap", resp.File.Name)
	require.Equal(t, int64(len(fileContent)), resp.File.Size)
	require.Equal(t, "MCAP uploaded successfully", resp.Message)
}