package routes

import (
	"encoding/json"
	"errors"
	"log"
	"my-blog-backend/models"
	"net/http"
	"strconv"
)

// validate post
func validatePost(title, category, content string) error {
	if title == "" || category == "" || content == "" {
		return errors.New("all fields are required")
	}

	return nil
}

// Create Post Handler
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newPost models.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = validatePost(newPost.Title, newPost.Category, newPost.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO posts (title, category, content) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, newPost.Title, newPost.Category, newPost.Content)
	if err != nil {
		log.Println("Failed to log post:", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post created successfully"))
}

// Get Post Handler
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`SELECT id, title, category, content, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		log.Println("Failed to retrieve posts: ", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Category, &post.Content, &post.CreatedAt)
		if err != nil {
			log.Println("Failed to scan post: ", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error: ", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}

// Delete Post Handler
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM posts WHERE id =$1`
	_, err = db.Exec(query, id)
	if err != nil {
		log.Println("Failed to delete post: ", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}

// UpdatePostHandler updates an existing post
func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the post ID from the query parameters
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Decode the request body into a Post struct
	var updatedPost models.Post
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE posts SET title = $1, category = $2, content = $3 WHERE id = $4`
	_, err = db.Exec(query, updatedPost.Title, updatedPost.Category, updatedPost.Content, id)
	if err != nil {
		log.Println("Failed to update post: ", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post updated successfully"))
}
