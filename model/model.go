package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int    `json:"id" db:"id"`
	UUID           string `json:"uuid" db:"uuid"`
	Age            string `json:"age" db:"age"`
	Gender         string `json:"gender" db:"gender"`
	FirstName      string `json:"firstName" db:"first_name"`
	LastName       string `json:"lastName" db:"last_name"`
	Nickname       string `json:"nickname" db:"nickname"`
	Email          string `json:"email" db:"email"`
	Password       string `json:"password" db:"password"`
	SessionToken   sql.NullString
	SessionExpiry  sql.NullTime
	LastTimeOnline time.Time  `json:"lastTimeOnline"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	UpdatedBy      *int       `json:"updated_by"`
}

type Post struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Categories string    `json:"categories"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Nickname   string    `json:"nickname"` // Join with users table
}

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

type Comment struct {
	UserName  string `json:"user_name"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Session struct {
	ID           int       `json:"id"`
	SessionToken string    `json:"session_token"`
	UserId       int       `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Gender    string `json:"gender"`
	Age       string `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	LoginType  string `json:"loginType"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
	Token   string `json:"token,omitempty"`
}

type WSMessage struct {
	Post                    string           `json:"post"`
	Comment                 string           `json:"comment"`
	MsgType                 string           `json:"msgType"`
	Updated                 bool             `json:"updated"`
	UserID                  int              `json:"user_id"`
	IsLikAction             bool             `json:"isLikeAction"`
	NumberOfReplies         int              `json:"numberOfReplies"`
	IsReplied               bool             `json:"isReplied"`
	ChattedUsers            []ChatUser       `json:"chattedUsers"`
	UnchattedUsers          []ChatUser       `json:"unchattedUsers"`
	PrivateMessage          PrivateMessage   `json:"privateMessage"`
	ReceiverUserID          int              `json:"receiver_user_id"`
	ReceiverUserName        string           `json:"receiver_nickname"`
	TypingNickname 			string 			 `json:"typing_nickname"`
	Messages                []PrivateMessage `json:"messages"`
	SendNotification        bool             `json:"notification"`
	GotAllMessagesRequested bool             `json:"allMessagesGot"`
}

type Chat struct {
	ID        int        `json:"id"`
	UUID      string     `json:"uuid"`
	User_id_1 int        `json:"user_id_1"`
	User_id_2 int        `json:"user_id_2"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *int       `json:"updated_by"`
}

type Message struct {
	ID int `json:"id"`
	// ChatUUID       string     `json:"chat_uuid"`
	ChatID         int        `json:"chat_id"`
	UserIDFrom     int        `json:"user_id_from"`
	SenderUsername string     `json:"sender_nickname"`
	Content        string     `json:"content"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
type ChatUser struct {
	// User         string         `json:"user"`
	Username     string         `json:"nickname"`
	UserID       int            `json:"user_id"`
	Age          string         `json:"age"`
	Gender       string         `json:"gender"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	LastActivity sql.NullString `json:"lastActivity"` // Changed to NullString
	ChatID       sql.NullInt64  `json:"chat_id"`
	IsOnline     bool           `json:"isOnline"`
}

type PrivateMessage struct {
	Message     Message `json:"message"` // Explicit lowercase
	IsCreatedBy bool    `json:"isCreatedBy"`
}
