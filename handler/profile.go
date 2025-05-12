package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	// "real-time-forum/auth"
	"encoding/json"
	"net/http"
	"real-time-forum/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// In your API handlers file (e.g., userHandlers.go)

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	if r.URL.Path != "/profile" {
		loggedIn, userID := CheckUserLoggedIn(r)
		if !loggedIn {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Query database for user information
		var user struct {
			Nickname  string `json:"nickname"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Gender    string `json:"gender"`
			Age       int    `json:"age"`
			Email     string `json:"email"`
		}

		query := `SELECT nickname, first_name, last_name, gender, age, email 
            FROM users WHERE id = ?`

		err := DB.QueryRow(query, userID).Scan(
			&user.Nickname,
			&user.FirstName,
			&user.LastName,
			&user.Gender,
			&user.Age,
			&user.Email,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Return user data as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// User struct represents the user data model
type User struct {
	ID             int        `json:"id"`
	UUID           string     `json:"uuid"`
	Type           string     `json:"type"`
	Age            string     `json:"age"`
	Gender         string     `json:"gender"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Username       string     `json:"nickname"`
	Email          string     `json:"email"`
	Password       string     `json:"password"`
	Status         string     `json:"status"`
	LastTimeOnline time.Time  `json:"lastTimeOnline"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	UpdatedBy      *int       `json:"updated_by"`
}

func InsertUser(user *User) (int, error) {

	// Generate UUID for the user if not already set
	if user.UUID == "" {
		uuid, err := utils.GenerateUuid()
		if err != nil {
			return -1, err
		}
		user.UUID = uuid
	}

	var existingEmail string
	var existingUsername string
	emailCheckQuery := `SELECT email, nickname FROM users WHERE email = ? OR nickname = ? LIMIT 1;`
	err := DB.QueryRow(emailCheckQuery, user.Email, user.Username).Scan(&existingEmail, &existingUsername)
	if err == nil {
		if existingEmail == user.Email {
			return -1, errors.New("duplicateEmail")
		}
		if existingUsername == user.Username {
			return -1, errors.New("duplicateUsername")
		}
	}

	//insertQuery := `INSERT INTO users (uuid, name, username, email, password) VALUES (?, ?, ?, ?, ?);`
	//result, insertErr := db.Exec(insertQuery, user.UUID, user.Username, user.Username, user.Email, user.Password)

	insertQuery := `INSERT INTO users (uuid, username, email, password, age, gender, firstname, lastname) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	result, insertErr := DB.Exec(insertQuery, user.UUID, user.Username, user.Email, user.Password, user.Age, user.Gender, user.FirstName, user.LastName)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation (duplicate entry)
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // 19 = UNIQUE constraint failed (SQLite error code)
				return -1, errors.New("user with this email or username already exists")
			}
		}
		return -1, insertErr // Other DB errors
	}

	// Retrieve the last inserted ID
	userId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return int(userId), nil
}

func AuthenticateUser(input, password string) (int, error) {
	// Open SQLite database/ Cl

	// Query to retrieve the hashed password stored in the database for the given username
	var userID int
	var storedHashedPassword string
	err := DB.QueryRow("SELECT id, password FROM users WHERE username = ? OR email = ?", input, input).Scan(&userID, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username not found
			return -1, errors.New("username or email not found")
		}
		return -1, err
	}

	// Compare the entered password with the stored hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		// Password is incorrect
		return -1, errors.New("password is incorrect")
	}

	// Successful login if no errors occurred
	return userID, nil
}

func FindUsernameByID(ID int) (string, error) {

	selectQuery := `
		SELECT
			username
		FROM users
			WHERE id = ?;
	`
	idRow, selectError := DB.Query(selectQuery, ID)
	if selectError != nil {
		return "", selectError
	}

	var name string
	for idRow.Next() {
		if err := idRow.Scan(&name); err != nil {
			fmt.Printf("Failed to scan row: %v\n", err)
		}
	}

	return name, nil
}
func FindUsername(userUUID string) (string, error) {
    var nickname string
    err := DB.QueryRow("SELECT nickname FROM users WHERE uuid = ?", userUUID).Scan(&nickname)
    if err != nil {
        log.Printf("[ERROR] FindUsername failed for UUID %s: %v", userUUID, err)
    }
    return nickname, err
}

// func FindUsername(userUUID string) (string, error) {

// 	var nickname string

// 	selectQuery := `
// 		SELECT
// 			nickname
// 		FROM users
// 			WHERE uuid = ?;
// 	`
// 	idRow, selectError := DB.Query(selectQuery, userUUID)
// 	if selectError != nil {
// 		return "", selectError
// 	}


// 	for idRow.Next() {
// 		if err := idRow.Scan(&nickname); err != nil {
// 			fmt.Printf("Failed to scan row: %v\n", err)
// 		}
// 	}

// 	return nickname, nil
// }

func UpdateOnlineTime(UserID int) error {

	updateQuery := `UPDATE users
	SET 
		last_activity = CURRENT_TIMESTAMP
	WHERE id = ?;`
	_, updateErr := DB.Exec(updateQuery, UserID)
	if updateErr != nil {
		fmt.Println(updateErr)
		return updateErr
	}
	return nil
}

func FindUserByUUID(UUID string) (int, error) {
	var id int
	err := DB.QueryRow("SELECT id FROM users WHERE uuid = ?", UUID).Scan(&id)
	log.Printf("[DEBUG] UUID %s → ID %d (error: %v)", UUID, id, err)
	return id, err
  }
  
// func FindUserByUUID(UUID string) (int, error) {

// 	selectQuery := `
// 		SELECT
// 			id
// 		FROM users
// 			WHERE uuid = ?;
// 	`
// 	idRow, selectError := DB.Query(selectQuery, UUID)
// 	if selectError != nil {
// 		return -1, selectError
// 	}

// 	var id int
// 	for idRow.Next() {
// 		if err := idRow.Scan(&id); err != nil {
// 			fmt.Printf("Failed to scan row: %v\n", err)
// 		}
// 	}

// 	return id, nil
// }
