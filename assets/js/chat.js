let ws;
let onlineUsernames = [];
let allUserData = {};
let isProcessingMessages = false;
let anotherUserClicked = false;
let chatboxOpen = false;

async function fetchOnlineUsers() {
  try {
    const response = await fetch("/api/online-users");
    if (!response.ok) {
      console.error("Failed to fetch online users");
      return;
    }

    onlineUsernames = await response.json();
  } catch (error) {
    console.error("Error fetching online users:", error);
  }
}

async function fetchAllChatUsers() {
  try {
    const response = await fetch("/api/users");
    if (!response.ok) {
      console.error("Failed to fetch all users");
      return;
    }

    allUserData = await response.json();
    const chatUsersList = document.getElementById("chat-users-list");
    chatUsersList.textContent = ""; // Clear the current list

    allUserData.forEach((user) => {
      if (user.username !== loggedInUser.username) {
        const li = document.createElement("li");

        if (onlineUsernames.includes(user.username)) {
          const link = document.createElement("a");
          link.href = `#chat-with-${user.username}`;
          link.onclick = (e) => {
            isProcessingMessages = false;
            anotherUserClicked = true;
            document.getElementById("message-display").innerHTML = "";
            e.preventDefault(); // Prevent default link behavior
            if (chatboxOpen) {
              chatboxOpen = false;
              document.getElementById("chatbox").style.display = "none";
            } else if (!chatboxOpen) {
              chatboxOpen = true;
              document.getElementById("chatbox").style.display = "block";
              handlePrivateChat(loggedInUser.username, user.username);
            }
          };
          link.textContent = `${user.username} (Online)`;
          li.appendChild(link);
        } else {
          li.textContent = user.username;
          li.onclick = async (e) => {
            isProcessingMessages = false;
            anotherUserClicked = true;
            document.getElementById("message-display").innerHTML = "";

            e.preventDefault
            if (chatboxOpen) {
              chatboxOpen = false;
              document.getElementById("chatbox").style.display = "none";
              document.getElementById("message-display").innerHTML = "";
            } else if (!chatboxOpen) {
              chatboxOpen = true;
              document.getElementById("chatbox").style.display = "block";

              document.getElementById("messages").style.display = "block";
              document.getElementById("message-display").style.display = "block";
              document.getElementById("message-display").innerHTML = "";
              // Fetch the chat ID and load the chat messages
              const chatID = await getChatIDForUsers(loggedInUser.username, user.username);
              if (chatID) {

                await fetchChatMessages(chatID, user.username);

              } else if (chatID === 0 || !chatID) {
                const chatHeader = document.getElementById("chat-header");
                document.getElementById("messageInput").style.display = "none";
                document.getElementById("send-btn").style.display = "none";

                chatHeader.textContent = `Chat with ${user.username}`;
                document.getElementById("message-display").innerHTML = "No chat history.";
              } else {
                console.error("Failed to retrieve chat ID");
              }
            }
          }
        }

        chatUsersList.appendChild(li);
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
    fetchAllChatUsers(); // Fetch all users to populate the list
  };

  ws.onmessage = function (event) {
    try {
      const message = JSON.parse(event.data);

      if (message.type === "message_content") {
        handlePrivateChat(message.sender, message.recipient);
        if (message.recipient === loggedInUser.username) {
          const res = {
            success: true,
            message: `You have a new message from ${message.sender}`,
            data: message,
          };
          showToast(res);
        }
        chatboxOpen = true;
      }

      fetchOnlineUsers();
      fetchAllChatUsers();
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
  document.getElementById("chatbox").style.display = "block";
  document.getElementById("messages").style.display = "block";
  const chatHeader = document.getElementById("chat-header");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");

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

  // Fetch the chat ID and load the chat messages
  const chatID = await getChatIDForUsers(senderUsername, recipientUsername);
  if (chatID) {
    if (loggedInUser.username === senderUsername
    ) {
      await fetchChatMessages(chatID, recipientUsername);

      messageInput.style.display = "block";
      sendButton.style.display = "block";
    } else if (
      loggedInUser.username === recipientUsername
    ) {
      await fetchChatMessages(chatID, senderUsername);

      messageInput.style.display = "block";
      sendButton.style.display = "block";
    }
  } else {
    console.error("Failed to retrieve chat ID");
  }
}

async function fetchChatMessages(chatID, recipientUsername) {
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("send-btn");
  const chatHeader = document.getElementById("chat-header");
  const messageDisplay = document.getElementById("message-display");

  messageDisplay.innerHTML = ""
  chatHeader.textContent = `Chat with ${recipientUsername}`;

  messageInput.style.display = "none";
  sendButton.style.display = "none";
  try {
    const response = await fetch(`/api/chat-messages/${chatID}`);
    if (!response.ok) {
      console.error("Failed to fetch chat messages");
      return;
    }

    let messages = await response.json();


    if (messages && messages.length != 0) {


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
          const parsedContent = JSON.parse(message.content);
          const messageElement = document.createElement("p");
          const details = document.createElement("div");
          details.classList.add("details");
          messageElement.classList.add("msg");
          details.textContent = `${parsedContent.sender} (${parsedContent.timestamp})`;
          messageElement.textContent = parsedContent.content;
          messageDisplay.prepend(details);
          messageDisplay.prepend(messageElement);
          if (loggedInUser.username === parsedContent.sender) {
            messageElement.classList.add("from-me");
            details.classList.add("my-details");
          }
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
    } else {
      messageDisplay.textContent = "No messages.";
    }
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
