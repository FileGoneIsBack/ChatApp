let socket;
let lastMessageID = 0;
function setupWebSocket() {
    socket = new WebSocket("wss://frosted.dev/api/websockets/messages");
    socket.onmessage = function(event) {
        let data = JSON.parse(event.data);
        console.log("New message:", data);
        // Render only new messages
        renderMessages(data);
    };

    socket.onopen = function() {
        console.log("WebSocket connection opened.");
        // Request messages since the last message ID
        socket.send(JSON.stringify({ type: 'requestNewMessages', lastMessageID }));
    };

    socket.onclose = function(event) {
        console.log("WebSocket connection closed.", event);
        // Reconnect if needed
        setTimeout(setupWebSocket, 500);
    };

    socket.onerror = function(error) {
        console.error("WebSocket error:", error);
        // Handle reconnection if needed
        socket.close(); 
    };
}

function renderMessages(messages) {
    const messagesContainer = document.getElementById('messageBody');
    if (!messagesContainer) {
        console.error('Element with id "messageBody" not found.');
        return;
    }
    messages.forEach(message => {
        if (message.id > lastMessageID) {
        // Create message elements
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${message.isMine ? 'self' : 'received'}`;
        
        const messageWrapperDiv = document.createElement('div');
        messageWrapperDiv.className = 'message-wrapper';
        
        const messageContentDiv = document.createElement('div');
        messageContentDiv.className = 'message-content';
        messageContentDiv.innerHTML = `<span>${message.content}</span>`;
        
        const messageOptionsDiv = document.createElement('div');
        messageOptionsDiv.className = 'message-options';
        
        const avatarDiv = document.createElement('div');
        avatarDiv.className = 'avatar avatar-sm';
        avatarDiv.innerHTML = `<img alt="" src="./..//_assets/media/6.png">`;
        
        const messageDateSpan = document.createElement('span');
        messageDateSpan.className = 'message-date';
        messageDateSpan.textContent = message.sentAt;
        
        const dropdownDiv = document.createElement('div');
        dropdownDiv.className = 'dropdown';
        dropdownDiv.innerHTML = `
            <a class="text-muted" href="#" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                <!-- SVG icon -->
            </a>
            <div class="dropdown-menu">
                <!-- Dropdown items -->
            </div>
        `;
        
        // Append elements
        messageOptionsDiv.appendChild(avatarDiv);
        messageOptionsDiv.appendChild(messageDateSpan);
        messageOptionsDiv.appendChild(dropdownDiv);
        
        messageWrapperDiv.appendChild(messageContentDiv);
        
        messageDiv.appendChild(messageWrapperDiv);
        messageDiv.appendChild(messageOptionsDiv);
        
        // Append the new message div to the container
        messagesContainer.appendChild(messageDiv);
        lastMessageID = message.id;
        }
    });

}

// Initialize WebSocket connection on window load
window.onload = function() {
    setupWebSocket();
};

function scrollToBottom() {
    const messagesContainer = document.getElementById('messageBody');
    // Add a scroll anchor
    const scrollAnchor = document.createElement('div');
    scrollAnchor.className = 'chat-finished';
    messagesContainer.appendChild(scrollAnchor);

    // Scroll to the anchor
    scrollAnchor.scrollIntoView({
        block: 'end',
        behavior: 'smooth'
    });
};

document.addEventListener('DOMContentLoaded', function() {
    // Get references to the textarea and send button
    const messageInput = document.getElementById('messageInput');
    const sendButton = document.querySelector('.send-icon');

    // Event listener for send button
    sendButton.addEventListener('click', function() {
        const message = messageInput.value.trim();
        scrollToBottom();
        if (message) {
            // Send the message to the server
            sendMessage(message);

            // Clear the textarea after sending
            messageInput.value = '';
        } else {
            console.error('Message cannot be empty');
        }
    });
});

function sendMessage(message) {
    fetch('/api/functions/send', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            message: message
        })
    })
    .then(response => {
        if (response.ok) {
            console.log('Message sent successfully');
        } else {
            console.error('Failed to send message');
        }
    })
    .catch(error => console.error('Error sending message:', error));
}
