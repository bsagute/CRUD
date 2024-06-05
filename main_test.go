package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	initConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/post", postHandler)

	server = httptest.NewServer(mux)
	defer server.Close()

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestConfig(t *testing.T) {
	viper.SetEnvPrefix("test")
	viper.SetDefault("server.port", "8080")
	viper.AutomaticEnv()

	expectedPort := "8080"
	if viper.GetString("server.port") != expectedPort {
		t.Errorf("expected port %s but got %s", expectedPort, viper.GetString("server.port"))
	}
}

func TestGetHandler(t *testing.T) {
	req, err := http.NewRequest("GET", server.URL+"/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"GET request successful"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPostHandler(t *testing.T) {
	data := map[string]string{"name": "test"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", server.URL+"/post", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"message":"POST request successful"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestServerShutdown(t *testing.T) {
	server := &http.Server{Addr: ":8080"}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("listen: %s\n", err)
		}
	}()
	time.Sleep(2 * time.Second) // Ensure the server starts

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("Server forced to shutdown: %s", err)
	}
}
