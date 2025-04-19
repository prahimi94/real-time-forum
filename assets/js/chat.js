let ws;
//let onlineUsernames = [];
let allUserData = {};
let isProcessingMessages = false;
//let anotherUserClicked = false;
let chatboxOpen = false;
let openChat = [(ID = null), (recipient = null)];
let loadedMessages = []; // Keep track of already loaded messages

let typingTimeout;
let isTyping = false;

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

function fetchAllChatUsers() {
  fetch("/api/users", {
    method: "GET",
    headers: { "Content-Type": "application/json" },
  })
    .then((res) =>
      res
        .json()
        .catch(() => ({ success: false, message: "Invalid JSON response" }))
    ) // Prevent JSON parse errors
    .then(async (data) => {
      allUserData = data.data;
      //await ShowAllChatUsers(allUserData);
    });
  // do more things and hdle errors
}

async function ShowAllChatUsers(allUserData) {
  if (!Array.isArray(allUserData)) {
    allUserData = Object.values(allUserData);
  }
  try {
    const chatUsersList = document.getElementById("chat-users-list");
    chatUsersList.innerHTML = ""; // Clear the current list
    const messageDisplay = document.getElementById("message-display");
    const chatbox = document.getElementById("chatbox");
    const chatHeader = document.getElementById("chat-header");
    const messages = document.getElementById("messages");
    const messageInput = document.getElementById("messageInput");
    const sendButton = document.getElementById("send-btn");

    allUserData.forEach(async (user) => {
      if (user.username === loggedInUser.username) return;
      const chatID = await getChatID(
        loggedInUser.username,
        user.username
      );
      const li = document.createElement("li");
      //const isOnline = onlineUsernames.includes(user.username);


      if (user.isOnline) {
        li.classList.add("isOnline");
        li.textContent = `${user.username} (Online)`;
        openChat = [chatID, user.username];
        if (chatboxOpen || chatbox.style.display === "block") {
          chatboxOpen = true;
          messageInput.style.display = "block";
          sendButton.style.display = "block";

          // Listen to input events (typing) in the messageInput field
          /* messageInput.removeEventListener("input", (event) =>
            handleTyping(event, loggedInUser.username, user.username)
          ); */// Remove previous event listener to avoid duplicates
          messageInput.addEventListener("input", (event) =>
            handleTyping(event, loggedInUser.username, user.username)
          );
        }

      } else { // if user is offline
        li.textContent = user.username;
        isProcessingMessages = false;
        openChat = [null, null];
        chatboxOpen = false;
        chatbox.style.display = "none";
        messageInput.style.display = "none";
        sendButton.style.display = "none";
        document.getElementById("typing").style.display = "none";
        stopTyping(loggedInUser.username, user.username); // Stop typing when user is offline
      }

      // Manage event when username is clicked
      li.onclick = null; // Clear previous handler
      li.onclick = async (e) => {
        e.preventDefault();

        //isProcessingMessages = false;
        //anotherUserClicked = true;
        loadedMessages = []; // Reset the loaded messages array
        messageDisplay.innerHTML = "";

        if (chatboxOpen) {
          chatboxOpen = false;
          openChat = [null, null];
          chatbox.style.display = "none";
          stopTyping(loggedInUser.username, user.username); // stop typing when chatbox is closed
        } else if (!chatboxOpen) {
          chatboxOpen = true;
          chatbox.style.display = "block";
          messageDisplay.scrollTop = messageDisplay.scrollHeight;
          messageDisplay.innerHTML = "";
          messages.style.display = "block";
          messageDisplay.style.display = "block";

          if (user.isOnline) {
            if (chatID) {
              chatboxOpen = true;
              openChat = [chatID, user.username];
              messageDisplay.innerHTML = ""; // Clear previous messages
              await showMessages(chatID, user.username);
              messageInput.style.display = "block";
              sendButton.style.display = "block";

              // Listen to input events (typing) in the messageInput field
              /*     messageInput.removeEventListener("input", (event) =>
                handleTyping(event, loggedInUser.username, user.username)
              ); */// Remove previous event listener to avoid duplicates
              messageInput.addEventListener("input", (event) =>
                handleTyping(event, loggedInUser.username, user.username)
              );
              await handlePrivateChat(loggedInUser.username, user.username);

            } else if (chatID === 0 || !chatID) {
              openChat = [null, null];
              chatboxOpen = true;
              messageInput.style.display = "none";
              sendButton.style.display = "none";
              chatHeader.textContent = `Chat with ${user.username}`;
              messageDisplay.innerHTML = ""; // Clear previous messages
              messageDisplay.innerHTML = "No chat history.";
              stopTyping(loggedInUser.username, user.username);
            }

          } else { //if user is offline

            if (chatID) {
              chatboxOpen = true;
              openChat = [chatID, user.username];
              messageDisplay.innerHTML = ""; // Clear previous messages
              await showMessages(chatID, user.username);
              messageInput.style.display = "none";
              sendButton.style.display = "none";
              stopTyping(loggedInUser.username, user.username);
            } else if (chatID === 0 || !chatID) {
              chatboxOpen = true;
              openChat = [null, null];
              messageInput.style.display = "none";
              sendButton.style.display = "none";
              chatHeader.textContent = `Chat with ${user.username}`;
              messageDisplay.innerHTML = ""; // Clear previous messages
              messageDisplay.innerHTML = "No chat history.";
              stopTyping(loggedInUser.username, user.username);
            }
          }
        }
      };

      chatUsersList.appendChild(li);
    });
  } catch (error) {
    console.error("Error fetching all users:", error);
  }
}

function connect() {
  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = async function () {
    console.log("Online: Connected to WebSocket server");
    allUserData = {};
    document.getElementById("chat-users-list").innerHTML = ""; // Clear the current list
    fetchAllChatUsers(); // Fetch all users to populate the list
    await ShowAllChatUsers(allUserData); // Show all users in the chat list
  };

  ws.onmessage = async function (event) {
    try {
      /*  onlineUsernames = await fetchData("/api/online-users");
      allUserData = await fetchData("/api/users");
      await showAllChatUsers(); */

      const message = JSON.parse(event.data);

      const chatbox = document.getElementById("chatbox");
      const messageDisplay = document.getElementById("message-display");
      const messageElement = document.createElement("p");
      const details = document.createElement("div");
      const messageInput = document.getElementById("messageInput");
      const sendButton = document.getElementById("send-btn");
      const typingIndicator = document.getElementById("typing");


      if (message.type === "private_chat") {
        //await handlePrivateChat(message.sender, message.recipient);
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
        if (
          chatID &&
          openChat[0] === chatID &&
          (openChat[1] === message.recipient ||
            openChat[1] === message.sender) /*  &&
        message.users.includes(message.sender) &&
          message.users.includes(message.recipient) */
        ) {
          chatboxOpen = true;
          openChat = [chatID, message.recipient];
          /* await showMessages(chatID, message.recipient); */

          details.classList.add("details");
          messageElement.classList.add("msg");
          timestamp = new Date(message.message.created_at).toLocaleString(); // Format the timestamp
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
          typingIndicator.style.display = "none";
          stopTyping(message.sender, message.recipient);
          // Listen to input events (typing) in the messageInput field
          /*     messageInput.removeEventListener("input", (event) =>
            handleTyping(event, loggedInUser.username, user.username)
          );  */// Remove previous event listener to avoid duplicates
          messageInput.addEventListener("input", (event) =>
            handleTyping(event, message.sender, message.recipient)
          );
        } else {
          chatboxOpen = false;
          openChat = [null, null];
          chatbox.style.display = "none";
          messageInput.style.display = "none";
          sendButton.style.display = "none";
          typingIndicator.style.display = "none";
          stopTyping(message.sender, message.recipient);
        }
      } /* else if (message.type === "show_all_users") {
        await ShowAllChatUsers(message.users);
      } */ else if (message.type === "typing") {
        // Process incoming typing notifications and display the typing indicator in the chat UI

        if (message.typing && message.sender === openChat[1]) {
          // Show typing indicator if the sender is the current chat recipient
          typingIndicator.textContent = `${message.sender} is typing...`;
          typingIndicator.style.display = "block";
        } else if (!message.typing && message.sender === openChat[1]) {
          // Hide typing indicator when typing stops
          typingIndicator.style.display = "none";
        }
      } else if (message.type === "fetch_all_users") {
        //fetchAllChatUsers();
        document.getElementById("chat-users-list").innerHTML = ""; // Clear the current list
        await ShowAllChatUsers(message.users);
      }
    } catch (error) {
      console.error("Failed to parse WebSocket message:", event.data, error);
    }
  };

  ws.onclose = function (event) {
    console.log("Offline: WebSocket connection closed on event:", event);
    stopTyping(loggedInUser.username, openChat[1]);
  };

  ws.onerror = function (error) {
    console.error("Offline: WebSocket error:", error);
    stopTyping(loggedInUser.username, user.username);
  };
}

async function handlePrivateChat(senderUsername, recipientUsername) {
  document.getElementById("chatbox").style.display = "block";
  document.getElementById("messages").style.display = "block";
  const chatHeader = document.getElementById("chat-header");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");

  messageInput.style.display = "block";
  sendButton.style.display = "block";
  // Listen to input events (typing) in the messageInput field
  /*  messageInput.removeEventListener("input", (event) =>
    handleTyping(event, loggedInUser.username, user.username)
  ); */ // Remove previous event listener to avoid duplicates
  messageInput.addEventListener("input", (event) =>
    handleTyping(event, senderUsername, recipientUsername)
  );

  if (loggedInUser.username === senderUsername) {
    chatHeader.textContent = `Chat with ${recipientUsername}`;
    messageInput.placeholder = `Type a message to ${recipientUsername}`;
    sendButton.onclick = null; // Clear previous handler
    sendButton.onclick = () => sendMessage(senderUsername, recipientUsername);
    document.getElementById("typing").style.display = "none";
    stopTyping(senderUsername, recipientUsername);
  } else if (loggedInUser.username === recipientUsername) {
    chatHeader.textContent = `Chat with ${senderUsername}`;
    messageInput.placeholder = `Type a message to ${senderUsername}`;
    sendButton.onclick = null; // Clear previous handler
    sendButton.onclick = () => sendMessage(recipientUsername, senderUsername);
    document.getElementById("typing").style.display = "none";
    stopTyping(recipientUsername, senderUsername);
  }

  // Send a message to the server to initiate or check the chat
  /*   if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          type: "private_chat",
          sender: senderUsername,
          recipient: recipientUsername,
        })
      );
    } */
}
let scrollHandler = null; // Reference to the current scroll handler

async function showMessages(chatID, recipientUsername) {
  if (openChat[0] !== chatID || openChat[1] !== recipientUsername) return; // Check if the chat is still open
  /*  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn"); */
  const chatHeader = document.getElementById("chat-header");
  const messageDisplay = document.getElementById("message-display");
  chatHeader.textContent = `Chat with ${recipientUsername}`;
  messageDisplay.innerHTML = ""; // Clear previous messages
  // Remove any previous scroll listener
  if (scrollHandler) {
    messageDisplay.removeEventListener("scroll", scrollHandler);
  }
  let messages = await fetchData(`/api/chat-messages/${chatID}`);

  if (!messages || messages.length === 0) {
    messageDisplay.textContent = "Chat started but no messages yet.";
    stopTyping(loggedInUser.username, recipientUsername);
    return;
  }
  // Reverse the messages to show the latest ones at the bottom
  // messages.reverse();

  //isProcessingMessages = false;
  //anotherUserClicked = false;

  // Throttle to process msgs in batches of 10
  const batchSize = 10;

  function renderMessages(batch, appendToTop = true) {
    //if (isProcessingMessages/*  || anotherUserClicked */) return; // Prevent multiple triggers
    //isProcessingMessages = true;

    // const batch = messages.splice(0, batchSize); // Get the next batch of messages
    //loadedMessages = [...batch, ...loadedMessages]; // Add to the loaded messages

    batch.forEach((message) => {
      //if (anotherUserClicked) return;

      const messageElement = document.createElement("p");
      const details = document.createElement("div");
      details.classList.add("details");
      messageElement.classList.add("msg");
      timestamp = new Date(message.created_at).toLocaleString(); // Format the timestamp
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


      document.getElementById("typing").style.display = "none";
      stopTyping(loggedInUser.username, recipientUsername);
    });

    // isProcessingMessages = false;

  }

  // Load initial batch (most recent messages)
  //messages.reverse(); // Ensure the latest messages are at the end
  const initialBatch = messages.slice(-batchSize);
  loadedMessages = [...initialBatch];
  renderMessages(initialBatch, false);

  // Scroll to the bottom of the chatbox to show the latest messages
  messageDisplay.scrollTop = messageDisplay.scrollHeight;

  // Add scroll event listener to load more messages when scrolling near the top
  scrollHandler = () => {
    if (openChat[0] !== chatID || openChat[1] !== recipientUsername) return; // Check if the chat is still open
    if (messageDisplay.scrollTop === 0 && loadedMessages.length < messages.length) {
      const previousHeight = messageDisplay.scrollHeight; // Store the current height
      const remainingMessages = messages.length - loadedMessages.length;
      const olderBatch = messages.slice(
        Math.max(0, remainingMessages - batchSize),
        remainingMessages
      ).reverse(); // Get the next batch of older messages
      loadedMessages = [...olderBatch, ...loadedMessages];
      renderMessages(olderBatch);
      // Maintain the scroll position after loading more messages
      messageDisplay.scrollTop = messageDisplay.scrollHeight - previousHeight;
    }
  }
  messageDisplay.addEventListener("scroll", scrollHandler);


  /*   if (!message.users.includes(recipientUsername)) {
    messageInput.style.display = "none";
    sendButton.style.display = "none";
  } */
}

async function getChatID(senderUsername, recipientUsername) {
  const data = await fetchData(
    `/api/get-chat-id?sender=${senderUsername}&recipient=${recipientUsername}`
  );
  return data.data; // Return the chatID from the response
}

function sendMessage(senderUsername, recipientUsername) {
  isProcessingMessages = false;
  let input = document.getElementById("messageInput");
  let message = input.value.trim(); // Remove leading/trailing whitespace
  if (message.length === 0) return; // Do not accept empty messages

  if (recipientUsername && ws && ws.readyState === WebSocket.OPEN) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          type: "private_chat",
          sender: senderUsername,
          recipient: recipientUsername,
          content: message,
          //timestamp: new Date().toLocaleString(),
          timestamp: new Date().toISOString(), // Use ISO 8601 format
        })
      );
    }
  }
  input.value = "";
}

/* TYPING */

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

function sendTypingStatus(isTyping, senderUsername, recipientUsername) {
  if (!senderUsername || !recipientUsername) return; // Ensure sender and recipient are provided
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