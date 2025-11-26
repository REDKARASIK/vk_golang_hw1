package main

import (
	"log"
	"net/http"
	"redditclone/internal/handler"
	"redditclone/internal/middleware"
	"redditclone/internal/post"
	"redditclone/internal/user"
	"strings"
	"time"
)

func stripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize repositories (in-memory)
	userRepo := user.NewMemoryRepo()
	postRepo := post.NewMemoryRepo()

	// Initialize handlers
	userHandler := handler.NewUserHandler(userRepo)
	postHandler := handler.NewPostHandler(postRepo)

	// Main router
	mux := http.NewServeMux()

	// --- Public routes ---
	// User routes
	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("POST /api/login", userHandler.Login)
	mux.HandleFunc("GET /api/user/{USER_LOGIN}", postHandler.ListByUser)

	// Post routes
	mux.HandleFunc("GET /api/posts", postHandler.List)
	mux.HandleFunc("GET /api/posts/{CATEGORY_NAME}", postHandler.ListByCategory)
	mux.HandleFunc("GET /api/post/{POST_ID}", postHandler.GetByID)

	// --- Authenticated routes ---
	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /api/posts", postHandler.Add)
	authMux.HandleFunc("POST /api/post/{POST_ID}", postHandler.AddComment)
	authMux.HandleFunc("DELETE /api/post/{POST_ID}/{COMMENT_ID}", postHandler.DeleteComment)
	authMux.HandleFunc("GET /api/post/{POST_ID}/upvote", postHandler.Upvote)
	authMux.HandleFunc("GET /api/post/{POST_ID}/downvote", postHandler.Downvote)
	authMux.HandleFunc("GET /api/post/{POST_ID}/unvote", postHandler.Unvote)
	authMux.HandleFunc("DELETE /api/post/{POST_ID}", postHandler.Delete)

	// Apply Auth middleware to the authenticated router
	mux.Handle("/api/", middleware.Auth(authMux))

	// --- Static file serving ---
	// Serve static files from the html directory
	staticHTMLHandler := http.FileServer(http.Dir("./static/html"))
	mux.Handle("/", staticHTMLHandler)

	// Serve other static assets like css, js
	staticHandler := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	log.Println("Starting server on :8080")
	server := &http.Server{
		Addr:         ":8080",
		Handler:      stripTrailingSlash(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
