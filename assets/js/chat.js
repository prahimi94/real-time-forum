let ws;
let onlineUsers = [];
let onlineUsernames = [];
let allUserData = {};
let isProcessingMessages = false;

async function fetchOnlineUsers() {
  try {
    const response = await fetch("/api/online-users");
    if (!response.ok) {
      console.error("Failed to fetch online users");
      return;
    }

    onlineUsernames = await response.json();
    updateOnlineUsersList(onlineUsernames);
  } catch (error) {
    console.error("Error fetching online users:", error);
  }
}

// Function to update the online users list in the frontend
function updateOnlineUsersList(onlineUsernames) {
  const onlineUsersList = document.getElementById("online-users-list");
  onlineUsersList.textContent = ""; // Clear the current list

  if (
    onlineUsernames.length === 0 ||
    (onlineUsernames.length === 1 && onlineUsernames[0] === loggedInUser.name)
  ) {
    onlineUsersList.textContent = "It's just you here.";
    return;
  }

  // Populate the list with onlineUsernames
  onlineUsernames.forEach((onlineUsername) => {
    if (onlineUsername === loggedInUser.name) return;

    const li = document.createElement("li");
    li.textContent = onlineUsername;
    li.style.cursor = "pointer";
    li.onclick = () => handlePrivateChat(loggedInUser.name, onlineUsername);
    onlineUsersList.appendChild(li);
  });
}

async function fetchAllUsers() {
  try {
    const response = await fetch("/api/users");
    if (!response.ok) {
      console.error("Failed to fetch all users");
      return;
    }

    allUserData = await response.json();
  
    const offlineUsersList = document.getElementById("offline-users-list");
    offlineUsersList.textContent = ""; // Clear the current list
  
    allUserData.forEach((user) => {
      if (!onlineUsernames.includes(user.name)) {
        const li = document.createElement("li");
        li.textContent = user.name;
        offlineUsersList.appendChild(li);
      }
    });
  } catch (error) {
    console.error("Error fetching all users:", error);
  }
}

function connect() {
  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = function () {
    console.log("Online: Connected to WebSocket server");
    fetchOnlineUsers(); // Fetch the latest online users when connected
    fetchAllUsers(); // Fetch all users to populate the list
  };

  ws.onmessage = function (event) {
    try {
      const message = JSON.parse(event.data);
      //console.log("message: ", message);
      if (message.type === "message_content") {
        handlePrivateChat(message.sender, message.recipient);
      }

      if (Array.isArray(onlineUsers)) {
        fetchOnlineUsers()
        fetchAllUsers();
        return;
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

// Handle private chat initiation
async function handlePrivateChat(senderUsername, recipientUsername) {
  document.getElementById("messages").style.display = "block";
  const chatHeader = document.getElementById("chat-header");
  chatHeader.textContent = `Chat with ${recipientUsername}`;
  document.getElementById(
    "messageInput"
  ).placeholder = `Type a message to ${recipientUsername}`;

  const sendButton = document.getElementById("send-btn");
  sendButton.onclick = () => sendMessage(senderUsername, recipientUsername);

  // Send a message to the server to initiate or check the chat
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(
      JSON.stringify({
        type: "private_chat",
        sender: senderUsername,
        recipient: recipientUsername,
      })
    );
  }

  // Fetch the chat ID and load the chat messages
  const chatID = await getChatIDForUsers(senderUsername, recipientUsername);
  if (chatID) {
    if (
      loggedInUser.name === senderUsername ||
      loggedInUser.name === recipientUsername
    ) {
      fetchChatMessages(chatID);
    }
  } else {
    console.error("Failed to retrieve chat ID");
  }
}

async function fetchChatMessages(chatID) {
  try {
    const response = await fetch(`/api/chat-messages/${chatID}`);
    if (!response.ok) {
      console.error("Failed to fetch chat messages");
      return;
    }

    const messages = await response.json();
    const messageDisplay = document.getElementById("message-display");
    messageDisplay.textContent = ""; // Clear existing messages

    isProcessingMessages = true;

    // Throttle to process msgs in batches of 10
    const batchSize = 10;
    function processBatch() {
      const batch = messages.splice(0, batchSize); // Get the next batch of messages

      batch.forEach((message) => {
        const parsedContent = JSON.parse(message.content);
        const messageElement = document.createElement("p");
        messageElement.textContent = `[${parsedContent.timestamp}] ${parsedContent.sender}: ${parsedContent.content}`;
        messageDisplay.appendChild(messageElement);
      });

      // Scroll to the bottom of the chatbox
      messageDisplay.scrollTop = messageDisplay.scrollHeight;

      if (batch.length > 0) {
        // If there are more messages, process the next batch after a delay
        setTimeout(processBatch, 500); // Adjust delay as needed
      } else {
        isProcessingMessages = false; // Mark processing as complete
      }
    }
    processBatch();
  } catch (error) {
    console.error("Error fetching chat messages:", error);
  }
}

async function getChatIDForUsers(senderUsername, recipientUsername) {
  try {
    const response = await fetch(
      `/api/get-chat-id?sender=${senderUsername}&recipient=${recipientUsername}`
    );
    if (!response.ok) {
      console.error("Failed to fetch chat ID");
      return null;
    }

    const data = await response.json();
    return data.chatID; // Return the chatID from the response
  } catch (error) {
    console.error("Error fetching chat ID:", error);
    return null;
  }
}

function sendMessage(senderUsername, recipientUsername) {
  let input = document.getElementById("messageInput");
  let message = input.value.trim(); // Remove leading/trailing whitespace
  if (message.length === 0) return; // Do not accept empty messages

  if (recipientUsername) {
    ws.send(
      JSON.stringify({
        type: "message_content",
        sender: senderUsername,
        recipient: recipientUsername,
        content: message,
        timestamp: new Date().toLocaleString(),
      })
    );
  }
  input.value = "";
}
