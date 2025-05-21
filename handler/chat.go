package handler

import (
	"database/sql"
	"fmt"
	"log"

	"real-time-forum/model"
	"real-time-forum/utils"
	"sort"
	"strings"
)

var DB *sql.DB

func InsertChat(user_id_1, user_id_2 int) (int, error) {
	log.Printf("Creating new chat between %d and %d", user_id_1, user_id_2)
	// Order user IDs to avoid duplicates
	if user_id_1 > user_id_2 {
		user_id_1, user_id_2 = user_id_2, user_id_1
	}

	uuid, err := utils.GenerateUuid()
	if err != nil {
		return 0, err
	}

	res, insertErr := DB.Exec(
		`INSERT INTO chats (uuid, user_id_1, user_id_2) VALUES (?, ?, ?)`,
		uuid, user_id_1, user_id_2,
	)
	if insertErr != nil {
		if sqliteErr, ok := insertErr.(interface{ Error() string }); ok &&
			strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed") {
			return FindChatIDbyUserIDS(user_id_1, user_id_2)
		}
		return 0, insertErr
	}
	chatID, err := res.LastInsertId()
	return int(chatID), err
}

func FindChatIDbyUserIDS(user_id_1, user_id_2 int) (int, error) {
	if user_id_1 > user_id_2 {
		user_id_1, user_id_2 = user_id_2, user_id_1
	}
	var chatID int
	err := DB.QueryRow(`
		SELECT id FROM chats 
		WHERE (user_id_1 = ? AND user_id_2 = ?) 
		OR (user_id_1 = ? AND user_id_2 = ?)`,
		user_id_1, user_id_2, user_id_2, user_id_1).Scan(&chatID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return chatID, err
}

func InsertMessage(content string, user_id_from int, chatID int) error {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("database error in InsertMessage", err)
		return err
	}
	// Update the chat's last activity
	_, updateErr := UpdateChat(chatID, user_id_from, tx)
	if updateErr != nil {
		fmt.Println("update error in InsertMessage", updateErr)
		tx.Rollback()
		return updateErr
	}

	insertQuery := `INSERT INTO messages (chat_id, user_id_from, content, created_at) 
                VALUES (?, ?, ?, CURRENT_TIMESTAMP)`;
	_, insertErr := tx.Exec(insertQuery, chatID, user_id_from, content)
	if insertErr != nil {
		fmt.Println("Insert error in InsertMessage", insertErr)
		tx.Rollback()
		return insertErr
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Error committing query at InsertMessage", err)
		return err
	}
	return nil
}

func UpdateChat(chatID int, userID int, tx *sql.Tx) (int, error) {
	// Perform the update
	query := `
		UPDATE chats
		SET updated_at = CURRENT_TIMESTAMP,
			updated_by = ?
		WHERE id = ?;
	`
	result, err := tx.Exec(query, userID, chatID)
	if err != nil {
		fmt.Println("Error, arguments:", chatID, userID)
		return 0, err
	}

	// Ensure at least one row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("rowsAffected error, arguments:", chatID, userID)
		return 0, err
	}
	if rowsAffected == 0 {
		fmt.Println("No rows afffected, arguments:", chatID, userID)
		return 0, sql.ErrNoRows // No rows were updated, meaning the UUID wasn't found
	}

	return chatID, nil
}

// ReadAllUsers retrieves all usernames: those the user has chatted with and those they haven't
func ReadAllUsers(userID int) ([]model.ChatUser, []model.ChatUser, error) {

	// Query the records
	rows, selectError := DB.Query(`
	SELECT u.nickname,
		u.id,
		u.age,
		u.gender,
		u.first_name,
		u.last_name,
		c.id AS chat_id,
		COALESCE(c.updated_at, c.created_at) AS last_activity,
		(SELECT EXISTS(
            SELECT 1 FROM messages m
            WHERE m.chat_id = c.id
            AND m.user_id_from = u.id
            AND m.user_id_from != ?
            AND m.id > IFNULL(
                (SELECT MAX(m2.id) FROM messages m2
                WHERE m2.chat_id = c.id
                AND m2.user_id_from = ?), 0)
        )) AS has_unread
	FROM users u
	LEFT JOIN chats c 
		ON (u.id = c.user_id_1 OR u.id = c.user_id_2)
		AND (c.user_id_1 = ? AND c.user_id_2 = u.id) 
		OR (c.user_id_2 = ? AND c.user_id_1 = u.id)
	WHERE u.id != ?
	ORDER BY last_activity DESC;
    `, userID, userID, userID, userID, userID)

	if selectError != nil {
		fmt.Println("Select error in ReadAllUsers:", selectError)
		return nil, nil, selectError
	}
	defer rows.Close()

	var chattedUsers []model.ChatUser
	var notChattedUsers []model.ChatUser

	// Iterate over rows and collect usernames
	for rows.Next() {
		var chatID sql.NullInt64
		var chatUser model.ChatUser
		var hasUnread bool

		err := rows.Scan(
			&chatUser.Username,
			&chatUser.UserID,
			&chatUser.Age, 
			&chatUser.Gender,
			&chatUser.FirstName,
			&chatUser.LastName,
			&chatID,
			&chatUser.LastActivity,
			&hasUnread,
		)
		if err != nil {
			return nil, nil, err
		}

		if chatID.Valid {
			chattedUsers = append(chattedUsers, chatUser)
		} else {
			notChattedUsers = append(notChattedUsers, chatUser)
		}
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	  // For each user, set the HasUnread flag based on database query
	for i := range chattedUsers {
        // Check if there are unread messages from this user
        hasUnread, err := CheckUnreadMessages(userID, chattedUsers[i].UserID)
        if err == nil {
            chattedUsers[i].HasUnread = hasUnread
        }
    }
	// Sort non-chatted users alphabetically
	sort.Slice(notChattedUsers, func(i, j int) bool {
		return strings.ToLower(notChattedUsers[i].Username) < strings.ToLower(notChattedUsers[j].Username)
	})

	return chattedUsers, notChattedUsers, nil

}

func CheckUnreadMessages(userID int, senderID int) (bool, error) {
    chatID, err := FindChatIDbyUserIDS(userID, senderID)
    if err != nil {
        return false, err
    }
    
    if chatID == 0 {
        return false, nil
    }
    
    var hasUnread bool
    query := `
    SELECT EXISTS(
        SELECT 1 FROM messages m
        WHERE m.chat_id = ?
        AND m.user_id_from = ?
        AND m.created_at > IFNULL(
            (SELECT read_at FROM message_read_receipts 
				WHERE chat_id = ? AND user_id = ?), 
            '1970-01-01')
    ) AS has_unread
    `
    
    err = DB.QueryRow(query, chatID, senderID, chatID, userID).Scan(&hasUnread)
    if err != nil {
        log.Printf("Error checking unread messages: %v", err)
        return false, err
    }
    
    return hasUnread, nil
}

func ClearUnreadMessages(userID int, senderID int) error {
    chatID, err := FindChatIDbyUserIDS(userID, senderID)
    if err != nil {
        log.Printf("Error finding chat: %v", err)
        return err
    }
    
    if chatID == 0 {
        log.Printf("No chat found between users %d and %d", userID, senderID)
        return nil
    }
    
    log.Printf("Clearing unread messages in chat %d for user %d from sender %d", 
				chatID, userID, senderID)
    
    // SQLite syntax for UPSERT
    _, err = DB.Exec(`
        INSERT INTO message_read_receipts (chat_id, user_id, read_at)
        VALUES (?, ?, CURRENT_TIMESTAMP)
        ON CONFLICT(chat_id, user_id) 
        DO UPDATE SET read_at = CURRENT_TIMESTAMP
    `, chatID, userID)
    
    if err != nil {
        log.Printf("Error updating read receipt: %v", err)
        return err
    }
    
    return nil
}

func ReadAllMessages(chatID int, numberOfMessages int, userID int, offset int) ([]model.PrivateMessage, error) {
	lastMessages := make([]model.PrivateMessage, 0)

	// Query messages along with the sender's username
	rows, selectError := DB.Query(`
        SELECT 
            m.id AS message_id, 
            m.chat_id,
            m.user_id_from, 
            u.nickname AS sender_nickname, 
            m.content, 
            m.updated_at, 
            m.created_at
        FROM messages m
        INNER JOIN users u 
            ON m.user_id_from = u.id
        WHERE m.chat_id = ?
        ORDER BY m.id DESC
        LIMIT ? OFFSET ?;
    `, chatID, numberOfMessages, offset)

	if selectError != nil {
		fmt.Println("Select error at ReadAllMessages:", selectError)
		return nil, selectError
	}
	defer rows.Close()

	// Iterate over rows and collect messages
	for rows.Next() {
		var message model.PrivateMessage

		err := rows.Scan(
			&message.Message.ID,
			&message.Message.ChatID,
			&message.Message.UserIDFrom,
			&message.Message.SenderUsername,
			&message.Message.Content,
			&message.Message.UpdatedAt,
			&message.Message.CreatedAt,
		)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			return nil, err
		}
		message.IsCreatedBy = (message.Message.UserIDFrom == userID)
		lastMessages = append(lastMessages, message)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		fmt.Printf("Rows iteration error: %v\n", err)
		return nil, err
	}

	if lastMessages == nil {
    lastMessages = []model.PrivateMessage{}
	}
	return lastMessages, nil
}
