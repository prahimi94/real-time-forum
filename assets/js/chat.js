let ws;
let onlineUsers;
let usernames = [];

async function fetchOnlineUsers() {
  try {
    const response = await fetch("/api/online-users");
    if (!response.ok) {
      console.error("Failed to fetch online users");
      return;
    }

    usernames = await response.json();
    updateOnlineUsersList(usernames);
  } catch (error) {
    console.error("Error fetching online users:", error);
  }
}

// Function to update the online users list in the frontend
function updateOnlineUsersList(usernames) {
  const onlineUsersList = document.getElementById("online-users-list");
  onlineUsersList.textContent = "";  // Clear the current list

  if (usernames.length === 0 || (usernames.length === 1 && usernames[0] === loggedInUser.name)) {
    onlineUsersList.textContent = "It's just you here.";
    return;
  }

  // Populate the list with usernames
  usernames.forEach((username) => {
    if (username === loggedInUser.name) return;

    const li = document.createElement("li");
    li.textContent = username;
    li.style.cursor = "pointer";
    li.onclick = () => openPrivateChat(username);
    onlineUsersList.appendChild(li);
  });
}

function openPrivateChat(username) {
  privateRecipient = username;
  document.getElementById("messages").style.display = "block";
  const chatHeader = document.getElementById("chat-header");
  chatHeader.textContent = `Chat with ${username}`;
  document.getElementById("messageInput").placeholder = `Type a message to ${username}`;

  // Send a message to the server to initiate or check the chat
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: "private_chat", recipient: username }));
  }
}


function connect() {
  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = function () {
    console.log("Online: Connected to WebSocket server");
    updateOnlineUsersList(onlineUsers);
  };

  ws.onmessage = function (event) {
    let messageDisplay = document.getElementById("messages");
    let message = event.data;

    try {
      // Check if the message is a JSON array (online users list)
      onlineUsers = JSON.parse(message);
      if (Array.isArray(onlineUsers)) {
        updateOnlineUsersList(onlineUsers);
        return;
      }
    } catch (e) {
      // Not a JSON array, proceed with normal message handling
    }

    // Append the message to the chatbox
    let messageElement = document.createElement("p");
    messageElement.textContent = message;
    messageDisplay.appendChild(messageElement);

    // Scroll to the bottom of the chatbox
    messageDisplay.scrollTop = messageDisplay.scrollHeight;
  };

  ws.onclose = function () {
    console.log("Offline: WebSocket connection closed, retrying...");
    updateOnlineUsersList(onlineUsers);
    setTimeout(connect, 1000); // Reconnect after 1 second
  };

  ws.onerror = function (error) {
    console.error("Offline: WebSocket error:", error);
  };
}

function sendMessage() {
  let input = document.getElementById("messageInput");
  let message = input.value.trim(); // Remove leading/trailing whitespace
  if (message.length === 0) return; // Do not accept empty messages

  if (privateRecipient) {
    ws.send(JSON.stringify({ type: "message_content", content: message }));
  }
  input.value = "";
}