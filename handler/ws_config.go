package handler

import (
	"log"
	"net/http"
	"real-time-forum/model"
	"real-time-forum/utils"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	_ "github.com/mattn/go-sqlite3"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Adjust for production!
		},
	}
	Clients   = make(map[string]*websocket.Conn)
	Broadcast = make(chan model.WSMessage)
	Mu        sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket handshake from", r.RemoteAddr)

	// 1. Get token from query parameter
	sessionToken := r.URL.Query().Get("token")
	log.Println("Token received:", sessionToken)
	if sessionToken == "" {
		http.Error(w, "Missing authentication token", http.StatusBadRequest)
		return
	}

	// 2. Validate session before upgrade
	user, expiry, err := utils.SelectSession(sessionToken)
	if err != nil || time.Now().After(expiry) {
		http.Error(w, "Unauthorized or expired token", http.StatusUnauthorized)
		return
	}
	log.Println("Session valid, proceeding to upgrade...")

	// 3. Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	log.Println("WebSocket upgraded successfully")
	defer conn.Close()

	// 4. Connection setup
	handleConnect(conn, user)

	// 5. Message processing loop
	for {
		var msg model.WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		if err := HandleWebSocketMessage(conn, msg, user); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}

	// 6. Disconnection cleanup
	handleDisconnect(user)
}

func handleConnect(conn *websocket.Conn, user model.User) {
	Mu.Lock()
	defer Mu.Unlock()
	key := strconv.Itoa(user.ID)
	Clients[key] = conn
	log.Printf("User %v connected. Total clients: %v", user.ID, len(Clients))
	go TellAllToUpdateClients()
}

func handleDisconnect(user model.User) {
	Mu.Lock()
	defer Mu.Unlock()
	key := strconv.Itoa(user.ID)
	delete(Clients, key)
	UpdateOnlineTime(user.ID)
	log.Printf("User %s disconnected. Total clients: %d", user.UUID, len(Clients))
	go TellAllToUpdateClients()
}

// Example TellAllToUpdateClients implementation
func TellAllToUpdateClients() {
	Mu.Lock()
	defer Mu.Unlock()
	if len(Clients) == 0 {
		return
	}
	log.Println("Broadcasting updateClients to all clients")
	msg := model.WSMessage{MsgType: "updateClients"}
	for _, client := range Clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Broadcast error: %v", err)
		}
	}
}
func HandleBroadcasts() {
	for msg := range Broadcast {
		switch msg.MsgType {
		case "sendMessage":
			receiverKey := strconv.Itoa(msg.ReceiverUserID)
			senderKey := strconv.Itoa(msg.UserID)
			// Send to message receiver
			if receiverConn, ok := Clients[receiverKey]; ok {
				receiverConn.WriteJSON(msg)
			}
			// Send confirmation to sender
			if senderConn, ok := Clients[senderKey]; ok {
				senderConn.WriteJSON(msg)
			}
		case "updateClients":
			Mu.Lock()
			for _, client := range Clients {
				client.WriteJSON(msg)
			}
			Mu.Unlock()
		}
	}
}
