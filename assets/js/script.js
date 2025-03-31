let categories = [];
let posts = [];

// Function to check if the session is active
async function checkSessionActive() {
    try {
        // Call the backend API to check session status
        const response = await fetch('/api/check-session', {
            method: 'GET',
            credentials: 'include', // Include cookies in the request
        });

        if (response.ok) {
            const data = await response.json();
            if (data.active) {
                return true; // Session is active
            } else {
                return false; // Session is not active
            }
        } else {
            return false; // Session is not active
        }
    } catch (error) {
        return false; // Session is not active
    }
}

function showNotAuthenticatedContainer() {
    const authContainer = document.getElementsByTagName('main')[0];
    authContainer.innerHTML = `
    <header>
        <nav class="navbar navbar-expand-lg navbar-light py-10" aria-label="Thirteenth navbar example">
            <div class="container-fluid">
                <a class="navbar-brand" href="/">
                    <img src="/img/logo.png" alt="Forum" class="me-2" style="height: 60px;">
                </a>
                <button class="navbar-toggler" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasNavbar" aria-controls="offcanvasNavbar" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="offcanvas offcanvas-end" tabindex="-1" id="offcanvasNavbar" aria-labelledby="offcanvasNavbarLabel">
                    <div class="offcanvas-header">
                        <h5 class="offcanvas-title" id="offcanvasNavbarLabel"><img src="/img/logo.png" alt="Forum" class="me-2" style="height: 50px;"></h5>
                        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
                    </div>
                    <div class="offcanvas-body">
                        <ul class="navbar-nav justify-content-end flex-grow-1 pe-3">
                            <li class="nav-item">
                                <a class="btn btn-success singIn-btn" href="/auth/"><i class="fa-solid fa-user pe-2"></i> Sign in</a>
                            </li>
                        </ul>
                        
                    </div>
                </div>

            <div class="collapse navbar-collapse" id="navbarsExample11">
                <div class="navbar-nav col-sm-12 col-lg-10 justify-content-lg-center">
                    
                </div>
                <div class="d-lg-flex col-sm-12 col-lg-2 justify-content-lg-end">
                    <a class="btn btn-success singIn-btn" href="/auth/"><i class="fa-solid fa-user pe-2"></i> Sign in</a>
                </div>
            </div>
            </div>
        </nav>
    </header>
    <header class="hero">
        <div class="container">
            <div class="row align-items-center">
                <div class="col-md-6 text-content">
                    <h1 class="text-header">Welcome To <br/> Future Lab!</h1>
                    <p style="color: #7c7e82;">Discover a thriving online community where ideas come to life! Our forum is the perfect place to connect with like-minded individuals, share your thoughts, ask questions, and explore engaging discussions across a variety of topics. Join us today and be part of something bigger!</p>
                    <!-- <a class="btn btn-success" href="/auth/">Join Now</a> -->
                </div>
                <div class="col-md-6">
                    <!-- <img src="/img/background.png" class="img-fluid" alt="Background Image"> -->
                    <!-- https://lottie.host/e4c9739c-85ef-441d-8064-2b6fc1c7ad60/rrBxdg5bgH.lottie -->
                    <dotlottie-player
                    src="/img/Animation-1739026245825.lottie"
                    background="transparent"
                    speed="1"
                    style="width: 100%;"
                    loop
                    autoplay
                    ></dotlottie-player>
                </div>
            </div>
        </div>
    </header>
    `;
}

function showAuthContainer() {
    const authContainer = document.getElementsByTagName('main')[0];
    authContainer.innerHTML = `
    <header>
        <nav class="navbar navbar-expand-lg navbar-light py-10" aria-label="Thirteenth navbar example">
            <div class="container-fluid">
                <a class="navbar-brand" href="/">
                    <img src="/img/logo.png" alt="Forum" class="me-2" style="height: 60px;">
                </a>
                <button class="navbar-toggler" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasNavbar" aria-controls="offcanvasNavbar" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="offcanvas offcanvas-end" tabindex="-1" id="offcanvasNavbar" aria-labelledby="offcanvasNavbarLabel">
                    <div class="offcanvas-header">
                        <h5 class="offcanvas-title" id="offcanvasNavbarLabel"><img src="/img/logo.png" alt="Forum" class="me-2" style="height: 50px;"></h5>
                        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
                    </div>
                    <div class="offcanvas-body">
                        <ul class="navbar-nav justify-content-end flex-grow-1 pe-3">
                            <li class="nav-item">
                                <a class="btn btn-success singIn-btn" href="/auth/"><i class="fa-solid fa-user pe-2"></i> Sign in</a>
                            </li>
                        </ul>
                        
                    </div>
                </div>

            <div class="collapse navbar-collapse" id="navbarsExample11">
                <div class="navbar-nav col-sm-12 col-lg-10 justify-content-lg-center">
                    
                </div>
                <div class="d-lg-flex col-sm-12 col-lg-2 justify-content-lg-end">
                    <a class="btn btn-success singIn-btn" href="/auth/"><i class="fa-solid fa-user pe-2"></i> Sign in</a>
                </div>
            </div>
            </div>
        </nav>
    </header>
    <header class="hero">
        <div class="container">
            <div class="row align-items-center">
                <div class="col-md-6 text-content">
                    <h1 class="text-header">Welcome To <br/> Future Lab!</h1>
                    <p style="color: #7c7e82;">Discover a thriving online community where ideas come to life! Our forum is the perfect place to connect with like-minded individuals, share your thoughts, ask questions, and explore engaging discussions across a variety of topics. Join us today and be part of something bigger!</p>
                    <!-- <a class="btn btn-success" href="/auth/">Join Now</a> -->
                </div>
                <div class="col-md-6">
                    <!-- <img src="/img/background.png" class="img-fluid" alt="Background Image"> -->
                    <!-- https://lottie.host/e4c9739c-85ef-441d-8064-2b6fc1c7ad60/rrBxdg5bgH.lottie -->
                    <dotlottie-player
                    src="/img/Animation-1739026245825.lottie"
                    background="transparent"
                    speed="1"
                    style="width: 100%;"
                    loop
                    autoplay
                    ></dotlottie-player>
                </div>
            </div>
        </div>
    </header>
    `;
}

function showAuthenticatedContainer() {
    const forumContainer = document.getElementsByTagName('main')[0];
    forumContainer.innerHTML = `
    <header>
            <nav class="navbar navbar-expand-lg navbar-light py-10" aria-label="Thirteenth navbar example">
                <div class="container-fluid">
                    <a class="navbar-brand" href="/">
                        <img src="/img/logo.png" alt="Forum" class="me-2" style="height: 60px;">
                    </a>
                <button class="navbar-toggler" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasNavbar" aria-controls="offcanvasNavbar" aria-label="Toggle navigation">
                        <span class="navbar-toggler-icon"></span>
                    </button>
                <div class="offcanvas offcanvas-end" tabindex="-1" id="offcanvasNavbar" aria-labelledby="offcanvasNavbarLabel">
                        <div class="offcanvas-header">
                        <h5 class="offcanvas-title" id="offcanvasNavbarLabel"><img src="/img/logo.png" alt="Forum" class="me-2" style="height: 50px;"></h5>
                        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
                        </div>
                        <div class="offcanvas-body">
                            <ul class="navbar-nav justify-content-end flex-grow-1 pe-3">
                                <li class="nav-item">
                                    <div class="me-3 text-center">
                                        <!-- todo bring profilePhoto -->
                                        <!-- { if .ProfilePhoto}}
                                    <img src="/uploads/{ .ProfilePhoto}}" alt="user image" class="rounded w-50 shadow">
                                    { else}} -->
                                    <i class="fa-solid fa-user rounded shadow" style="font-size: 4rem;padding: 1rem;"></i>
                                        <!-- { end}} -->
                                    </div>
                                </li>
                                <li class="nav-item text-center pt-3">
                                    <!-- todo bring name -->
                                    <!-- <span class="me-3">Welcome, { .Name }}</span> -->
                                    <span class="me-3">Welcome</span>
                                </li>
                                <li class="nav-item text-center pb-2 pt-2">
                                    <!-- todo bring email -->
                                    <!-- <span class="me-3">{ .Email }}</span> -->
                                    <span class="me-3">Email</span>
                                </li>
                                <li class="nav-item text-center pb-3">
                                <a class="btn btn-outline-secondary" href="/profile"><i class="fa-regular fa-address-card"></i></a>
                                    <a class="btn btn-danger" href="/logout/"><i class="fas fa-power-off"></i></a>
                                </li>
                                <li><div class="divaider mb-2"></div></li>
                                    <li class="nav-item">
                                        <form class="d-flex mt-3" role="search" method="get" action="/filterPosts">
                                        <input class="form-control me-2" type="search" placeholder="Search" aria-label="Search" name="post_info">
                                        <button class="btn btn-outline-success" type="submit"><i class="fa-solid fa-magnifying-glass"></i></button>
                                        </form>
                                    </li>
                                <li><div class="divaider mb-2 mt-2"></div></li>
                                <li class="nav-item">
                                    <a class="nav-link" href="/">Home</a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="/newPost/">Create Post</a>
                                    <!-- <a class="nav-link" href="/newPost/"><i class="fa-solid fa-plus pe-2"></i> Create Post</a> -->
                                </li>
                                <li class="nav-item dropdown">
                                <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">

                                My posts
                                </a>
                                <ul class="dropdown-menu">
                                    <li><a class="dropdown-item" href="javascript:fetchMyCreatedPosts()">My created posts</a></li>
                                    <li><a class="dropdown-item" href="javascript:fetchMyLikedPosts()">My liked posts</a></li>
                                </ul>
                            </li>
                            </ul>

                        </div>
                    </div>

                    <div class="collapse navbar-collapse" id="navbarsExample11">
                        <div class="navbar-nav col-sm-12 col-lg-10 justify-content-lg-center">
                            <div class="input-p-main">
                                <form class="d-flex" role="search" method="get" action="/filterPosts">
                            <input class="form-control me-2 input-p" type="search" placeholder="Search" aria-label="Search" name="post_info">
                            <button class="btn btn-outline-success" type="submit"><i class="fa-solid fa-magnifying-glass"></i></button>
                                </form>
                            </div>
                        </div>
                        <div class="d-lg-flex col-sm-12 col-lg-2 justify-content-lg-end">
                            <!-- todo check signIn -->
                            <div class="navbar-brand d-flex align-items-center navbar-logined" style="display: grid;grid-auto-flow: column;grid-column-gap: 10px">
                                <a class="nav-link link-un me-3" href="/">Home</a>

                                <div class="btn-group">
                                    <button type="button" class="btn btn-secondary dropdown-toggle btn-dropdown" data-bs-toggle="dropdown" aria-expanded="false">

                                        My posts
                                    </button>
                                    <ul class="dropdown-menu dropdown-menu-end">
                                        <li><a class="dropdown-item" href="javascript:fetchMyCreatedPosts()">My created posts</a></li>
                                        <li><a class="dropdown-item" href="javascript:fetchMyLikedPosts()">My liked posts</a></li>
                                    </ul>
                                </div>
                                <a class="btn btn-success createPost-btn" href="/newPost/"><i class="fa-solid fa-plus pe-2"></i> Create Post</a>

                                <div class="btn-group">
                                    <button type="button" class="btn btn-secondary dropdown-toggle dropdown-toggle-user rounded-circle" style="width: 43px; height: 43px;padding: 0;overflow: hidden;" data-bs-toggle="dropdown" aria-expanded="false">
                                        <div class="rounded-circle">
                                            <!-- todo bring .ProfilePhoto -->
                                            <!-- { if .ProfilePhoto}}
                                             <img src="/uploads/{ .ProfilePhoto}}" alt="user image" style="width: 43px; height: 43px; cursor: pointer;" >
                                            { else}} -->
                                            <i class="fa-solid fa-user"></i>
                                            <!-- { end}} -->
                                        </div>
                                    </button>
                                    <ul class="dropdown-menu dropdown-menu-end">
                                        <li><a class="dropdown-item" href="/profile"><i class="fa-regular fa-address-card me-2"></i> Profile</a></li>
                                        <li><a class="dropdown-item" href="/logout/"><i class="fas fa-power-off me-2"></i> Log Out</a></li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </nav>
        </header>


        <div class="container">
            <div class="row">
                <div class="col-md-3">
                    <div class="col-left-userLogin">
                        <ul class="left-menu">
                            <li><a class="left-item" href="/"><i class="fa-solid fa-house"></i> Home</a></li>
                            <li><a class="left-item" href="javascript:fetchMyCreatedPosts()"><i
                                        class="fa-solid fa-bullhorn"></i> My created posts</a></li>
                            <li><a class="left-item" href="javascript:fetchMyLikedPosts()"><i
                                        class="fa-solid fa-heart"></i> My liked posts</a></li>
                        </ul>

                        <p class="section-title-list">Categories</p>
                        <!-- this part is filled in script.js by fetchCategories -->
                        <ul class="left-menu" id="categoriesContainer"></ul>
                    </div>

                </div>
                <div class="col-md-6">
                    <!-- {template "posts" .}} -->
                    <!-- this part is filled in script.js by fetchPosts -->
                    <div class="row" id="newPostContainer">
                        <div class="col-sm-12 col-md-12 mb-3">
                            <div class="post-card">
                                <h4 class="text-center mb-4">New Post</h4>
                                <form enctype="multipart/form-data" action="/submitPost" method="post">
                                    <div class="mb-3">
                                        <div class="mt-3">
                                            <select id="categories" name="categories" required
                                                class="form-control multiSelect" multiple="multiple"
                                                data-placeholder="Select categories">
                                            </select>
                                        </div>
                                    </div>
                                    <div class="mb-3">
                                        <input type="text" class="form-control" placeholder="Title" required
                                            name="title">
                                    </div>
                                    <div class="mb-3">
                                        <textarea class="form-control" placeholder="Description" required rows="4"
                                            name="description"></textarea>
                                    </div>
                                    <div class="mb-3 text-center border p-3 rounded">
                                        <p class="text-muted">Attach an image or video (optional)</p>
                                        <input type="file" class="form-control" name="postFiles" multiple>
                                    </div>
                                    <button type="submit" class="btn btn-success w-100">Post</button>
                                </form>
                            </div>
                        </div>
                    </div>
                    <div class="row" id="postsContainer"></div>
                </div>
                <div class="col-md-3">
                    <div class="col-right-userLogin">
                        <div class="text-center w-100 py-4">
                        <!-- todo bring profilePhoto -->
                        <!-- { if .LoginUser.ProfilePhoto}}
                            <img src="/uploads/ { .LoginUser.ProfilePhoto}}" alt="user image" class="rounded w-50 shadow">
                        { else}} -->
                            <i class="fa-solid fa-user rounded shadow" style="font-size: 4rem;padding: 1rem;"></i>
                            <!-- { end}} -->
                        </div>
                        <!-- todo bring LoginUser.Name -->
                        <!-- <div class="text-center title-username"> { .LoginUser.Name}}</div> -->
                        <div class="text-center title-username">.LoginUser.Name</div>
                        <div class="info-box-username">
                            <!-- todo bring .LoginUser.Email -->
                            <!-- <p> { .LoginUser.Email}}</p> -->
                            <p>.LoginUser.Email</p>
                            <div class="info-box-logout">
                                <p style="margin-bottom: 1rem;"><a href="/profile"><i
                                            class="fa-regular fa-address-card me-2"></i> Profile</a></p>
                                <a href="/logout/"><i class="fas fa-power-off me-2" style="color: #c44343;"></i> Log
                                    Out</a>
                            </div>
                        </div>

                        <!-- ONLINE USERS -->
                        <div id="online-users" class="online-users-container">
                            <h2>Online Users</h2>
                            <ul id="online-users-list"></ul>
                        </div>
                        <!-- CHATBOX -->
                        <div class="chatbox" style="display: none;">
                            <h1>CHATBOX</h1>
                            <div id="messages"></div>
                            <div>
                                <input type="text" id="messageInput" placeholder="Enter your message">
                                <button onclick="sendMessage()">Send</button>
                            </div>
                        </div>
                        <!-- END CHATBOX -->
                    </div>
                </div>
            </div>
        </div>
    `;
}

async function fetchCategories() {
    const response = await fetch('/api/categories/');
    categories = await response.json();

    // load categories in sidebar
    const categoriesContainer = document.getElementById('categoriesContainer');
    categoriesContainer.innerHTML = "";
    categories.forEach(category => {
        const categoryITag = document.createElement('i');
        categoryITag.setAttribute('class', category.icon)
        categoryITag.style.color = category.color;
        categoryITag.style.padding = '0.75rem 1rem';

        const categoryATag = document.createElement('a');
        categoryATag.setAttribute('class', 'left-item');
        categoryATag.href = `/posts/${category.name}`;
        categoryATag.appendChild(categoryITag);
        categoryATag.innerHTML += category.name;

        const categoryLiTag = document.createElement('li');
        categoryLiTag.appendChild(categoryATag);

        categoriesContainer.appendChild(categoryLiTag);
    });

    // laod categories in select for new post
    const selectCategoriesContainer = document.getElementById('categories');
    const categoryOptions = categories.map(category =>
        `<option value="${category.id}">${category.name}</option>`
    ).join('');
    selectCategoriesContainer.innerHTML = categoryOptions;
}

async function fetchPosts() {
    const response = await fetch('/api/posts/');
    posts = await response.json();

    fillPostsInHtml(posts);
}

async function fetchMyCreatedPosts() {
    const response = await fetch('/api/myCreatedPosts/');
    posts = await response.json();

    fillPostsInHtml(posts.Posts);
}

async function fetchMyLikedPosts() {
    const response = await fetch('/api/myLikedPosts/');
    posts = await response.json();

    fillPostsInHtml(posts.Posts);
}

function fillPostsInHtml(posts) {
    // load posts for home page
    const postsContainer = document.getElementById('postsContainer');
    postsContainer.innerHTML = "";

    if (posts === null || posts.length === 0) {
        postsContainer.innerHTML = '<div class="col-md-12 text-center">No posts found!</div>';
        return;
    }
    posts.forEach(post => {
        const postImage = post.user.profile_photo
            ? `<img class="bd-placeholder-img flex-shrink-0 me-2 rounded" role="img" src="/uploads/${post.user.profile_photo}" width="45" height="45"/>`
            : `<div style="padding: 7px;"><i class="fa-solid fa-user" style="font-size: 2rem;"></i></div>`;

        const postCategories = post.categories.map(category =>
            `<span class="badge-p text-dark"><a href="/posts/${category.name}">${category.name}</a></span>`
        ).join('');

        const postFiles = post.post_files.map(post_file =>
            `<div class="col-md-12">
                <img src="/uploads/${post_file.file_uploaded_name}" alt="post image" class="rounded mb-1" style="width: 100%; max-height: 400px;">
            </div>`
        ).join('');

        const formattedDateTime = post.created_at.replace('T', ' ').replace('Z', '');

        const postHTML = `
            <div class="col-sm-12 col-md-12 mb-3">
                <div class="post-card">
                    <h5 class="mt-2 post-title">
                        <a href="/post/${post.uuid}" class="nav-link">${post.title}</a>
                    </h5>
                    <div>
                        <div class="d-flex text-body-secondary pt-3 m-posts">
                            ${postImage}
                            <div class="pb-3 mb-0 small lh-sm w-100 mb-3 ms-2 mt-1">
                                <div class="d-flex justify-content-between mb-1 m-posts-userInfo">
                                    <span class="post-user">${post.user.username}</span>
                                    <span class="text-right m-posts-ctg">${postCategories}</span>
                                </div>
                                <span class="d-block post-dateTime">${formattedDateTime}</span>
                            </div>
                        </div>
                    </div>
                    <p class="post-description">${post.description}</p>
                    ${postFiles}
                </div>
            </div>
        `;

        postsContainer.insertAdjacentHTML('beforeend', postHTML);
    });
}

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
            const link = document.createElement("a");
            link.href = "javascript:void(0);"; // Prevent default navigation
            link.textContent = username;
            link.addEventListener("click", () => openChatWithUser(username)); // Add click event
            li.appendChild(link);
            onlineUsersList.appendChild(li);
        });
    } catch (error) {
        console.error("Error fetching online users:", error);
    }
}

/* WEBSOCKET FOR CHAT */
let ws;
let activeChatUser = null; // Track the currently active chat user

function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = function () {
        console.log("Connected to WebSocket server");
    };

    ws.onmessage = function (event) {
        let messageDisplay = document.getElementById("messages");
        let message = JSON.parse(event.data);

        // Check if the message is for the active chat user
        if (message.sender === activeChatUser || message.receiver === activeChatUser) {
        // Append the message to the chatbox
        let messageElement = document.createElement("p");
        messageElement.textContent = `${message.sender}: ${message.text}`;
        messageDisplay.appendChild(messageElement);

        // Scroll to the bottom of the chatbox
        messageDisplay.scrollTop = messageDisplay.scrollHeight;
        }
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
    const input = document.getElementById("messageInput");
    const message = input.value;

    if (!activeChatUser) {
        alert("Please select a user to chat with.");
        return;
    }

    // Send the message as a JSON object
    const messagePayload = {
        receiver: activeChatUser,
        text: message,
    };

    ws.send(JSON.stringify(messagePayload));
    input.value = "";

    // Optionally, display the sent message in the chatbox
    const messageDisplay = document.getElementById("messages");
    const messageElement = document.createElement("p");
    messageElement.textContent = `You: ${message}`;
    messageDisplay.appendChild(messageElement);

    // Scroll to the bottom of the chatbox
    messageDisplay.scrollTop = messageDisplay.scrollHeight;
}
function openChatWithUser(username) {
    activeChatUser = username; // Set the active chat user
    const chatbox = document.querySelector(".chatbox");
    const messages = document.getElementById("messages");

    // Show the chatbox
    chatbox.style.display = "block";

    // Clear previous messages and set a header for the chat
    messages.innerHTML = `<h3>Chat with ${username}</h3>`;

    // Load previous messages with this user from the server
    fetch(`/api/chat/${username}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to fetch chat messages: ${response.status} ${response.statusText}`);
            }
            return response.json();
        })
        .then(chatMessages => {
            if (!chatMessages || chatMessages.length === 0) {
                messages.innerHTML = `<p>Type a message to start chatting with ${username}!</p>`;
                return;
            }
            chatMessages.forEach(message => {
                const messageElement = document.createElement("p");
                messageElement.textContent = `${message.sender}: ${message.text}`;
                messages.appendChild(messageElement);
            });

            // Scroll to the bottom of the chatbox
            messages.scrollTop = messages.scrollHeight;
        })
        .catch(error => console.error("Error fetching chat messages:", error));
}
/* END OF WEBSOCKET FOR CHAT */

addEventListener("DOMContentLoaded", async function () {
    console.log('DOMContentLoaded');

    // Call the function to check session status
    const sessionActive = await checkSessionActive();
    if (sessionActive) {
        console.log('Session is active');
        fetchCategories();
        fetchPosts();
        connect();
        fetchOnlineUsers();
        showAuthenticatedContainer();
    } else {
        console.log('Session is NOT active'); 
        showNotAuthenticatedContainer();
    }
});