package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var (
	secretKey    = "secret"
	currentGuess = generateHiddenNumber()
)

const (
	predefinedToken    = "admin_token"
	predefinedUsername = "admin"
	predefinedPassword = "password" // You had hardcoded the password in the login check
	tokenHeaderKey     = "Authorization"
	bearerPrefix       = "Bearer "
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		tokenMiddleware(http.HandlerFunc(guessHandler)).ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Minute,
		WriteTimeout:   time.Minute,
		IdleTimeout:    75 * time.Second,
	}

	fmt.Println("Server is running on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Server failed to start:", err)
	}
	srv.ListenAndServe()
}

func generateHiddenNumber() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(100) + 1 // random number from 1 to 100
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//decode json body of the request intio the credential  struct
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		//if decode fail return bad request
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//check if the username and password is correct
	if credentials.Username == predefinedUsername && credentials.Password == predefinedPassword {
		//generate token and send it to the client
		token := generateToken(credentials.Username)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	} else {
		//if the username and password is incorrect return unauthorized
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	var guessData struct {
		Guess int `json:"guess"`
	}
	if err := json.NewDecoder(r.Body).Decode(&guessData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if guessData.Guess == currentGuess {
		currentGuess = generateHiddenNumber()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Correct guess!"})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "Try again!"})
	}
}

func tokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if validateToken(token) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func generateToken(username string) string {
	// This is a placeholder. In a real application, use JWT or another secure method.
	return fmt.Sprintf("%s_token", username)
}

func validateToken(token string) bool {
	return token == bearerPrefix+predefinedToken
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
