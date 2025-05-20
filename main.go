package main

import (
	"log"
	"net/http"
	"real-time-forum/handler"
	"real-time-forum/utils"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Initialize database
	var err error

	utils.DB = utils.OpenDBConnection()

	defer utils.DB.Close()

	// Optionally, test the connection
	if err = utils.DB.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Database connection established")

	// Set the database connection for handlers
	handler.DB = utils.DB

	// API routes
	http.HandleFunc("/api/categories", handler.GetCategoriesHandler)
	http.HandleFunc("/api/posts/list", handler.GetPostsHandler)
	http.HandleFunc("/api/user/profile", handler.GetUserProfile)
	// Add these to your API routes
	http.HandleFunc("/api/comments", handler.PostCommentHandler)  // POST new comment
	http.HandleFunc("/api/comments/", handler.GetCommentsHandler) // GET comments for post
	http.HandleFunc("/api/posts", handler.CreatePostHandler)      // POST to create a new post

	// New API routes for user authentication
	http.HandleFunc("/api/register", handler.RegisterHandler)
	http.HandleFunc("/api/login", handler.LoginHandler)
	http.HandleFunc("/api/logout", handler.LogoutHandler)
	http.HandleFunc("/api/validate-session", handler.ValidateSessionHandler)

	http.HandleFunc("/ws", handler.HandleConnections)

	// Parse the HTML template
	indexTemplate, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Serve static files from the "assets" directory
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Handle the root route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title string
		}{
			Title: "My Single Page App",
		}

		if err := indexTemplate.Execute(w, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
		}
	})

	log.Println("Server starting at http://localhost:8999")
	if err := http.ListenAndServe(":8999", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
