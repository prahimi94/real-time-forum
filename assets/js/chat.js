let ws;
let onlineUsernames = [];
let activeChat = [(ID = null), (recipient = null)];
let loadedMessages = []; // Keep track of loaded messages for func showMessage
let fetchedMsg;

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

// Handles listing all users
async function ShowAllChatUsers(users) {
  if (!Array.isArray(users)) {
    users = Object.values(users);
  }
  try {
    const chatUsersList = document.getElementById("chat-users-list");
    chatUsersList.innerHTML = "";
    const chatbox = document.getElementById("chatbox");
    const messageInput = document.getElementById("messageInput");
    const sendButton = document.getElementById("send-btn");

    for (const user of users) {
      if (user.username === loggedInUser.username) continue;
      const chatID = await getChatID(loggedInUser.username, user.username);
      const li = document.createElement("li");

      if (user.isOnline) {
        li.classList.add("isOnline");
        li.textContent = `${user.username} 👋`;

        li.onclick = null;
        li.onclick = async () =>
          openChatbox(chatID, loggedInUser.username, user.username);

        if (chatbox.style.display === "block") {
          messageInput.style.display = "block";
          sendButton.style.display = "block";
          await handlePrivateChat(loggedInUser.username, user.username);
          // Listen to input events (typing) in the messageInput field
          messageInput.removeEventListener("input", (event) =>
            handleTyping(event, loggedInUser.username, user.username)
          ); // Remove previous event listener to avoid duplicates
          messageInput.addEventListener("input", (event) =>
            handleTyping(event, loggedInUser.username, user.username)
          );
        }
      } else {
        li.textContent = user.username;
        li.onclick = null;
        li.onclick = async () =>
          openChatbox(chatID, loggedInUser.username, user.username);

        chatbox.style.display = "none";
        stopTyping(loggedInUser.username, user.username);
      }
      chatUsersList.appendChild(li);
    }
  } catch (error) {
    console.error("Error fetching all users:", error);
  }
}

// Handles when username is clicked
async function openChatbox(chatID, senderUsername, recipientUsername) {
  if (chatID === 0 || chatID === null) {
    chatID = await getChatID(senderUsername, recipientUsername);
  }

  const messageDisplay = document.getElementById("message-display");
  const chatbox = document.getElementById("chatbox");
  const chatHeader = document.getElementById("chat-header");
  const messages = document.getElementById("messages");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");
  loadedMessages = []; // Reset the loaded messages array
  messageDisplay.innerHTML = "";

  if (chatbox.style.display === "block") {
    activeChat = [null, null];
    chatbox.style.display = "none";
    stopTyping(senderUsername, recipientUsername); // stop typing when chatbox is closed
  } else if (chatbox.style.display === "none") {
    activeChat = [chatID, recipientUsername];
    chatbox.style.display = "block";
    messages.style.display = "block";
    messageDisplay.style.display = "block";
    messageDisplay.scrollTop = messageDisplay.scrollHeight;
    messageDisplay.innerHTML = "";
    chatHeader.textContent = `Chat with ${recipientUsername}`;
    messageInput.style.display = "block";
    sendButton.style.display = "block";
    // Listen to input events (typing) in the messageInput field
    messageInput.removeEventListener("input", (event) =>
      handleTyping(event, loggedInUser.username, user.username)
    ); // Remove previous event listener to avoid duplicates
    messageInput.addEventListener("input", (event) =>
      handleTyping(event, senderUsername, recipientUsername)
    );

    if (onlineUsernames.has(recipientUsername)) {
      if (chatID) {
        await showMessages(chatID, recipientUsername);
      } else if (chatID === 0 || !chatID) {
        messageDisplay.innerHTML = "No chat history.";
      }
      await handlePrivateChat(senderUsername, recipientUsername);
    } else {
      //if user is offline
      messageInput.style.display = "none";
      sendButton.style.display = "none";
      stopTyping(senderUsername, recipientUsername);
      if (chatID) {
        await showMessages(chatID, recipientUsername);
      } else if (chatID === 0 || !chatID) {
        chatHeader.textContent = `Chat with ${recipientUsername}`;
        messageDisplay.innerHTML = "";
        messageDisplay.innerHTML = "No chat history.";
      }
    }
  }
}

/* --- WEBSOCKET CONNECTION --- */
function connect() {
  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = async function () {
    console.log("Online: Connected to WebSocket server", loggedInUser.username);
  };

  ws.onmessage = async function (event) {
    try {
      const message = JSON.parse(event.data);

      const chatbox = document.getElementById("chatbox");
      const messageDisplay = document.getElementById("message-display");
      const messageElement = document.createElement("p");
      const details = document.createElement("div");
      const messageInput = document.getElementById("messageInput");
      const sendButton = document.getElementById("send-btn");
      const typingIndicator = document.getElementById("typing");

      if (message.type === "private_chat") {
        // Show notification only to the recipient
        if (loggedInUser.username === message.recipient) {
          const res = {
            success: true,
            message: `You have a new message from ${message.sender}`,
            data: message,
          };
          showToast(res);
        }
        fetchedMsg.push(message.message);
        loadedMessages.push(message.message);
        // Move sender and recipient to top of the chat-user-list only for them
        reorderChatUserList(message.sender, message.recipient);

        if (
          message.message.chat_id &&
          (activeChat[1] === message.recipient ||
            activeChat[1] === message.sender) &&
          chatbox.style.display === "block"
        ) {
          activeChat = [message.message.chat_id, message.recipient];

          await handlePrivateChat(message.sender, message.recipient);

          if (
            messageDisplay.innerHTML === "No chat history." ||
            messageDisplay.innerHTML === "Chat started but no messages yet."
          ) {
            messageDisplay.innerHTML = "";
          }

          // Append new message to messageDisplay
          details.classList.add("details");
          messageElement.classList.add("msg");
          timestamp = new Date(message.message.created_at).toLocaleString();
          details.textContent = `${message.sender} (${timestamp})`;
          messageElement.textContent = message.message.content;

          if (loggedInUser.username === message.sender) {
            messageElement.classList.add("from-me");
            details.classList.add("my-details");
          }

          messageDisplay.appendChild(messageElement);
          messageDisplay.appendChild(details);
          messageDisplay.scrollTop = messageDisplay.scrollHeight;

          messageInput.style.display = "block";
          sendButton.style.display = "block";
          // Listen to input events (typing) in the messageInput field
          messageInput.removeEventListener("input", (event) =>
            handleTyping(event, loggedInUser.username, user.username)
          ); // Remove previous event listener to avoid duplicates
          messageInput.addEventListener("input", (event) =>
            handleTyping(event, message.sender, message.recipient)
          );
        }
      } else if (message.type === "typing") {
        // Process incoming typing notifs and display typing indicator
        if (
          message.typing &&
          loggedInUser.username === message.recipient &&
          message.sender === activeChat[1]
        ) {
          // Show typing indicator if the recipient is the current chat recipient
          typingIndicator.style.display = "block";
          typingIndicator.innerHTML = `
          <span>${message.sender} is typing</span>
          <div class="dot-flashing"></div>
        `;
        } else if (!message.typing) {
          // Hide typing indicator when typing stops
          typingIndicator.style.display = "none";
          typingIndicator.innerHTML = ""; // Clear the content
        }
      } else if (message.type === "fetch_all_users") {
        onlineUsernames = new Set(
          message.users
            .filter((user) => user.isOnline)
            .map((user) => user.username)
        );
        const chatUsersList = document.getElementById("chat-users-list");
        const currentUsers = Array.from(chatUsersList.children).map((li) =>
          li.textContent.replace(" 👋", "")
        );

        // If chat-user-list is empty or a new user registers, reset the list ShowAllChatUsers
        if (
          currentUsers.length === 0 ||
          currentUsers.length + 1 !== message.users.length
        ) {
          await ShowAllChatUsers(message.users);
        } else {
          let chatUsersLi = chatUsersList.getElementsByTagName("li");

          // Update only the online/offline style without rearranging the list
          for (const li of chatUsersLi) {
            const username = li.textContent.replace(" 👋", "");
            const isOnline = onlineUsernames.has(username);

            if (isOnline) {
              if (!li.classList.contains("isOnline")) {
                li.classList.add("isOnline");
                li.textContent = `${username} 👋`;
              }
            } else {
              if (li.classList.contains("isOnline")) {
                li.classList.remove("isOnline");
                li.textContent = username;
              }
            }
          }
        }
      }
    } catch (error) {
      console.error("Failed to parse WebSocket message:", event.data, error);
    }
  };

  ws.onclose = function (event) {
    console.log(
      loggedInUser.username,
      "Offline: WebSocket connection closed on event:",
      event
    );
    stopTyping(loggedInUser.username, activeChat[1]);
  };

  ws.onerror = function (error) {
    console.error("Offline: WebSocket error:", error);
    stopTyping(loggedInUser.username, activeChat[1]);
  };
}

// Function to reorder chat-user-list (when new private_chat msg arrives)
function reorderChatUserList(sender, recipient) {
  const chatUsersList = document.getElementById("chat-users-list");
  const users = Array.from(chatUsersList.children);

  if (loggedInUser.username === sender) {
    const recipientElement = users.find((li) =>
      li.textContent.includes(recipient)
    );
    if (recipientElement) chatUsersList.removeChild(recipientElement);
    if (recipientElement && recipient !== sender)
      chatUsersList.prepend(recipientElement);
  } else if (loggedInUser.username === recipient) {
    const senderElement = users.find((li) => li.textContent.includes(sender));
    if (senderElement) chatUsersList.removeChild(senderElement);
    if (senderElement) chatUsersList.prepend(senderElement);
  }
}

// Handles messageInput and send
async function handlePrivateChat(senderUsername, recipientUsername) {
  document.getElementById("chatbox").style.display = "block";
  document.getElementById("messages").style.display = "block";
  const chatHeader = document.getElementById("chat-header");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");
  messageInput.style.display = "block";
  sendButton.style.display = "block";
  // Listen to input events (typing) in the messageInput field
  messageInput.removeEventListener("input", (event) =>
    handleTyping(event, senderUsername, recipientUsername)
  ); // Remove previous event listener to avoid duplicates
  messageInput.addEventListener("input", (event) =>
    handleTyping(event, senderUsername, recipientUsername)
  );

  if (loggedInUser.username === senderUsername) {
    chatHeader.textContent = `Chat with ${recipientUsername}`;
    messageInput.placeholder = `Type a message to ${recipientUsername}`;
    sendButton.onclick = null;
    sendButton.onclick = () => sendMessage(senderUsername, recipientUsername);
    stopTyping(senderUsername, recipientUsername);
  } else if (loggedInUser.username === recipientUsername) {
    chatHeader.textContent = `Chat with ${senderUsername}`;
    messageInput.placeholder = `Type a message to ${senderUsername}`;
    sendButton.onclick = null;
    sendButton.onclick = () => sendMessage(recipientUsername, senderUsername);
    stopTyping(recipientUsername, senderUsername);
  }
}

// Shows messages in the chatbox/messageDisplay w/throttle 10msgs at a time
let scrollHandler = null; // Reference to the current scroll handler

async function showMessages(chatID, recipientUsername) {
  if (activeChat[0] !== chatID || activeChat[1] !== recipientUsername) return; // Check if the chat is still open
  const chatHeader = document.getElementById("chat-header");
  const messageDisplay = document.getElementById("message-display");
  chatHeader.textContent = `Chat with ${recipientUsername}`;
  messageDisplay.innerHTML = ""; // Clear previous messages

  fetchedMsg = await fetchData(`/api/chat-messages/${chatID}`);
  if (!fetchedMsg || fetchedMsg.length === 0) {
    messageDisplay.innerHTML = "Chat started but no messages yet.";
    stopTyping(loggedInUser.username, recipientUsername);
    return;
  }

  // Throttle to process msgs in batches of 10
  const batchSize = 10;

  function renderMessages(batch, appendToTop = true) {
    batch.forEach((message) => {
      const messageElement = document.createElement("p");
      const details = document.createElement("div");
      details.classList.add("details");
      messageElement.classList.add("msg");
      timestamp = new Date(message.created_at).toLocaleString();
      details.textContent = `${message.created_by_username} (${timestamp})`;
      messageElement.textContent = message.content;

      if (loggedInUser.username === message.created_by_username) {
        messageElement.classList.add("from-me");
        details.classList.add("my-details");
      }

      if (appendToTop) {
        messageDisplay.prepend(details);
        messageDisplay.prepend(messageElement);
      } else {
        messageDisplay.append(messageElement);
        messageDisplay.append(details);
      }

      stopTyping(loggedInUser.username, recipientUsername);
    });
  }
  // Load initial batch (most recent messages)
  const initialBatch = fetchedMsg.slice(-batchSize);
  loadedMessages = [...initialBatch];
  renderMessages(initialBatch, false);

  // Scroll to the bottom of the chatbox to show the latest messages
  messageDisplay.scrollTop = messageDisplay.scrollHeight;

  // Add scroll event listener to load more messages when scrolling near the top
  scrollHandler = () => {
    if (activeChat[0] !== chatID || activeChat[1] !== recipientUsername) return; // Check if the chat is still open
    if (
      messageDisplay.scrollTop === 0 &&
      loadedMessages.length < fetchedMsg.length
    ) {
      const previousHeight = messageDisplay.scrollHeight; // Store the current height
      const remainingMessages = fetchedMsg.length - loadedMessages.length;
      const olderBatch = fetchedMsg
        .slice(Math.max(0, remainingMessages - batchSize), remainingMessages)
        .reverse(); // Get the next batch of older messages, in reverse order
      loadedMessages = [...olderBatch, ...loadedMessages];
      renderMessages(olderBatch);
      // Maintain the scroll position after loading more messages
      messageDisplay.scrollTop = messageDisplay.scrollHeight - previousHeight;
    }
  };
  messageDisplay.addEventListener("scroll", scrollHandler);
}

async function getChatID(senderUsername, recipientUsername) {
  const data = await fetchData(
    `/api/get-chat-id?sender=${senderUsername}&recipient=${recipientUsername}`
  );
  return data.data; // Return the chatID from the response
}

// Send message to server
function sendMessage(senderUsername, recipientUsername) {
  let input = document.getElementById("messageInput");
  let messageContent = input.value.trim();
  if (messageContent.length === 0) return;

  if (recipientUsername && ws && ws.readyState === WebSocket.OPEN) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          type: "private_chat",
          sender: senderUsername,
          recipient: recipientUsername,
          content: messageContent,
          //timestamp: new Date().toLocaleString(),
          timestamp: new Date().toISOString(), // Use ISO 8601 format
        })
      );
    }
  }
  input.value = "";
}

/* TYPING */

let typingTimeout;
let isTyping = false;

function handleTyping(event, senderUsername, recipientUsername) {
  if (!recipientUsername) return; // Ensure a chat is active

  if (!isTyping) {
    isTyping = true;
    sendTypingStatus(true, senderUsername, recipientUsername);
  }

  clearTimeout(typingTimeout);
  typingTimeout = setTimeout(() => {
    isTyping = false;
    sendTypingStatus(false, senderUsername, recipientUsername);
    stopTyping(senderUsername, recipientUsername);
  }, 2000); // 2 seconds of inactivity
}

// Send typing status to the server
function sendTypingStatus(isTyping, senderUsername, recipientUsername) {
  if (!senderUsername || !recipientUsername) return;
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(
      JSON.stringify({
        type: "typing",
        sender: senderUsername,
        recipient: recipientUsername,
        typing: isTyping,
      })
    );
  }
}

// Add a function to handle when typing stops or input loses focus
function stopTyping(senderUsername, recipientUsername) {
  if (isTyping) {
    isTyping = false;
    sendTypingStatus(false, senderUsername, recipientUsername);
  }
}
