/* WEBSOCKET FOR CHAT */
let ws;

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
      const onlineUsers = JSON.parse(message);
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
  let message = input.value;
  ws.send(message);
  input.value = "";
}
/* END OF WEBSOCKET FOR CHAT */

async function fetchOnlineUsers() {
  try {
    const response = await fetch("/api/online-users");
    if (!response.ok) {
      console.error("Failed to fetch online users");
      return;
    }

    const usernames = await response.json();
    updateOnlineUsersList(usernames);
  } catch (error) {
    console.error("Error fetching online users:", error);
  }
}

// Function to update the online users list in the frontend
function updateOnlineUsersList(usernames) {
  const onlineUsersList = document.getElementById("online-users-list");

  // Clear the current list
  onlineUsersList.textContent = "";

  if (usernames.length === 0 || (usernames.length === 1 && usernames[0] === "OnlineUsername")) { //TODO: "OnlineUsername" to be changed with correct var
    onlineUsersList.textContent = "It's just you here.";
    return;
  }

  // Populate the list with usernames
  usernames.forEach((username) => {
    const li = document.createElement("li");
    li.textContent = username;
    onlineUsersList.appendChild(li);
  });
}
