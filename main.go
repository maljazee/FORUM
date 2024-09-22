package main

import (
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
)

func main() {
	// Initialize database
	err := database.InitDB("forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	// Set up routes
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/create-post", handlers.CreatePost)
	http.HandleFunc("/post", handlers.GetPost)
	http.HandleFunc("/create-comment", handlers.CreateComment)
	http.HandleFunc("/like", handlers.LikePost)
	http.HandleFunc("/dislike", handlers.DislikePost)

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
