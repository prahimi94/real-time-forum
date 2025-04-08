let ws;
let onlineUsers = [];
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
  onlineUsersList.textContent = ""; // Clear the current list

  if (
    usernames.length === 0 ||
    (usernames.length === 1 && usernames[0] === loggedInUser.name)
  ) {
    onlineUsersList.textContent = "It's just you here.";
    return;
  }

  // Populate the list with usernames
  usernames.forEach((username) => {
    if (username === loggedInUser.name) return;

    const li = document.createElement("li");
    li.textContent = `${username} is online!`;
    li.style.cursor = "pointer";
    li.onclick = () => handlePrivateChat(username);
    onlineUsersList.appendChild(li);
  });
}

function connect() {
  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = function () {
    console.log("Online: Connected to WebSocket server");
    fetchOnlineUsers(); // Fetch the latest online users when connected
  };

  ws.onmessage = function (event) {
    try {
      const message = JSON.parse(event.data);
      console.log("message: ", message);
      if (message.type) {
        handleMessage(message);
      } else {
        console.warn("Unknown message format:", message);
      }
    } catch (error) {
      console.error("Failed to parse WebSocket message:", event.data, error);
    }
  };

  ws.onclose = function () {
    console.log("Offline: WebSocket connection closed, retrying...");
    onlineUsers = []; // Clear the onlineUsers list on disconnect
    updateOnlineUsersList(onlineUsers);
    setTimeout(connect, 1000); // Reconnect after 1 second
  };

  ws.onerror = function (error) {
    console.error("Offline: WebSocket error:", error);
  };
}

// Centralized message handler
function handleMessage(message) {
  switch (message.type) {
    /* case "online_users_list":
      console.log("Received online users list:", message.data);
      handleOnlineUsersList(message.data);
      break; */

    case "private_chat":
      handlePrivateChat(message.data);
      break;

    case "message_content":
      handleMessageContent(message);
      break;

    default:
      console.warn("Unhandled message type:", message.type);
  }
}

// Handle online users list
/* function handleOnlineUsersList(data) {
  console.log("handleOnlineUsersList(data):", data);
  if (Array.isArray(data)) {
    onlineUsers = data; // Update the onlineUsers array
    updateOnlineUsersList(onlineUsers);
  } else {
    console.warn("Invalid data for online_users_list:", data);
  }
} */

// Handle private chat initiation
function handlePrivateChat(username) {
  privateRecipient = username;
  document.getElementById("messages").style.display = "block";
  const chatHeader = document.getElementById("chat-header");
  chatHeader.textContent = `Chat with ${username}`;
  document.getElementById(
    "messageInput"
  ).placeholder = `Type a message to ${username}`;

  // Send a message to the server to initiate or check the chat
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: "private_chat", recipient: username }));
  }
}

// Handle incoming chat messages
function handleMessageContent(message) {
  if (loggedInUser.name === message.recipient || loggedInUser.name === message.sender) {
    const messageDisplay = document.getElementById("message-display");
    messageDisplay.style.display = "block";
    const messageElement = document.createElement("p");
    messageElement.textContent = `[${message.timestamp}] ${message.sender}: ${message.content}`;
    messageDisplay.appendChild(messageElement);

    // Scroll to the bottom of the chatbox
    messageDisplay.scrollTop = messageDisplay.scrollHeight;
  }
}

function sendMessage() {
  let input = document.getElementById("messageInput");
  let message = input.value.trim(); // Remove leading/trailing whitespace
  if (message.length === 0) return; // Do not accept empty messages

  if (privateRecipient) {
    ws.send(
      JSON.stringify({
        type: "message_content",
        sender: loggedInUser.name,
        recipient: privateRecipient,
        content: message,
        timestamp: new Date().toLocaleString(),
      })
    );
  }
  input.value = "";
}
