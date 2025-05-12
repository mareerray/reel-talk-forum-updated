// WebSocket connection
const socket = new WebSocket('ws://yourserver.com/ws');

// Handle WebSocket errors and connection closure
socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};

socket.onclose = () => {
    console.warn('WebSocket connection closed');
};

// Listen for messages
socket.onmessage = (event) => {
    const message = JSON.parse(event.data);

    if (message.msgType === 'updateClients') {
        // Fetch the updated user list from the server
        fetch('/api/users/list')
            .then(response => response.json())
            .then(users => {
                const userListContainer = document.querySelector('.user-list-container');
                userListContainer.innerHTML = ''; // Clear existing list

                users.forEach(user => {
                    const userElement = document.createElement('div');
                    userElement.className = 'user-item';
                    userElement.setAttribute('data-user-id', user.id); // Add user ID for status updates
                    userElement.textContent = user.nickname; // Display user nickname
                    userListContainer.appendChild(userElement);
                });
            })
            .catch(error => console.error('Error fetching user list:', error));
    }
    else if (message.msgType === 'newMessage') {
        // Handle new message
        const chatContainer = document.querySelector('.chat-container');
        const messageElement = document.createElement('div');
        messageElement.className = 'chat-message';
        messageElement.textContent = `${message.sender}: ${message.content}`;
        chatContainer.appendChild(messageElement);
    }
    else if (message.msgType === 'userStatus') {
        // Handle user status update
        const statusElement = document.querySelector(`.user-item[data-user-id="${message.userId}"]`);
        if (statusElement) {
            statusElement.classList.toggle('online', message.status === 'online');
            statusElement.classList.toggle('offline', message.status === 'offline');
        }
    }
};

// Send a message
function sendMessage() {
    const messageInput = document.querySelector('.message-input');
    const message = {
        msgType: 'newMessage',
        content: messageInput.value,
        sender: 'currentUserId' // Replace with actual user ID
    };
    socket.send(JSON.stringify(message));
    messageInput.value = ''; // Clear input after sending
}