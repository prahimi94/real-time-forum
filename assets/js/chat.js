let ws;
let onlineUsernames = [];
let allUserData = {};
let isProcessingMessages = false;
let anotherUserClicked = false;
let chatboxOpen = false;
let openChat = [ID = null, recipient = null];

async function fetchData(apiEndpoint) {
  try {
    const response = await fetch(apiEndpoint);
    if (!response.ok) {
      throw new Error(`Error fetching ${apiEndpoint}: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error(error);
    return null;
  }
}

async function showAllChatUsers() {
  const chatUsersList = document.getElementById("chat-users-list");
  chatUsersList.textContent = ""; // Clear the current list
  const messageDisplay = document.getElementById("message-display");
  const chatbox = document.getElementById("chatbox");
  const chatHeader = document.getElementById("chat-header");
  const messages = document.getElementById("messages");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");


  allUserData.forEach((user) => {
    if (user.username === loggedInUser.username) return;

    const li = document.createElement("li");
    const isOnline = onlineUsernames.includes(user.username);

    if (isOnline) {
      li.classList.add("isOnline");
      li.textContent = `${user.username} (Online)`;
      messageInput.style.display = "block";
      sendButton.style.display = "block";
    } else {
      li.textContent = user.username;
      messageInput.style.display = "none";
      sendButton.style.display = "none";
    }

    li.onclick = async (e) => {
      e.preventDefault();

      isProcessingMessages = false;
      anotherUserClicked = true;
      messageDisplay.innerHTML = "";

      if (chatboxOpen || (chatboxOpen && openChat[1] === user.username)) {
        chatboxOpen = false;
        chatbox.style.display = "none";
      } else if (!chatboxOpen || (chatboxOpen && openChat[1] !== user.username)) {
        chatboxOpen = true;
        chatbox.style.display = "block";
        if (isOnline) {
          // Fetch the chat ID and load the chat messages
          const chatID = await getChatID(loggedInUser.username, user.username);
          if (chatID) {
            openChat = [chatID, user.username];
            messageDisplay.scrollTop = messageDisplay.scrollHeight;
            await showMessages(chatID, user.username);
            messageInput.style.display = "block";
            sendButton.style.display = "block";
          } else if (chatID === 0 || !chatID) {
            openChat = [null, null];
            messageInput.style.display = "none";
            sendButton.style.display = "none";
            chatHeader.textContent = `Chat with ${user.username}`;
            messageDisplay.innerHTML = "No chat history.";
          } else {
            console.error("Failed to retrieve chat ID");
          }
          await handlePrivateChat(loggedInUser.username, user.username);

        } else {

          chatbox.style.display = "block";
          messages.style.display = "block";
          messageDisplay.style.display = "block";
          messageDisplay.innerHTML = "";

          // Fetch the chat ID and load the chat messages
          const chatID = await getChatID(loggedInUser.username, user.username);
          if (chatID) {
            openChat = [chatID, user.username];
            await showMessages(chatID, user.username);
            messageInput.style.display = "none";
            sendButton.style.display = "none";
          } else if (chatID === 0 || !chatID) {
            openChat = [null, null];
            messageInput.style.display = "none";
            sendButton.style.display = "none";
            chatHeader.textContent = `Chat with ${user.username}`;
            messageDisplay.innerHTML = "No chat history.";
          } else {
            console.error("Failed to retrieve chat ID");
          }
        }
      }
    }

    chatUsersList.appendChild(li);

  });

}

function connect() {
  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = async function () {
    console.log("Online: Connected to WebSocket server");
    onlineUsernames = await fetchData("/api/online-users");
    allUserData = await fetchData("/api/users");
    await showAllChatUsers();
  };

  ws.onmessage = async function (event) {
    try {
      onlineUsernames = await fetchData("/api/online-users");
      allUserData = await fetchData("/api/users");
      await showAllChatUsers();

      const message = JSON.parse(event.data);

      const messageDisplay = document.getElementById("message-display");

      if (message.type === "message_content") {
        // Show notification only to the recipient
        if (loggedInUser.username === message.recipient) {
          const res = {
            success: true,
            message: `You have a new message from ${message.sender}`,
            data: message,
          };
          showToast(res);
        }

        const chatID = await getChatID(message.sender, message.recipient);
        if (chatID && openChat[0] === chatID && (openChat[1] === message.recipient || openChat[1] === message.sender) && (onlineUsernames.includes(message.sender) && onlineUsernames.includes(message.recipient))) {
          const messageElement = document.createElement("p");
          const details = document.createElement("div");
          details.classList.add("details");
          messageElement.classList.add("msg");
          details.textContent = `${message.sender} (${message.timestamp})`;
          messageElement.textContent = message.content;

          if (loggedInUser.username === message.sender) {
            messageElement.classList.add("from-me");
            details.classList.add("my-details");
          }

          messageDisplay.appendChild(messageElement);
          messageDisplay.appendChild(details);
          messageDisplay.scrollTop = messageDisplay.scrollHeight;

          document.getElementById("messageInput").style.display = "block";
          document.getElementById("send-btn").style.display = "block";
        }
      }
    } catch (error) {
      console.error("Failed to parse WebSocket message:", event.data, error);
    }
  };

  ws.onclose = function () {
    console.log("Offline: WebSocket connection closed, retrying...");
    setTimeout(connect, 1000); // Reconnect after 1 second
  };

  ws.onerror = function (error) {
    console.error("Offline: WebSocket error:", error);
  };
}

async function handlePrivateChat(senderUsername, recipientUsername) {
  if (!onlineUsernames.includes(recipientUsername) || !onlineUsernames.includes(senderUsername)) return;

  document.getElementById("chatbox").style.display = "block";
  document.getElementById("messages").style.display = "block";
  const chatHeader = document.getElementById("chat-header");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");
  
  messageInput.style.display = "block";
  sendButton.style.display = "block";

  if (loggedInUser.username === senderUsername) {
    chatHeader.textContent = `Chat with ${recipientUsername}`;
    messageInput.placeholder = `Type a message to ${recipientUsername}`;
    sendButton.onclick = () => sendMessage(senderUsername, recipientUsername);
  } else if (loggedInUser.username === recipientUsername) {
    chatHeader.textContent = `Chat with ${senderUsername}`;
    messageInput.placeholder = `Type a message to ${senderUsername}`;
    sendButton.onclick = () => sendMessage(recipientUsername, senderUsername);
  }

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
}

async function showMessages(chatID, recipientUsername) {
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");
  const chatHeader = document.getElementById("chat-header");
  const messageDisplay = document.getElementById("message-display");

  chatHeader.textContent = `Chat with ${recipientUsername}`;

  let messages = await fetchData(`/api/chat-messages/${chatID}`);
  if (!messages || messages.length === 0) {
    messageDisplay.textContent = "Chat started but no messages yet.";
    return;
  }

  // Reverse the messages to show the latest ones at the bottom
  messages.reverse();

  isProcessingMessages = false;
  anotherUserClicked = false;

  // Throttle to process msgs in batches of 10
  const batchSize = 10;
  let loadedMessages = []; // Keep track of already loaded messages

  function processBatch() {
    if (isProcessingMessages || anotherUserClicked) return; // Prevent multiple triggers
    isProcessingMessages = true;

    const batch = messages.splice(0, batchSize); // Get the next batch of messages
    loadedMessages = [...batch, ...loadedMessages]; // Add to the loaded messages

    batch.forEach((message) => {
      if (anotherUserClicked) return;

      const msg = JSON.parse(message.content);
      const messageElement = document.createElement("p");
      const details = document.createElement("div");
      details.classList.add("details");
      messageElement.classList.add("msg");
      details.textContent = `${msg.sender} (${msg.timestamp})`;
      messageElement.textContent = msg.content;

      if (loggedInUser.username === msg.sender) {
        messageElement.classList.add("from-me");
        details.classList.add("my-details");
      }

      messageDisplay.prepend(details);
      messageDisplay.prepend(messageElement);
    });

    isProcessingMessages = false;
  }

  // Initial batch load (latest 10 messages)
  processBatch();

  // Scroll to the bottom of the chatbox to show the latest messages
  messageDisplay.scrollTop = messageDisplay.scrollHeight;

  // Add scroll event listener to load more messages when scrolling near the top
  messageDisplay.addEventListener("scroll", () => {
    if (messageDisplay.scrollTop === 0 && messages.length > 0) {
      const previousHeight = messageDisplay.scrollHeight; // Store the current height
      processBatch();
      // Maintain the scroll position after loading more messages
      messageDisplay.scrollTop = messageDisplay.scrollHeight - previousHeight;
    }
  });

  if (!onlineUsernames.includes(recipientUsername)) {
    messageInput.style.display = "none";
    sendButton.style.display = "none";
  }

}

async function getChatID(senderUsername, recipientUsername) {
  const data = await fetchData(`/api/get-chat-id?sender=${senderUsername}&recipient=${recipientUsername}`);
  return data.chatID; // Return the chatID from the response
}

function sendMessage(senderUsername, recipientUsername) {
  isProcessingMessages = false;
  let input = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");
  let message = input.value.trim(); // Remove leading/trailing whitespace
  if (message.length === 0) return; // Do not accept empty messages

  sendButton.disabled = true; // Prevent spam sending

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
  sendButton.disabled = false; // Re-enable

}
