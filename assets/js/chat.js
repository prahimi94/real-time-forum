async function fetchOnlineUsers() {
    try {
        const response = await fetch('/api/online-users');
        if (!response.ok) {
            console.error("Failed to fetch online users");
            return;
        }

        const usernames = await response.json();
        const onlineUsersList = document.getElementById("online-users-list");

        // Clear the current list
        onlineUsersList.textContent = "";

        if (usernames.length === 0) {
            onlineUsersList.textContent = 'No other users online';
            return;
        }

        // Populate the list with usernames
        usernames.forEach(username => {
            const li = document.createElement("li");
            li.textContent = username;
            onlineUsersList.appendChild(li);
        });
    } catch (error) {
        console.error("Error fetching online users:", error);
    }
}

/* WEBSOCKET FOR CHAT */
let ws;

function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = function () {
        console.log("Connected to WebSocket server");
    };

    ws.onmessage = function (event) {
        let messageDisplay = document.getElementById("messages");
        let message = event.data;

        // Append the message to the chatbox
        let messageElement = document.createElement("p");
        messageElement.textContent = message;
        messageDisplay.appendChild(messageElement);

        // Scroll to the bottom of the chatbox
        messageDisplay.scrollTop = messageDisplay.scrollHeight;
    };

    ws.onclose = function () {
        console.log("WebSocket connection closed, retrying...");
        setTimeout(connect, 1000); // Reconnect after 1 second
    };

    ws.onerror = function (error) {
        console.error("WebSocket error:", error);
    };
}

function sendMessage() {
    let input = document.getElementById("messageInput");
    let message = input.value;
    ws.send(message);
    input.value = "";
}
/* END OF WEBSOCKET FOR CHAT */