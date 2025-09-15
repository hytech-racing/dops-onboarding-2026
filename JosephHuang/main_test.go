package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)
func TestRootHandler(t *testing.T) {
 req := httptest.NewRequest(http.MethodGet, "/", nil) // simulate request
 rr := httptest.NewRecorder() // simulate response
 rootHandler(rr, req) // call the root handler
 if status := rr.Code; status != http.StatusOK { // check status code
 t.Errorf("expected status code %d, got %d", http.StatusOK, status)
 }
 expected := "Hello World"
 if rr.Body.String() != expected { // check the response body
 t.Errorf("expected response body %q, got %q", expected, rr.Body.String())
 }
}