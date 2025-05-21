document.getElementById('sendButton').addEventListener('click', sendMessage);
document.getElementById('messageInput').addEventListener('keydown', (event) => {
    if (event.key === 'Enter') {
        sendMessage();
    }
});

let currentChatID;
let currentChatUser;
let typingTimeout;
let isTypingIndicatorVisible = false;
let currentMessagePage = 1;
let hasMoreMessages = true;
const messagesPerPage = 10;
const displayedMessageIds = new Set();

messageInput.addEventListener('input', () => {
    if (!currentChatUser?.user_id) return;
    // Only send typing event if indicator isn't already visible
    if (!isTypingIndicatorVisible) {
        isTypingIndicatorVisible = true;
        const typingPayload = {
            msgType: "typing",
            chat_id: currentChatID,
            receiver_user_id: currentChatUser.user_id,
        };
        socket.send(JSON.stringify(typingPayload));
    }
    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
        isTypingIndicatorVisible = false;
        const stopPayload = {
            msgType: "stopped_typing",
            chat_id: currentChatID,
            receiver_user_id: currentChatUser.user_id,
        };
        socket.send(JSON.stringify(stopPayload));
    }, 2000);
    // 5 seconds is the widely recommended default for chat apps, and will feel smoother for
});


let socket;

function initializeApp() {
    const token = localStorage.getItem('sessionToken');
    console.log('WebSocket token:', token);
    if (!token) {
        return;
    }
    socket = new WebSocket(`ws://localhost:8999/ws?token=${encodeURIComponent(token)}`);

    socket.onclose = () => {
        console.log("WebSocket closed. Reconnecting...");
        setTimeout(initializeApp, 3000); // Reconnect after 3 seconds
    };
        
    // WebSocket Message Listener - Sent
    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);

            if (messageHandlers[data.msgType]) {
                messageHandlers[data.msgType](data);
            } else {
                console.warn("No handler for message type:", data.msgType);
            }
        } catch (error) {
            console.error("Message handling error:", error);
        }
    };
}

// ------------------------------------------------------------------------------ //

// Message Handlers - Received
const messageHandlers = {
    listOfChat: (data) => {
        renderUserLists(data.chattedUsers, data.unchattedUsers);
    },

    updateClients: (data) => {
        console.log("Status update - refreshing list...");
        requestUserListViaWebSocket();
    },

    sendMessage: (data) => {
        const chatId = data.privateMessage?.message?.chat_id;
        // Always refresh the user list when receiving any message
        // This ensures notifications appear even from new users
        requestUserListViaWebSocket();
        
        if (chatId === currentChatID) {
            addMessageToChat({
                message: data.privateMessage.message, 
                isCreatedBy: data.privateMessage.isCreatedBy 
            });
            setTimeout(() => {
                const chatMessages = document.querySelector('.chat-bubbles');
                chatMessages.scrollTop = chatMessages.scrollHeight;
            }, 50);
        }
    },

    chatCreated: (data) => {
        if (!data?.privateMessage?.message) {
            console.error("Invalid chatCreated message:", data);
            return;
        }
        
        currentChatID = data.privateMessage.message.chat_id;
        console.log("Chat created with ID:", currentChatID);
        showChat({
            receiverUserID: currentChatUser?.user_id,
            receiverUserName: currentChatUser?.nickname,
            chatID: currentChatID,
            privateMessages: []
        });
        loadMoreMessages(currentChatID, 1);
    },

    messages: (data) => {
        hasMoreMessages = data.hasMore;

        const chatMessages = document.querySelector('.chat-bubbles');
        const messagesArr = Array.isArray(data.messages) ? data.messages : [];

        if (data.page === 1) {
            displayedMessageIds.clear();
            renderMessages(data.messages, chatMessages);
            messagesArr.forEach(pm => displayedMessageIds.add(pm.message?.id || pm.id));
            chatMessages.scrollTop = chatMessages.scrollHeight;
        } else {
            const oldHeight = chatMessages.scrollHeight;
            const fragment = document.createDocumentFragment();

            // Filter out duplicates
            const newMessages = messagesArr.filter(pm => {
                const id = pm.message?.id || pm.id;
                return !displayedMessageIds.has(id);
            });

            newMessages.slice().reverse().forEach(pm => {
                createChatBubble(pm, fragment, pm.isCreatedBy);
                displayedMessageIds.add(pm.message?.id || pm.id);
            });
            chatMessages.insertBefore(fragment, chatMessages.firstChild);
            chatMessages.scrollTop = chatMessages.scrollHeight - oldHeight;
        }
    },

    typing: (data) => {
        // Check if the typing user is the one we're currently chatting with
        if (data.user_id === currentChatUser?.user_id) {
            showTypingIndicator(data.typing_nickname || "User");
        }
    },

    stopped_typing: (data) => {
        if (data.user_id === currentChatUser?.user_id) {
            hideTypingIndicator();
        }
    }
};

// ------------------------------------------------------------------------------ //

function requestUserListViaWebSocket() {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
            msgType: "getUsers"
        }));
    } else {
        console.error("WebSocket not connected");
    }
}

function renderUserLists(chattedUsers, unchattedUsers) {
    const userList = document.getElementById('userList');
    userList.innerHTML = ''; // Clear previous contents

    // Add Chatted Users Section
    if (chattedUsers?.length > 0) {
        // Header
        const chattedHeader = document.createElement('li');
        chattedHeader.className = 'list-header';
        chattedHeader.textContent = 'Active Chats ▾';
        userList.appendChild(chattedHeader);

        // Users
        chattedUsers.forEach(user => {
            const li = document.createElement('li');
            li.className = 'user-item';
            const statusEmoji = user.isOnline ? ' 🟢' : ' ⚪';
            const notificationEmoji = user.hasUnread ? '🍿' : '';

            li.textContent =  statusEmoji + ' ' + user.nickname + ' ' + notificationEmoji;
            li.addEventListener('click', () => openChatWithUser(user));
            userList.appendChild(li);
        });
    }

    // Add Unchatted Users Section
    if (unchattedUsers?.length > 0) {
        // Header
        const unchattedHeader = document.createElement('li');
        unchattedHeader.className = 'list-header';
        unchattedHeader.textContent = 'Other Users ▾';
        userList.appendChild(unchattedHeader);

        // Users
        unchattedUsers.forEach(user => {
            const li = document.createElement('li');
            li.className = 'user-item';
            const statusEmoji = user.isOnline ? ' 🟢' : ' ⚪';
            const notificationEmoji = user.hasUnread ? '🍿' : '';
            
            li.textContent = statusEmoji + ' ' + user.nickname + ' ' + notificationEmoji;
            li.addEventListener('click', () => openChatWithUser(user));
            userList.appendChild(li);
        });
    }
}

async function openChatWithUser(user) {
    currentChatUser = {
        user_id: user.user_id,
        nickname: user.nickname
    };
     // Reset pagination when opening a new chat
    currentMessagePage = 1;
    try {
        socket.send(JSON.stringify({
            msgType: "getOrCreateChat",
            receiver_user_id: user.user_id,
            clearUnread: true
        }));
        setTimeout(() => {
            requestUserListViaWebSocket();
        }, 100);
        // The rest is handled by the chatCreated/messages handlers
    } catch (error) {
        console.error('Chat error:', error);
    }
}

function showChat(msg) {
    const chatBox = document.getElementById('chatBox');
    if (!chatBox) {
        console.error('No chat container found with id="chatBox" or class="message-list"');
        return;
    }
    chatBox.innerHTML = '';
    // Chat header
    const chatTitle = document.createElement('div');
    chatTitle.classList.add('chat-header');
    chatTitle.textContent = `Chat with ${msg.receiverUserName || 'Unknown'}`;
    chatBox.appendChild(chatTitle);
    // Chat messages container
    const chatMessages = document.createElement('div');
    chatMessages.classList.add('chat-bubbles');
    chatMessages.id = `chat_${msg.chatID}`;
    chatBox.appendChild(chatMessages);
    
    // Render initial messages
    renderMessages(msg.privateMessages, chatMessages);
    // Scroll to bottom to show latest messages
    chatMessages.scrollTop = chatMessages.scrollHeight;

    // Infinite scroll up (throttled)
    let isThrottled = false;

    chatMessages.addEventListener('scroll', () => {
        if (isThrottled) return;
        if (!hasMoreMessages) return; 
        isThrottled = true;
        setTimeout(() => {
          // Load more messages when scrolling to top
            if (chatMessages.scrollTop <= 10) {
                currentMessagePage++;
                loadMoreMessages(currentChatID, currentMessagePage);
            }
            isThrottled = false;
        }, 500); //execute at most once every 500 milliseconds (2 times per second)
    });
}

function loadMoreMessages(chatId, page) {
    if (!chatId) return;    
    socket.send(JSON.stringify({
        msgType: "getMessages",
        privateMessage: {
            message: {
                chat_id: chatId
            }
        },
        page: page
        
    }));
}

function renderMessages(messages, container, options = {}) {
    if (!Array.isArray(messages)) return;
    messages.slice().reverse().forEach(pm => {
        createChatBubble(pm, container, pm.isCreatedBy);
    });
}
function addMessageToChat(pm) {
    const chatMessages = document.querySelector('.chat-bubbles');
    createChatBubble(pm, chatMessages, pm.isCreatedBy);
}

function createChatBubble(pm, chatMessages, isSelf) {
    const msg = pm.message || pm;
    const messageDiv = document.createElement('div');
    messageDiv.className = `message-bubble${isSelf ? ' self' : ''}`;

    if (!chatMessages) {
        showChat({
        receiverUserID: currentChatUser?.user_id,
        receiverUserName: currentChatUser?.nickname,
        chatID: currentChatID,
        privateMessages: []
        });
        chatMessages = document.querySelector('.chat-bubbles');
        if (!chatMessages) {
            console.error('chatMessages container still does not exist!');
            return;
        }
    }

    // Nickname (optional - remove if not needed)
    if (msg.sender_nickname || msg.SenderUsername) {
        const nicknameDiv = document.createElement('div');
        nicknameDiv.className = 'message-nickname';
        nicknameDiv.textContent = msg.sender_nickname || msg.SenderUsername;
        messageDiv.appendChild(nicknameDiv);
    }

    // Message Content
    const contentDiv = document.createElement('div');
    contentDiv.className = 'message-content';
    contentDiv.textContent = msg.content;
    messageDiv.appendChild(contentDiv);

    // Timestamp 
    const timestampDiv = document.createElement('div');
    timestampDiv.className = 'message-timestamp';
    
    try {
        const timestamp = new Date(msg.created_at);
        const year = timestamp.getFullYear();
        const month = String(timestamp.getMonth() + 1).padStart(2, '0'); // Months are 0-indexed
        const day = String(timestamp.getDate()).padStart(2, '0');
        const hours = String(timestamp.getHours()).padStart(2, '0');
        const minutes = String(timestamp.getMinutes()).padStart(2, '0');
        const seconds = String(timestamp.getSeconds()).padStart(2, '0');
        timestampDiv.textContent = `${day}-${month}-${year} ${hours}:${minutes}:${seconds}`;
    } catch {
        timestampDiv.textContent = 'Now';
    }
    
    messageDiv.appendChild(timestampDiv);
    chatMessages.appendChild(messageDiv);
}

function sendMessage() {
    const messageInput = document.getElementById('messageInput');
    const message = messageInput.value.trim();
    if (!message || !currentChatID || !currentChatUser?.user_id) {
        console.error("Missing required fields for message");
        return;
    }
    const sessionToken = localStorage.getItem('sessionToken');
    if (!sessionToken) {
        console.error("No session token");
        return;
    }
    // Construct the WebSocket message as expected by the backend
    const wsMessage = {
        msgType: "sendMessage",
        receiver_user_id: currentChatUser.user_id,
        PrivateMessage: {
            Message: {
                chat_id: currentChatID,
                content: message,
                sender_nickname: currentChatUser.nickname
            }
        },
        token: sessionToken 
    };
    try {
        if (socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(wsMessage));
            messageInput.value = '';
        } else {
            console.error("WebSocket not connected");
            // Optionally, try to reconnect
        }
    } catch (error) {
        console.error("Message send error:", error);
        // Optionally, show error to user
    }
}

// ------------------------------------------------------------------------------ //

function showTypingIndicator(nickname = "User") {
    const typingDiv = document.getElementById('typing-indicator');
    if (typingDiv) {
        typingDiv.innerHTML = `${nickname} is typing <span class="typing-dots"><span class="popcorn">🍿</span>
            <span class="popcorn">🍿</span><span class="popcorn">🍿</span></span>`;
        typingDiv.style.display = 'block';
    }
}

function hideTypingIndicator() {
    const typingDiv = document.getElementById('typing-indicator');
    if (typingDiv) {
        typingDiv.style.display = 'none';
    }
}


