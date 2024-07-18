package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"sync"

	"net/http"
	"time"
)

var (
	secretKey    = "secret"
	currentGuess = generateHiddenNumber()
)

var mu sync.Mutex

var userToken string

const (
	predefinedToken    = "admin_token"
	predefinedUsername = "mickie"
	predefinedPassword = "kuay" // You had hardcoded the password in the login check
	tokenHeaderKey     = "Authorization"
	bearerPrefix       = "Bearer"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", loginHandler)
	mux.Handle("POST /guess", tokenMiddleware(http.HandlerFunc(guessHandler)))

	// mux.Handle("POST /xxx", corsMiddleware(tokenMiddleware(http.HandlerFunc(guessHandler))))

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        corsMiddleware(mux),
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
	nBig, err := rand.Int(rand.Reader, big.NewInt(27))
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64())
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	var credentials LoginParams

	//decode json body of the request intio the credential  struct
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		//if decode fail return bad request
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v", credentials)
	//check if the username and password is correct
	if credentials.Username == predefinedUsername && credentials.Password == predefinedPassword {
		//generate token and send it to the client
		token := generateToken()
		userToken = token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		currentGuess = generateHiddenNumber()
		fmt.Println("correct number:", currentGuess)
	} else {
		//if the username and password is incorrect return unauthorized
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	var guessData struct {
		Guess string `json:"guess"`
	}
	if err := json.NewDecoder(r.Body).Decode(&guessData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// fmt.Println("Current guess:", currentGuess, "User guess:", guessData.Guess)

	if guessData.Guess == fmt.Sprint(currentGuess) {

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Correct guess!"})
		currentGuess = generateHiddenNumber()
		fmt.Println("Next correct number:", currentGuess)
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

func generateToken() string {
	var token [32]byte
	_, err := io.ReadFull(rand.Reader, token[:])
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(token[:])
}

func validateToken(token string) bool {
	return token == bearerPrefix+" "+userToken
}

func corsMiddleware(next http.Handler) http.Handler {
	// fmt.Println("CORS middleware")
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
