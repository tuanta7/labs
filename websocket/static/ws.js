let ws = null;
const statusEl = document.getElementById('status');
const messagesEl = document.getElementById('messages');
const connectBtn = document.getElementById('connectBtn');
const disconnectBtn = document.getElementById('disconnectBtn');
const messageInput = document.getElementById('messageInput');
const sendBtn = document.getElementById('sendBtn');

function updateStatus(connected) {
    if (connected) {
        statusEl.textContent = 'Connected';
        statusEl.className = 'connected';
        connectBtn.disabled = true;
        disconnectBtn.disabled = false;
        messageInput.disabled = false;
        sendBtn.disabled = false;
    } else {
        statusEl.textContent = 'Disconnected';
        statusEl.className = 'disconnected';
        connectBtn.disabled = false;
        disconnectBtn.disabled = true;
        messageInput.disabled = true;
        sendBtn.disabled = true;
    }
}

function addMessage(text, isError = false) {
    const msgEl = document.createElement('div');
    msgEl.className = isError ? 'message error' : 'message';
    msgEl.textContent = new Date().toLocaleTimeString() + ' - ' + text;
    messagesEl.appendChild(msgEl);
    messagesEl.scrollTop = messagesEl.scrollHeight;
}

function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const url = protocol + '//' + window.location.host + '/ws';

    try {
        ws = new WebSocket(url);

        ws.onopen = function() {
            addMessage('Connected to WebSocket');
            updateStatus(true);
        };

        ws.onmessage = function(event) {
            addMessage('Received: ' + event.data);

            // Auto-respond to PING with PONG
            if (event.data === 'PING') {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    ws.send('PONG');
                    addMessage('Sent: PONG (auto-response)');
                }
            }
        };

        ws.onerror = function(error) {
            addMessage('WebSocket error: ' + error, true);
        };

        ws.onclose = function() {
            addMessage('Disconnected from WebSocket');
            updateStatus(false);
        };
    } catch (error) {
        addMessage('Connection error: ' + error.message, true);
    }
}

function disconnect() {
    if (ws) {
        ws.close();
        ws = null;
    }
}

function sendMessage() {
    const message = messageInput.value.trim();
    if (message && ws && ws.readyState === WebSocket.OPEN) {
        ws.send(message);
        addMessage('Sent: ' + message);
        messageInput.value = '';
    }
}

messageInput.addEventListener('keypress', function(event) {
    if (event.key === 'Enter') {
        sendMessage();
    }
});

updateStatus(false);