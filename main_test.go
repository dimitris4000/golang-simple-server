package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestHandler(t *testing.T) {
	expected := fmt.Sprintf("%s\n", version)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	VersionHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(data) != fmt.Sprintf("%s", expected) {
		t.Errorf("Expected '%s' but got '%v'", expected, string(data))
	}
}

func TestReadinessRequestHandler(t *testing.T) {
	expected := "OK\n"

	req := httptest.NewRequest(http.MethodGet, "/rediness", nil)
	w := httptest.NewRecorder()
	ReadinessHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	got, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(got) != fmt.Sprintf("%s", expected) {
		t.Errorf("Expected '%s' but got '%v'", expected, string(got))
	}
}

func TestShutdownRequestHandler(t *testing.T) {
	expected := "Shutdown initiated\n"

	req := httptest.NewRequest(http.MethodPost, "/rediness", nil)
	w := httptest.NewRecorder()
	ShutdownHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	got, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(got) != fmt.Sprintf("%s", expected) {
		t.Errorf("Expected '%s' but got '%v'", expected, string(got))
	}
}

func TestLivenessRequestHandler(t *testing.T) {
	expected := "OK"

	req := httptest.NewRequest(http.MethodGet, "/rediness", nil)
	w := httptest.NewRecorder()
	LivenessHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	got, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(got) != fmt.Sprintf("%s\n", expected) {
		t.Errorf("Expected '%s\n' but got '%v'", expected, string(got))
	}
}
