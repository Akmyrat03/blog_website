package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func validateUserInput(username, password, confirmPassword string) error {
	if username == "" || password == "" || confirmPassword == "" {
		return errors.New("all fields are required")
	}

	//Check if username already exists
	for _, u := range users {
		if u.Username == username {
			return errors.New("username already exists")
		}
	}

	//Check if password and confirm password match
	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	return nil
}

// SignUp handler
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//validate user input
	err = validateUserInput(input.Username, input.Password, input.ConfirmPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	//Insert the user into the database
	query := `INSERT INTO public.users (username, password) VALUES ($1, $2)`
	_, err = db.Exec(query, input.Username, string(hashedPassword))
	if err != nil {
		log.Println("Failed to insert user: ", err)
		http.Error(w, "Failed to sign up", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("sign up successful"))
}

// Login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var loginUser User
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}

	var storedUser User
	query := `SELECT username, password FROM public.users WHERE username = $1`
	err = db.QueryRow(query, loginUser.Username).Scan(&storedUser.Username, &storedUser.Password)
	if err != nil {
		log.Println("Failed to fetch user:", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginUser.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// get users handler
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`SELECT id, username, password FROM users`)
	if err != nil {
		log.Println("Failed to fetch users: ", err)
		http.Error(w, "Failed to fetch users: ", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			log.Println("Failed to scan user: ", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error: ", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}
