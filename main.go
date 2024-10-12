package main

import (
	"fmt"
	"my-blog-backend/routes"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to My Blog API")
}

func main() {

	routes.InitDB()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/signup", routes.SignUpHandler)
	http.HandleFunc("/users", routes.GetUsersHandler)
	http.HandleFunc("/login", routes.LoginHandler)

	http.HandleFunc("/posts", routes.CreatePostHandler)
	http.HandleFunc("/posts/all", routes.GetPostsHandler)
	http.HandleFunc("/posts/delete", routes.DeletePostHandler)
	http.HandleFunc("/posts/update", routes.UpdatePostHandler)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
