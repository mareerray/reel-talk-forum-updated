package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"real-time-forum/model"

	"github.com/gorilla/websocket"
)

type MessageHandler func(conn *websocket.Conn, msg model.WSMessage, user model.User) error

var messageHandlers = map[string]MessageHandler{
	"getUsers":       handleGetUsers,
	"sendMessage":    handleSendMessage,
	"getOrCreateChat": handleGetOrCreateChat,
    "getMessages":     handleGetMessages,
	"typing":         handleTyping,
	"stopped_typing": handleStoppedTyping,
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func HandleWebSocketMessage(conn *websocket.Conn, msg model.WSMessage, user model.User) error {
	if handler, exists := messageHandlers[msg.MsgType]; exists {
		return handler(conn, msg, user)
	}
	log.Printf("No handler for message type: %s", msg.MsgType)
	return nil
}

func handleGetUsers(conn *websocket.Conn, msg model.WSMessage, user model.User) error {

	chatted, unchatted, err := ReadAllUsers(user.ID)
	if err != nil {
		log.Println("Error fetching users:", err)
		return err
	}

	// Set online status
	for i := range chatted {
		key := strconv.Itoa(chatted[i].UserID)
		chatted[i].IsOnline = Clients[key] != nil
	}
	for i := range unchatted {
		key := strconv.Itoa(unchatted[i].UserID)
		unchatted[i].IsOnline = Clients[key] != nil
	}

	return conn.WriteJSON(model.WSMessage{
		MsgType:        "listOfChat",
		UserID:         user.ID,
		ChattedUsers:   chatted,
		UnchattedUsers: unchatted,
	})
}

func handleSendMessage(conn *websocket.Conn, msg model.WSMessage, user model.User) error {

	if msg.PrivateMessage.Message.ChatID == 0 ||
		msg.ReceiverUserID == 0 ||
		msg.PrivateMessage.Message.Content == "" {
		log.Println("Invalid message format")
		return nil
	}

	// Insert to database
	err := InsertMessage(
		msg.PrivateMessage.Message.Content,
		user.ID,
		msg.PrivateMessage.Message.ChatID,
	)
	if err != nil {
		log.Printf("Message insert failed: %v", err)
		return err
	}

	// Fetch the newly created message with timestamp
	messages, err := ReadAllMessages(msg.PrivateMessage.Message.ChatID, 1, user.ID, 0)
	if err != nil {
		log.Printf("Failed to fetch new message: %v", err)
		return err
	}

	// Update the message in the response with the database values
	if len(messages) > 0 {
		msg.PrivateMessage = messages[0]
	}

	// Set sender's user ID
    msg.UserID = user.ID


    // Create receiver message (isCreatedBy = false)
    receiverMsg := msg
    receiverMsg.PrivateMessage.IsCreatedBy = false
	receiverMsg.SendNotification = true

	// Send to receiver
	receiverKey := strconv.Itoa(msg.ReceiverUserID)
	if receiverConn, ok := Clients[receiverKey]; ok {
		receiverConn.WriteJSON(receiverMsg)
	}
	// Create sender confirmation (isCreatedBy = true)
    senderMsg := msg
    senderMsg.PrivateMessage.IsCreatedBy = true

	// Confirm to sender
	return conn.WriteJSON(senderMsg)
}

func handleGetOrCreateChat(conn *websocket.Conn, msg model.WSMessage, user model.User) error {
    // Extract receiver user ID from message
    receiverID := msg.ReceiverUserID
    if receiverID == 0 {
        log.Println("Invalid receiver user ID")
        return nil
    }

    // Find existing chat or create new one
    chatID, err := FindChatIDbyUserIDS(user.ID, receiverID)
    if err != nil {
        log.Printf("Error finding chat: %v", err)
        return err
    }

    if chatID == 0 {
        chatID, err = InsertChat(user.ID, receiverID)
        if err != nil {
            log.Printf("Error creating chat: %v", err)
            return err
        }
    }

	if msg.ClearUnread {
        err := ClearUnreadMessages(user.ID, receiverID)
        if err != nil {
            log.Printf("Error clearing unread messages: %v", err)
        }
    }

	response := model.WSMessage{
		MsgType: "chatCreated",
		UserID:  user.ID,
		PrivateMessage: model.PrivateMessage{
			Message: model.Message{
				ChatID: chatID,
			},
			IsCreatedBy: false,
		},
	}
	return conn.WriteJSON(response)
}

func handleGetMessages(conn *websocket.Conn, msg model.WSMessage, user model.User) error {
    chatID := msg.PrivateMessage.Message.ChatID
    if chatID == 0 {
        log.Println("Invalid chat ID")
        return nil
    }

	numberOfMessages := 10
	page := msg.Page
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * numberOfMessages
	messages, err := ReadAllMessages(chatID, numberOfMessages, user.ID, offset)

    if err != nil {
        log.Printf("Error reading messages: %v", err)
        return err
    }
	log.Printf("Fetching messages for chat %d page %d with limit %d and offset %d", 
			chatID, page, numberOfMessages, offset)
			
	hasMore := len(messages) == numberOfMessages

    return conn.WriteJSON(model.WSMessage{
        MsgType: "messages",
        UserID:  user.ID,
        Messages: messages,
		Page: page,
		HasMore: hasMore,
    })
}

func handleTyping(conn *websocket.Conn, msg model.WSMessage, user model.User) error {
    receiverKey := strconv.Itoa(msg.ReceiverUserID)
    Mu.Lock()
    defer Mu.Unlock()
    if receiverConn, ok := Clients[receiverKey]; ok {
        msg.UserID = user.ID
        msg.TypingNickname = user.Nickname
        return receiverConn.WriteJSON(msg)
    }
    return nil
}


func handleStoppedTyping(conn *websocket.Conn, msg model.WSMessage, user model.User) error {
	return handleTyping(conn, msg, user) // Same logic as typing
}
