package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type InMemStore struct {
	mu   sync.RWMutex
	urls map[string]string
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		urls: make(map[string]string),
	}
}

func (s *InMemStore) Set(code, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[code] = url
}

func (s *InMemStore) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists := s.urls[code]
	return url, exists
}

type Application struct {
	store   *InMemStore
	baseURL string
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Message  string `json:"message"`
	ShortURL string `json:"short_url"`
}

func generateShortKey() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (app *Application) handleShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "Payload JSON invalide ou champ 'url' manquant", http.StatusBadRequest)
		return
	}

	code, err := generateShortKey()
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	app.store.Set(code, req.URL)

	resp := ShortenResponse{
		Message:  "URL enregistrée",
		ShortURL: fmt.Sprintf("%s/%s", app.baseURL, code),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (app *Application) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	targetURL, exists := app.store.Get(code)
	if !exists {
		http.NotFound(w, r)
		return
	}

	// 307 préserve la méthode et n'est pas mis en cache agressivement par les navigateurs clients
	http.Redirect(w, r, targetURL, http.StatusTemporaryRedirect)
}

func main() {
	app := Application{
		store:   NewInMemStore(),
		baseURL: "http://localhost:8080",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /url", app.handleShortenURL)
	mux.HandleFunc("GET /{code}", app.handleRedirect)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server ready on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error : %v", err)
		}
	}()

	<-shutdownChan
	log.Println("Server stop")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Can't stop server")
	}

	log.Println("Server stopped")

}
