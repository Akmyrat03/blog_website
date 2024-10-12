package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

var posts []Post
var nextID = 1

// Create a new post
func CreatePost(title, content, category string) Post {
	post := Post{
		ID:        nextID,
		Title:     title,
		Content:   content,
		Category:  category,
		CreatedAt: time.Now(),
	}

	posts = append(posts, post)
	nextID++
	return post
}

// Get all posts
func GetAllPosts() []Post {
	return posts
}

// Delete a post by ID
func DeletePost(id int) bool {
	for index, post := range posts {
		if post.ID == id {
			posts = append(posts[:index], posts[index+1:]...)
			return true
		}
	}
	return false
}

// Update a post by ID
func UpdatePost(id int, title string, content string, category string) bool {
	for index, post := range posts {
		if post.ID == id {
			posts[index].Title = title
			posts[index].Content = content
			posts[index].Category = category
			return true
		}
	}
	return false
}
