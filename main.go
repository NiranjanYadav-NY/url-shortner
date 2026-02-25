package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

// In-Memory URL Store

type URLStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		data: make(map[string]string),
	}
}

var urlStore = NewURLStore()


// Short ID Generator

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateShortID() string {
	shortID := make([]byte, 6)
	for i := range shortID {
		shortID[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortID)
}

// JSON Structures

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// Handlers

// POST /shorten
func ShortenURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Method not allowed"})
		return
	}

	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid URL"})
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "URL must start with http:// or https://"})
		return
	}

	shortID := generateShortID()

	urlStore.mu.Lock()
	urlStore.data[shortID] = req.URL
	urlStore.mu.Unlock()

	shortURL := fmt.Sprintf("http://%s/%s", r.Host, shortID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShortenResponse{ShortURL: shortURL})
}

// GET /{shortID}
func RedirectURL(w http.ResponseWriter, r *http.Request) {

	shortID := strings.TrimPrefix(r.URL.Path, "/")

	if len(shortID) != 6 {
		http.NotFound(w, r)
		return
	}

	urlStore.mu.RLock()
	originalURL, exists := urlStore.data[shortID]
	urlStore.mu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// Homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// Static files
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

// MAIN

func main() {

	// Serve static files correctly
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API
	http.HandleFunc("/shorten", ShortenURL)

	// Home + Redirect
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			HomeHandler(w, r)
			return
		}
		RedirectURL(w, r)
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}