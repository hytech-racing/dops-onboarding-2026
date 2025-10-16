package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(UploadMcap))
	defer ts.Close()

	type want struct {
		status  int
		hasFile bool
		name    string
		size    int
	}

	tests := []struct {
		filename string
		want     want
	}{
		{
			filename: "test.mcap",
			want:     want{status: http.StatusOK, hasFile: true, name: "test.mcap", size: 5120},
		},
		{
			filename: "test.txt",
			want:     want{status: http.StatusUnsupportedMediaType, hasFile: false},
		},
	}

	upload := func(filename string) (*http.Response, error) {
		body := &bytes.Buffer{}
		w := multipart.NewWriter(body)

		part, err := w.CreateFormFile("file", filename)
		if err != nil {
			return nil, err
		}
		if _, err := io.CopyN(part, bytes.NewReader(make([]byte, 5120)), 5120); err != nil {
			return nil, err
		}
		if err := w.Close(); err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, ts.URL, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		return http.DefaultClient.Do(req)
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			resp, err := upload(tc.filename)
			if err != nil {
				t.Fatalf("upload error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.want.status {
				t.Fatalf("status: want %d got %d", tc.want.status, resp.StatusCode)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil && err != io.EOF {
				t.Fatalf("decode error: %v", err)
			}

			fileVal, hasFile := result["file"]
			if tc.want.hasFile != hasFile {
				t.Fatalf("file presence mismatch: want %v got %v (payload: %#v)", tc.want.hasFile, hasFile, result)
			}

			if tc.want.hasFile {
				m, ok := fileVal.(map[string]interface{})
				if !ok || m == nil {
					t.Fatalf("file field not an object: %#v", fileVal)
				}

				name, _ := m["name"].(string)
				if name != tc.want.name {
					t.Fatalf("name: want %q got %q", tc.want.name, name)
				}

				sizeF, _ := m["size"].(float64)
				if int(sizeF) != tc.want.size {
					t.Fatalf("size: want %d got %d", tc.want.size, int(sizeF))
				}
			} else {

				if errMsg, ok := result["error"].(string); !ok || errMsg == "" {
					t.Logf("no 'error' field in JSON (payload: %#v); this is fine if handler returns plain text", result)
				}
			}
		})
	}
}
