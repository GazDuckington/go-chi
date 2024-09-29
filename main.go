package main

import (
	"bobot/database"
	"bobot/middleware"
	repository "bobot/repositories"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

type user_payload map[string]interface{}

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	sample_user := user_payload{"user_id": 124}
	_, tokenString, _ := tokenAuth.Encode(sample_user)
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
	database.ConnectDatabase()
}

func main() {
	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, router())
}

func router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.LogCalls)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})

		// Create Entry
		r.Post("/entries", repository.CreateEntry)

		// Update Entry
		r.Put("/entries/{id}", repository.UpdateEntry)

		// Delete Entry
		r.Delete("/entries/{id}", repository.DeleteEntry)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous"))
		})

		r.Get("/entries", repository.GetAllEntry)
		r.Get("/entries/num/{num}", repository.FindEntryByNumber)
		r.Get("/entries/id/{id}", repository.FindEntryByID)
		r.Get("/entries/search", repository.GetEntriesByPattern)

		// auths
		r.Post("/register", repository.CreateUser)
		r.Post("/login", repository.LoginUser)
	})

	return r
}
