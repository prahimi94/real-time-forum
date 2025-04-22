let categories = [];
let posts = [];
let toast;
let loggedInUser;

// Function to check if the session is active
async function checkSessionActive() {
    try {
        // Call the backend API to check session status
        const response = await fetch('/api/check-session', {
            method: 'GET',
            credentials: 'include', // Include cookies in the request
        });

        if (response.ok) {
            const res = await response.json();
            const data = res.data;
            if (data.Active) {
                loggedInUser = data.LoginUser;
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
    <div class="containerAuth">
        <div class="form-box login">
            <form method="post" id="loginForm">
                <h1>Login</h1>
                <div class="input-box">
                    <input type="text" placeholder="Username Or Email" name="username" required>
                    <i class='bx bxs-user'></i>
                </div>
                <div class="input-box">
                    <input id="login-password" type="password" placeholder="Password" name="password" required>
                    <i id="toggle-login-password-icon" class='bx bxs-show' style="cursor: pointer;" ></i>
                </div>
                <!-- <div class="forgot-link">
                    <a href="#">Forgot Password?</a>
                </div> -->
                <button type="submit" onclick="loginFunc()" class="btnAuth">Login</button>
            </form>
        </div>

        <div class="form-box register">
            <form method="post" id="registerForm">
                <h1>Registration</h1>
                <div class="input-box">
                    <input type="text" name="username" placeholder="Nick Name" required>
                    <i class='bx bxs-user'></i>
                </div>
                <div class="input-box">
                    <input type="text" name="firstname" placeholder="First Name" required>
                    <i class='bx bxs-user'></i>
                </div>
                <div class="input-box">
                    <input type="text" name="lastname" placeholder="Last Name" required>
                    <i class='bx bxs-user'></i>
                </div>
                <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="gender" id="genderMale" value="male">
                    <label class="form-check-label" for="genderMale">
                        Male
                        <i class='bx bx-male'></i>
                    </label>
                </div>
                <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="gender" id="genderFemale" value="female">
                    <label class="form-check-label" for="genderFemale">
                        Female
                        <i class='bx bx-female'></i>
                    </label>
                </div>
                <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="gender" id="genderNeutral" value="neutral" checked>
                    <label class="form-check-label" for="genderNeutral">
                        Neutral
                        <i class='bx bx-body'></i>
                    </label>
                </div>
                <div class="input-box">
                    <input type="number" name="age" placeholder="Age" required>
                    <i class='bx bxs-calendar'></i>
                </div>
                <div class="input-box">
                    <input type="email" name="email" placeholder="Email" required>
                    <i class='bx bxs-envelope' ></i>
                </div>
                <div class="input-box">
                    <input id="register-password" type="password" placeholder="Password" name="password" required>
                    <i id="toggle-register-password-icon" class='bx bxs-show'  style="cursor: pointer;" ></i>
                </div>
                <button type="submit" onclick="registerFunc()" class="btnAuth">Register</button>
            </form>
        </div>

        <div class="toggle-box">
            <div class="toggle-panel toggle-left">
                <h1 style="color: #fff;">Hello, Welcome!</h1>
                <p>Don't have an account?</p>
                <button class="btnAuth register-btn">Register</button>
            </div>

            <div class="toggle-panel toggle-right">
                <h1 style="color: #fff;">Welcome Back!</h1>
                <p>Already have an account?</p>
                <button class="btnAuth login-btn">Login</button>
            </div>
        </div>
    </div>    
    `;


    const container = document.querySelector('.containerAuth');
    const registerBtn = document.querySelector('.register-btn');
    const loginBtn = document.querySelector('.login-btn');

    registerBtn.addEventListener('click', () => {
        container.classList.add('active');
    })

    loginBtn.addEventListener('click', () => {
        container.classList.remove('active');
    })

    const toggleLoginPasswordIcon = document.getElementById('toggle-login-password-icon');
    const loginPasswordInput = document.getElementById('login-password');

    toggleLoginPasswordIcon.addEventListener('click', function() {
        // Toggle password visibility
        if (loginPasswordInput.type === 'password') {
            loginPasswordInput.type = 'text'; // Show password
            toggleLoginPasswordIcon.classList.remove('bxs-show'); // Remove eye icon
            toggleLoginPasswordIcon.classList.add('bxs-hide');
        } else {
            loginPasswordInput.type = 'password'; // Hide password
            toggleLoginPasswordIcon.classList.remove('bxs-hide'); // Remove eye icon
            toggleLoginPasswordIcon.classList.add('bxs-show');
        }
    });

    const toggleRegisterPasswordIcon = document.getElementById('toggle-register-password-icon');
    const loginRegisterInput = document.getElementById('register-password');

    toggleRegisterPasswordIcon.addEventListener('click', function() {
        // Toggle password visibility
        if (loginRegisterInput.type === 'password') {
            loginRegisterInput.type = 'text'; // Show password
            toggleRegisterPasswordIcon.classList.remove('bxs-show'); // Remove eye icon
            toggleRegisterPasswordIcon.classList.add('bxs-hide');
        } else {
            loginRegisterInput.type = 'password'; // Hide password
            toggleRegisterPasswordIcon.classList.remove('bxs-hide'); // Remove eye icon
            toggleRegisterPasswordIcon.classList.add('bxs-show');
        }
    });
}

function showAuthenticatedContainer() {
    const forumContainer = document.getElementsByTagName('main')[0];
    const loggedInUserProfilePhoto = loggedInUser.profile_photo
        ? `<img src="/uploads/${loggedInUser.profile_photo}" alt="user image" class="rounded w-50 shadow">`
        : `<i class="fa-solid fa-user rounded shadow" style="font-size: 4rem;padding: 1rem;"></i>`;

    const loggedInUserProfilePhoto2 = loggedInUser.profile_photo
        ? `<img src="/uploads/${loggedInUser.profile_photo}" alt="user image" style="width: 43px; height: 43px; cursor: pointer;" >`
        : `<i class="fa-solid fa-user"></i>`;


    forumContainer.innerHTML = `
    <header>
            <nav class="navbar navbar-expand-lg navbar-light py-10" aria-label="Thirteenth navbar example">
                <div class="container-fluid">
                    <a class="navbar-brand" href="/">
                        <img src="/img/logo.png" alt="Forum" class="me-2" style="height: 60px;">
                    </a>
                <!-- <button class="navbar-toggler" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasNavbar" aria-controls="offcanvasNavbar" aria-label="Toggle navigation">
                        <span class="navbar-toggler-icon"></span>
                    </button> -->
                <div class="offcanvas offcanvas-end" tabindex="-1" id="offcanvasNavbar" aria-labelledby="offcanvasNavbarLabel">
                        <div class="offcanvas-header">
                        <h5 class="offcanvas-title" id="offcanvasNavbarLabel"><img src="/img/logo.png" alt="Forum" class="me-2" style="height: 50px;"></h5>
                        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
                        </div>
                        <div class="offcanvas-body">
                            <ul class="navbar-nav justify-content-end flex-grow-1 pe-3">
                                <li class="nav-item">
                                    <div class="me-3 text-center">
                                        ${loggedInUserProfilePhoto}
                                    </div>
                                </li>
                                <li class="nav-item text-center pt-3">                                    
                                    <span class="me-3">Welcome, ${loggedInUser.firstname}</span>
                                </li>
                                <li class="nav-item text-center pb-2 pt-2">
                                    <span class="me-3">${loggedInUser.email}</span>
                                </li>
                                <li class="nav-item text-center pb-3">
                                
                                  <a type="button" class="btn btn-outline-secondary" onclick="showProfile()"><i class="fa-regular fa-address-card"></i></a>
                                    <a type="button" class="btn btn-danger" onclick="logoutFunc()"><i class="fas fa-power-off"></i></a>
                                </li>
                                <li><div class="divaider mb-2"></div></li>
                                    <li class="nav-item">
                                        <!--<form id="searchForm" class="d-flex mt-3" role="search">
                                        <input class="form-control me-2" type="search" placeholder="Search" aria-label="Search" name="post_info">
                                        <button onclick="filterPosts()" class="btn btn-outline-success"><i class="fa-solid fa-magnifying-glass"></i></button>
                                        </form> -->
                                    </li>
                                <li><div class="divaider mb-2 mt-2"></div></li>
                                <li class="nav-item">
                                    <a class="nav-link" href="/">Home</a>
                                </li>
                                <li class="nav-item">
                                    <!-- todo -->
                                    <!-- <a class="nav-link" href="/newPost/">Create Post</a> -->
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
                                <form id="searchForm" class="d-flex" role="search">
                                    <input class="form-control me-2 input-p" type="search" placeholder="Search" aria-label="Search" name="post_info">
                                    <button onclick="filterPosts()" class="btn btn-outline-success"><i class="fa-solid fa-magnifying-glass"></i></button>
                                </form>
                            </div>
                        </div>
                        <div class="d-lg-flex col-sm-12 col-lg-2 justify-content-lg-end">
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

                                <div class="btn-group">
                                    <button type="button" class="btn btn-secondary dropdown-toggle dropdown-toggle-user rounded-circle" style="width: 43px; height: 43px;padding: 0;overflow: hidden;" data-bs-toggle="dropdown" aria-expanded="false">
                                        <div class="rounded-circle">
                                            ${loggedInUserProfilePhoto2}
                                        </div>
                                    </button>
                                    <ul class="dropdown-menu dropdown-menu-end">
                                        <li><a type="button" class="dropdown-item" onclick="showProfile()"><i class="fa-regular fa-address-card me-2"></i> Profile</a></li> 
                                        <li><a type="button" class="dropdown-item" onclick="logoutFunc()"><i class="fas fa-power-off me-2"></i> Log Out</a></li>
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
                <div class="col-md-6" id="middlePanel">
                    <!-- {template "posts" .}} -->
                    <!-- this part is filled in script.js by fetchPosts -->
                    <div class="row" id="newPostContainer">
                        <div class="col-sm-12 col-md-12 mb-3">
                            <div class="post-card">
                                <h4 class="text-center mb-4">New Post</h4>
                                <form id="newPostForm" enctype="multipart/form-data">
                                    <div class="mb-3">
                                        <div class="mt-3">
                                            <div class="custom-multiselect" id="custom-multiselect-submit">
                                                <div class="select-box" id="categories" onclick="categoriesClicked('custom-multiselect-submit')">Select categories</div>
                                                <div class="options-list" onchange="categoriyListChanged('categoriesList')" id="categoriesList"> </div>
                                            </div>
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
                                    <button onclick="submitPost()" class="btn btn-success w-100">Post</button>
                                </form>
                            </div>
                        </div>
                    </div>
                    <div class="row" id="postsContainer"></div>
                </div>
                <div class="col-md-3">
                    <div class="col-right-userLogin">
                        <div class="text-center w-100 py-4">
                        ${loggedInUserProfilePhoto}
                        </div>
                        <div class="text-center title-username">${loggedInUser.username}</div>
                        <div class="info-box-username">
                            <p> ${loggedInUser.email}</p>
                            <div class="info-box-logout">
                                <p style="margin-bottom: 1rem;"><a type="button" onclick="showProfile()"><i
                                            class="fa-regular fa-address-card me-2"></i> Profile</a></p>
                                <a type="button" onclick="logoutFunc()"><i class="fas fa-power-off me-2" style="color: #c44343;"></i> Log
                                    Out</a>
                            </div>
                        </div>
                        <!-- ALL CHAT USERS -->
                        <div id="chat-users" class="chat-users-container">
                            <h3>All Chat Users</h3>
                            <ul id="chat-users-list"></ul>
                        </div>
                        <!-- CHATBOX -->
                        <div id="chatbox" class="chatbox" style="display: none">
                            <h3 id="chat-header"></h3>
                            <div id="messages" style="display: none">
                                <div id="message-display"></div>
                                <div id="typing" style="display: none;"></div>
                                <input type="text" id="messageInput" style="display: none">
                                <button id="send-btn" style="display: none">Send</button>
                            </div>
                        </div>
                        <!-- END CHATBOX -->
                    </div>
                </div>
            </div>
        </div>
    `;


    // laod categories in select for new post
    const selectCategoriesContainer = document.getElementById('categoriesList');
    const categoryOptions = categories.map(category =>
        `<label><input type="checkbox" name="categories" value="${category.id}" /> ${category.name}</label>`
    ).join('');
    selectCategoriesContainer.innerHTML = categoryOptions;


    const multiSelect = document.querySelector(".custom-multiselect");
    document.addEventListener("click", (e) => {
        if (!multiSelect.contains(e.target)) {
            multiSelect.classList.remove("open");
        }
    });
}

async function fetchCategories() {
    const response = await fetch('/api/categories/');
    res = await response.json();
    if (!res.success) {
        showToast(res.message);
        return;
    }
    categories = res.data

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
        categoryATag.href = `javascript:fetchCategoryPosts('${category.name}')`;
        categoryATag.appendChild(categoryITag);
        categoryATag.innerHTML += category.name;

        const categoryLiTag = document.createElement('li');
        categoryLiTag.appendChild(categoryATag);

        categoriesContainer.appendChild(categoryLiTag);
    });

    // laod categories in select for new post
    const selectCategoriesContainer = document.getElementById('categoriesList');
    const categoryOptions = categories.map(category =>
        `<label><input type="checkbox" name="categories" value="${category.id}" /> ${category.name}</label>`
    ).join('');
    selectCategoriesContainer.innerHTML = categoryOptions;

    const multiSelect = document.querySelector(".custom-multiselect");
    document.addEventListener("click", (e) => {
        if (!multiSelect.contains(e.target)) {
            multiSelect.classList.remove("open");
        }
    });
}

async function fetchPost(postId, postUuid) {
    const response = await fetch('/api/post/' + postUuid);
    res = await response.json();
    if (!res.success) {
        showToast(res.message);
        return;
    }

    const post = res.data.Post;
    const comments = res.data.Comments;
    updatePostHtml(post, comments, postId);
}

async function fetchPosts() {
    const response = await fetch('/api/bunchOfPosts/1');
    res = await response.json();
    if (!res.success) {
        showToast(res.message);
        return;
    }
    posts = res.data

    fillPostsInHtml(posts);
}

async function filterPosts() {
    const form = document.getElementById('searchForm');
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });

    const searchTerm = document.querySelector('input[name="post_info"]').value;
    const response = await fetch('/api/filterPosts/' + searchTerm);
    res = await response.json();
    if (!res.success) {
        showToast(res.message);
        return;
    }
    posts = res.data

    fillPostsInHtml(posts.Posts, 'Search results for: ' + searchTerm);
}

async function fetchMyCreatedPosts() {
    const response = await fetch('/api/myCreatedPosts/');
    res = await response.json();
    showToast(res);
    posts = res.data

    fillPostsInHtml(posts.Posts, 'My created post');
}

async function fetchMyLikedPosts() {
    const response = await fetch('/api/myLikedPosts/');
    res = await response.json();
    showToast(res);
    posts = res.data

    fillPostsInHtml(posts.Posts, 'My liked post');
}

async function fetchCategoryPosts(category) {
    const response = await fetch('/api/posts/' + category);
    res = await response.json();
    showToast(res);
    posts = res.data

    fillPostsInHtml(posts.Posts, category + ' posts');
}

async function showProfile() {
    const middlePanel = document.getElementById('middlePanel');

    const loggedInUserProfilePhoto = loggedInUser.profile_photo
        ? `<img src="/uploads/${loggedInUser.profile_photo}" alt="user image" class="rounded-circle center" style="width: 150px; height: 150px; cursor: pointer;" >`
        : `<i class="fa-solid fa-user" style="font-size: 3rem;"></i>`;

    const genderMaleChecked = loggedInUser.gender == 'male'
        ? `checked` : ``;
    const genderFemaleChecked = loggedInUser.gender == 'female'
        ? `checked` : ``;
    const genderNeutralChecked = loggedInUser.gender == 'neutral'
        ? `checked` : ``;

    middlePanel.innerHTML = `
    <div class="container-form">

        <div class="mb-3" style="text-align: center;">
            <div class="rounded-circle">
                ${loggedInUserProfilePhoto}
            </div>
        </div>

        <h4 class="text-center mb-4">Edit Profile</h4>
        
        <form id="profileForm" enctype="multipart/form-data" method="post">
            
            <div class="mb-3">
                <label for="firstname">First Name:</label>
                <input type="text" class="form-control" placeholder="First Name" name="firstname" value="${loggedInUser.firstname}">
            </div>
            <div class="mb-3">
                <label for="lastname">Last Name:</label>
                <input type="text" class="form-control" placeholder="Last Name" name="lastname" value="${loggedInUser.lastname}">
            </div>
            <div class="mb-3">
                <label for="gender">Gender:</label>
                <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="gender" id="genderMale" value="male" ${genderMaleChecked}>
                    <label class="form-check-label" for="genderMale">
                        Male
                        <i class='bx bx-male'></i>
                    </label>
                </div>
                <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="gender" id="genderFemale" value="female" ${genderFemaleChecked}>
                    <label class="form-check-label" for="genderFemale">
                        Female
                        <i class='bx bx-female'></i>
                    </label>
                </div>
                <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="gender" id="genderNeutral" value="neutral" ${genderNeutralChecked}>
                    <label class="form-check-label" for="genderNeutral">
                        Neutral
                        <i class='bx bx-body'></i>
                    </label>
                </div>
            </div>
            <div class="mb-3">
                <label for="age">Age:</label>
                <input type="number" class="form-control" placeholder="Age" name="age" value="${loggedInUser.age}">
            </div>
            <label for="profile_photo">Profile photo:</label>
            <div class="mb-3 text-center border p-3 rounded">
                <p class="text-muted">Attach an image (optional)</p>
                <input type="file" class="form-control" name="profile_photo">
            </div>
            <button type="submit" onclick="updateUser()" class="btn btn-success w-100">Edit user</button>
        </form>
    </div>
    `
}

async function updateUser() {
    const form = document.getElementById('profileForm');
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const response = await fetch('/api/updateUser', {
        method: 'POST',
        body: new FormData(form),
    });
    res = await response.json();
    showToast(res);

    if (res.success) {
        await checkSessionActive()
    }
}

async function submitPost() {
    const form = document.getElementById('newPostForm');
    const optionsList = document.getElementById("categoriesList");
    const selectBox = document.getElementById("categories");

    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });

    const formData = new FormData(form);

    // Add selected categories
    optionsList.querySelectorAll('input[name="categories"]:checked').forEach(input => {
        formData.append('categories[]', input.value);
    });

    const response = await fetch('/api/submitPost', {
        method: 'POST',
        body: formData,
    });

    const res = await response.json();
    showToast(res);

    form.reset();
    optionsList.querySelectorAll('input[name="categories"]').forEach(input => {
        input.checked = false;
    });
    selectBox.textContent = "Select options";

    fetchPosts();
}

async function likePost(id, uuid, actionType) {
    const form = document.getElementById('likePostForm-' + id);
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const formData = new FormData();
    formData.append('post_id', id);
    formData.append('actionType', actionType);

    const response = await fetch('/api/likePost', {
        method: 'POST',
        body: formData,
    });
    res = await response.json();
    showToast(res);

    form.reset();
    fetchPost(id, uuid);
}

async function updatePost(id, uuid) {
    console.log('updatePost');
    const form = document.getElementById('updatePostForm-' + id);
    const optionsList = document.getElementById("updatePostCategoriesList-" + uuid);
    const selectBox = document.getElementById("update_post_categories-" + uuid);

    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });

    const formData = new FormData(form);
    console.log('updatePost formData', formData);

    // Add selected categories
    optionsList.querySelectorAll('input[name="update_post_categories-'+ uuid +'"]:checked').forEach(input => {
        formData.append('update_post_categories[]', input.value);
    });
    console.log('updatePost formData after', formData);

    const response = await fetch('/api/updatePost', {
        method: 'POST',
        body: formData,
    });
    res = await response.json();
    showToast(res);

    form.reset();

    optionsList.querySelectorAll('input[name="update_post_categories-'+ uuid +'"]').forEach(input => {
        input.checked = false;
    });
    selectBox.textContent = "Select options";

    const updatePostModal = document.getElementById('updatePostModal-' + id);
    // const modalInstance = bootstrap.Modal.getInstance(updatePostModal);
    // modalInstance.hide();
    updatePostModal.style.display = "none";
    fetchPost(id, uuid);   
}

async function deletePost(id, uuid) {
    const form = document.getElementById('deletePostForm-' + id);
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const response = await fetch('/api/deletePost', {
        method: 'POST',
        body: new FormData(form),
    });
    res = await response.json();
    showToast(res);

    form.reset();
    const deletePostModal = document.getElementById('deletPostModal-' + id);
    // const modalInstance = bootstrap.Modal.getInstance(deletePostModal);
    // modalInstance.hide();
    deletePostModal.style.display = "none";
    removePostHtml(id);
}

async function sumbitComment(id, uuid) {
    const form = document.getElementById('commentForm-' + id);
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const response = await fetch('/api/submitComment', {
        method: 'POST',
        body: new FormData(form),
    });
    res = await response.json();
    showToast(res);

    form.reset();
    fetchPost(id, uuid);
}

async function likeComment(comment_id, id, uuid, actionType) {
    const form = document.getElementById('likeCommentForm-' + comment_id);
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const formData = new FormData();
    formData.append('comment_id', comment_id);
    formData.append('actionType', actionType);

    const response = await fetch('/api/likeComment', {
        method: 'POST',
        body: formData,
    });
    res = await response.json();
    showToast(res);

    form.reset();
    fetchPost(id, uuid);
}

async function updateComment(comment_id, id, uuid) {
    const form = document.getElementById('updateCommentForm-' + comment_id);
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const response = await fetch('/api/updateComment', {
        method: 'POST',
        body: new FormData(form),
    });
    res = await response.json();
    showToast(res);

    form.reset();
    const updateCommentModal = document.getElementById('updateCommentModal-' + comment_id);
    // const modalInstance = bootstrap.Modal.getInstance(updateCommentModal);
    // modalInstance.hide();
    updateCommentModal.style.display = "none";
    fetchPost(id, uuid);   
}

async function deleteComment(comment_id, id, uuid) {
    const form = document.getElementById('deleteCommentForm-' + comment_id);
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const response = await fetch('/api/deleteComment', {
        method: 'POST',
        body: new FormData(form),
    });
    res = await response.json();
    showToast(res);

    form.reset();
    const deleteCommentModal = document.getElementById('deletCommentModal-' + comment_id);
    // const modalInstance = bootstrap.Modal.getInstance(deleteCommentModal);
    // modalInstance.hide();
    deleteCommentModal.style.display = "none";
    fetchPost(id, uuid);   
}

async function loginFunc() {
    const form = document.getElementById('loginForm');
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const response = await fetch('/api/login', {
        method: 'POST',
        body: new FormData(form),
    });
    res = await response.json();
    showToast(res);

    if(res.success) {
        form.reset();
        // const authModal = document.getElementById('authModal');
        // const modalInstance = bootstrap.Modal.getInstance(authModal);
        // modalInstance.hide();
        await checkSession();
    }
}

async function registerFunc() {
    const form = document.getElementById('registerForm');
    form.addEventListener('submit', (event) => {
        event.preventDefault();
    });
    
    const response = await fetch('/api/register', {
        method: 'POST',
        body: new FormData(form),
    });
    res = await response.json();
    showToast(res);

    if(res.success) {
        form.reset();
        // const authModal = document.getElementById('authModal');
        // const modalInstance = bootstrap.Modal.getInstance(authModal);
        // modalInstance.hide();
        await checkSession();
    }
}

async function logoutFunc() {
    const response = await fetch('/api/logout', {
        method: 'GET',
        credentials: 'include',
    });
    res = await response.json();
    showToast(res);

    await checkSession();
}

function fillPostsInHtml(posts, actionSubject = '', fillType = 'reWrite') {

    if (fillType == 'reWrite') {
        // load posts for home page
        const middlePanel = document.getElementById('middlePanel');
        middlePanel.innerHTML = `
        <div class="row" id="newPostContainer"></div>
        <div class="row" id="postsContainer"></div>
        `;
        const newPostContainer = document.getElementById('newPostContainer');
        const postsContainer = document.getElementById('postsContainer');
        
        newPostContainer.innerHTML = `
        <div class="row" id="newPostContainer">
            <div class="col-sm-12 col-md-12 mb-3">
                <div class="post-card">
                    <h4 class="text-center mb-4">New Post</h4>
                    <form id="newPostForm" enctype="multipart/form-data">
                        <div class="mb-3">
                            <div class="mt-3">
                                <div class="custom-multiselect" id="custom-multiselect-submit">
                                    <div class="select-box" id="categories" onclick="categoriesClicked('custom-multiselect-submit')">Select categories</div>
                                    <div class="options-list" onchange="categoriyListChanged('categoriesList', 'categories')" id="categoriesList"> </div>
                                </div>
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
                        <button onclick="submitPost()" class="btn btn-success w-100">Post</button>
                    </form>
                </div>
            </div>
        </div>
        `;
        console.log(newPostContainer.innerHTML)
        postsContainer.innerHTML = `
        <div class='col-md-12 text-center'>
            <h4 class='text-muted'>${actionSubject}</h4>
        </div>`;
        
        if (posts === null || posts.length === 0) {
            postsContainer.innerHTML += '<div class="col-md-12 text-center">No posts found!</div>';
            return;
        }
        
        postsContainer.innerHTML += '<div class="accordion accordion-flush" id="accordionFlushExample">';        
    } else {
        const postsContainer = document.getElementById('postsContainer');
    }
    
    posts.forEach(post => {
        const postImage = post.user.profile_photo
            ? `<img class="bd-placeholder-img flex-shrink-0 me-2 rounded" role="img" src="/uploads/${post.user.profile_photo}" width="45" height="45"/>`
            : `<div style="padding: 7px;"><i class="fa-solid fa-user" style="font-size: 2rem;"></i></div>`;

        const postCategories = post.categories.map(category =>
            `<span class="badge-p text-dark"><a href="javascript:fetchCategoryPosts('${category.name}')">${category.name}</a></span>`
        ).join('');

        const postFiles = post.post_files.map(post_file =>
            `<div class="col-md-12">
                <img src="/uploads/${post_file.file_uploaded_name}" alt="post image" class="rounded mb-1" style="width: 100%; max-height: 400px;">
            </div>`
        ).join('');

        const formattedDateTime = post.created_at.replace('T', ' ').replace('Z', '');

        const postButtons = 
        post.user_id === loggedInUser.id
                        ?`
                        <div style="float: right;margin-top: -16px;">
                            <div class="row py-3 ms-2">
                                <div class="btn-group">
                                    <div class="mydropdown">
                                        <a class="mydropdown-toggle btn btn-light" type="button">
                                            <i class="fa-solid fa-ellipsis"></i>
                                        </a>
                                        <ul class="mydropdown-menu dropdown-menu dropdown-menu-end" style="border: 1px solid #c2c2c270;">
                                            <li>
                                                <a type="button" class="dropdown-item open-modal" data-target="updatePostModal-${post.id}">
                                                    <i class="fas fa-edit me-2"></i>Edit Post
                                                </a>                                        
                                            </li>
                                            <li>
                                                <a type="button" class="dropdown-item open-modal" data-target="deletPostModal-${post.id}"><i class="fa-solid fa-trash me-2"></i>Delete Post</a> 
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <!-- Modal -->
                        <div id="updatePostModal-${post.id}" class="modal fade bd-example-modal-lg">
                            <div class="modal-dialog modal-lg">
                                <div class="modal-content">
                                    <div class="modal-header">
                                    <h1 class="modal-title fs-5" id="editPostModalLabel">Edit post</h1>
                                    <a type="button" class="btn-close close-modal" data-bs-dismiss="modal" aria-label="Close"></a>
                                    </div>
                                    <form id="updatePostForm-${post.id}" method="post" enctype="multipart/form-data">
                                    <div class="modal-body">
                                            <input type="hidden" name="post_uuid" value="${post.uuid}">
                                            <input type="hidden" name="post_id" value="${post.id}">
                                            <div class="mb-3">
                                                <div class="mt-3">
                                                    <div class="custom-multiselect" id="custom-multiselect-update-${post.uuid}">
                                                        <div class="select-box" id="update_post_categories-${post.uuid}" onclick="categoriesClicked('custom-multiselect-update-${post.uuid}')">Select categories</div>
                                                        <div class="options-list" onchange="categoriyListChanged('updatePostCategoriesList-${post.uuid}', 'update_post_categories-${post.uuid}')" id="updatePostCategoriesList-${post.uuid}"> </div>
                                                    </div>
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
                                    </div>
                                    <div class="modal-footer">
                                    <a type="button" class="btn btn-secondary close-modal" data-bs-dismiss="modal">Close</a>
                                    <a onclick="updatePost(${post.id}, '${post.uuid}')" class="btn btn-success">Save changes</a>
                                </form>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="modal fade" id="deletPostModal-${post.id}" tabindex="-1" aria-labelledby="deletPostModalLabel-${post.id}" aria-hidden="true">
                            <form id="deletePostForm-${post.id}" method="post">
                                <input type="hidden" name="id" value="${post.id}">
                                <div class="modal-dialog modal-dialog-centered">
                                    <div class="modal-content">
                                        <div class="modal-header bg-danger text-white">
                                            <h5 class="modal-title" id="deletPostModalLabel-${post.id}">Confirm Deletion</h5>
                                            <a type="button" class="btn-close btn-close-white close-modal" data-bs-dismiss="modal" aria-label="Close"></a>
                                        </div>
                                        <div class="modal-body">
                                            <p class="mb-0">Are you sure you want to delete this item? This action cannot be undone.</p>
                                        </div>
                                        <div class="modal-footer">
                                            <a type="button" class="btn btn-secondary close-modal" data-bs-dismiss="modal">Cancel</a>
                                            <a onclick="deletePost(${post.id}, '${post.uuid}')" class="btn btn-danger" id="confirmDelete">Delete</a>
                                        </div>
                                    </div>
                                </div>
                            </form>
                        </div>
                        `
                        : ``;

        const postLikeElement = post.liked
            ? `<a onclick="likePost(${post.id}, '${post.uuid}', 'like')" value="like" name="like_post" class="btn btn-success"><i class="fa-solid fa-thumbs-up"></i></a>`
            : `<a onclick="likePost(${post.id}, '${post.uuid}', 'like')" value="like" name="like_post" class="btn btn-outline-success"><i class="fa-regular fa-thumbs-up"></i></a>`;

        const postDislikeElement = post.disliked
            ? `<a onclick="likePost(${post.id}, '${post.uuid}', 'dislike')" value="dislike" name="dislike_post" value="dislike" class="btn btn-danger"><i class="fa-solid fa-thumbs-down"></i></a>`
            : `<a onclick="likePost(${post.id}, '${post.uuid}', 'dislike')" value="dislike" name="dislike_post" value="dislike" class="btn btn-outline-danger"><i class="fa-regular fa-thumbs-down"></i></a>`;

        const postHTML = `
            <div class="accordion-item">
                <h2 class="accordion-header" id="flush-heading-${post.id}">
                    <div class="col-sm-12 col-md-12 mb-3">
                        <div class="post-card" id="post-${post.id}">
                            <div onclick="accordionHeadedClicked(${post.id})">
                                <button onclick="fetchPost(${post.id}, '${post.uuid}')" class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapse-${post.id}" aria-expanded="false" aria-controls="flush-collapseOne">

                                <div class="d-flex flex-column justify-content-between">
                                    <h5 class="mt-2 post-title">
                                        ${post.title}
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

                                </button>
                            </div>


                            <div class="mt-4">
                                <span class="like-inpost"><i class="fa-solid fa-thumbs-up"></i> ${post.number_of_likes}</span>
                                <span class="dislike-inpost"><i class="fa-solid fa-thumbs-down"></i> ${post.number_of_dislikes}</span>
                                
                                ${postButtons}

                                <div style="float: right;margin-top: -16px;">
                                    <div class="row py-3">
                                        <form id="likePostForm-${post.id}">
                                            <input type="hidden" id="post_id" name="post_id" value="${post.iD}">

                                            ${postLikeElement}
                                        

                                            ${postDislikeElement}
                                        </form>
                                    </div>
                                </div>

                            </div> <!-- mt-4 -->
                        </div> <!-- post-card -->
                    </div> <!-- col-sm-12 col-md-12 mb-3 -->
                </h2>

                <div id="flush-collapse-${post.id}" class="accordion-collapse collapse" aria-labelledby="flush-heading-${post.id}" data-bs-parent="#accordionFlushExample">
                    <div class="accordion-body post-card bg-border-box" id="post-comments-${post.id}">
                    </div>
                </div>
            </div>
        `;
    
        if (fillType == 'reWrite') {
            postsContainer.insertAdjacentHTML('beforeend', postHTML);
        } else {
            postsContainer.innerHTML += postHTML;
        }
    });

    if(fillType == 'reWrite') {
        // laod categories in select for new post
        const selectCategoriesContainer = document.getElementById('categoriesList');
        const categoryOptions = categories.map(category =>
            `<label><input type="checkbox" name="categories" value="${category.id}" /> ${category.name}</label>`
        ).join('');
        selectCategoriesContainer.innerHTML = categoryOptions;


        // laod categories in select for new post
        const updatePostCategoryOptions = categories.map(category =>
            `<label><input type="checkbox" name="update_post_categories" value="${category.id}" /> ${category.name}</label>`
        ).join('');
        document.querySelectorAll('[id^="updatePostCategoriesList"]').forEach(container => {
            container.innerHTML = updatePostCategoryOptions;
        });
    }

    
}

function updatePostHtml(post, comments, postId) {
    const postImage = post.user.profile_photo
            ? `<img class="bd-placeholder-img flex-shrink-0 me-2 rounded" role="img" src="/uploads/${post.user.profile_photo}" width="45" height="45"/>`
            : `<div style="padding: 7px;"><i class="fa-solid fa-user" style="font-size: 2rem;"></i></div>`;

    const postCategories = post.categories.map(category =>
        `<span class="badge-p text-dark"><a href="javascript:fetchCategoryPosts('${category.name}')">${category.name}</a></span>`
    ).join('');

    const postFiles = post.post_files.map(post_file =>
        `<div class="col-md-12">
            <img src="/uploads/${post_file.file_uploaded_name}" alt="post image" class="rounded mb-1" style="width: 100%; max-height: 400px;">
        </div>`
    ).join('');

    const formattedDateTime = post.created_at.replace('T', ' ').replace('Z', '');

    const postButtons = 
    post.user_id === loggedInUser.id
                    ?`
                    <div style="float: right;margin-top: -16px;">
                        <div class="row py-3 ms-2">
                            <div class="btn-group">
                                <div class="mydropdown">
                                    <a class="mydropdown-toggle btn btn-light" type="button">
                                        <i class="fa-solid fa-ellipsis"></i>
                                    </a>
                                    <ul class="mydropdown-menu dropdown-menu dropdown-menu-end" style="border: 1px solid #c2c2c270;">
                                        <li>
                                            <a type="button" class="dropdown-item open-modal" data-target="updatePostModal-${post.id}">
                                                <i class="fas fa-edit me-2"></i>Edit Post
                                            </a>                                        
                                        </li>
                                        <li>
                                            <a type="button" class="dropdown-item open-modal" data-target="deletPostModal-${post.id}"><i class="fa-solid fa-trash me-2"></i>Delete Post</a> 
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                    <!-- Modal -->
                    <div class="modal fade bd-example-modal-lg" id="updatePostModal-${post.id}" tabindex="-1" aria-labelledby="editPostModalLabel" aria-hidden="true">
                        <div class="modal-dialog modal-lg">
                            <div class="modal-content">
                                <div class="modal-header">
                                <h1 class="modal-title fs-5" id="editPostModalLabel">Edit post</h1>
                                <a type="button" class="btn-close close-modal" data-bs-dismiss="modal" aria-label="Close"></a>
                                </div>
                                <form id="updatePostForm-${post.id}" method="post" enctype="multipart/form-data">
                                <div class="modal-body">
                                        <input type="hidden" name="post_uuid" value="${post.uuid}">
                                        <input type="hidden" name="post_id" value="${post.id}">
                                         <div class="mb-3">
                                            <div class="mt-3">
                                                <div class="custom-multiselect" id="custom-multiselect-update-${post.uuid}">
                                                    <div class="select-box" id="update_post_categories-${post.uuid}" onclick="categoriesClicked('custom-multiselect-update-${post.uuid}')">Select categories</div>
                                                    <div class="options-list" onchange="categoriyListChanged('updatePostCategoriesList-${post.uuid}', 'update_post_categories-${post.uuid}')" id="updatePostCategoriesList-${post.uuid}"> </div>
                                                </div>
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
                                </div>
                                <div class="modal-footer">
                                <a type="button" class="btn btn-secondary close-modal" data-bs-dismiss="modal">Close</a>
                                <a onclick="updatePost(${post.id}, '${post.uuid}')" class="btn btn-success">Save changes</a>
                            </form>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="modal fade" id="deletPostModal-${post.id}" tabindex="-1" aria-labelledby="deletPostModalLabel-${post.id}" aria-hidden="true">
                        <form id="deletePostForm-${post.id}" method="post">
                            <input type="hidden" name="id" value="${post.id}">
                            <div class="modal-dialog modal-dialog-centered">
                                <div class="modal-content">
                                    <div class="modal-header bg-danger text-white">
                                        <h5 class="modal-title" id="deletPostModalLabel-${post.id}">Confirm Deletion</h5>
                                        <button type="button" class="btn-close btn-close-white close-modal" data-bs-dismiss="modal" aria-label="Close"></button>
                                    </div>
                                    <div class="modal-body">
                                        <p class="mb-0">Are you sure you want to delete this item? This action cannot be undone.</p>
                                    </div>
                                    <div class="modal-footer">
                                        <button type="button" class="btn btn-secondary close-modal" data-bs-dismiss="modal">Cancel</button>
                                        <button onclick="deletePost(${post.id}, '${post.uuid}')" class="btn btn-danger" id="confirmDelete">Delete</button>
                                    </div>
                                </div>
                            </div>
                        </form>
                    </div>
                    `
                    : ``;

    const postLikeElement = post.liked
        ? `<button onclick="likePost(${post.id}, '${post.uuid}','like')" value="like" name="like_post" class="btn btn-success"><i class="fa-solid fa-thumbs-up"></i></button>`
        : `<button onclick="likePost(${post.id}, '${post.uuid}','like')" value="like" name="like_post" class="btn btn-outline-success"><i class="fa-regular fa-thumbs-up"></i></button>`;

    const postDislikeElement = post.disliked
        ? `<button onclick="likePost(${post.id}, '${post.uuid}','dislike')" value="dislike" name="dislike_post" value="dislike" class="btn btn-danger"><i class="fa-solid fa-thumbs-down"></i></button>`
        : `<button onclick="likePost(${post.id}, '${post.uuid}','dislike')" value="dislike" name="dislike_post" value="dislike" class="btn btn-outline-danger"><i class="fa-regular fa-thumbs-down"></i></button>`;


    const postElement = document.getElementById('post-' + postId);
    const postCommentsElement = document.getElementById('post-comments-' + postId);
    if (postElement) {
        postElement.innerHTML = `
        <div onclick="accordionHeadedClicked(${post.id})">
            <button onclick="fetchPost(${post.id}, '${post.uuid}')" class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapse-${post.id}" aria-expanded="false" aria-controls="flush-collapseOne">
                <div class="d-flex flex-column justify-content-between">
                    <h5 class="mt-2 post-title">
                        ${post.title}
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

            </button>
        </div>
            

            <div class="mt-4">
                <span class="like-inpost"><i class="fa-solid fa-thumbs-up"></i> ${post.number_of_likes}</span>
                <span class="dislike-inpost"><i class="fa-solid fa-thumbs-down"></i> ${post.number_of_dislikes}</span>
                
                ${postButtons}

                <div style="float: right;margin-top: -16px;">
                    <div class="row py-3">
                        <form id="likePostForm-${post.id}">
                            <input type="hidden" id="post_id" name="post_id" value="${post.iD}">

                            ${postLikeElement}
                        

                            ${postDislikeElement}
                        </form>
                    </div>
                </div>

            </div>
        `;

        let postCommentsHtml = `<div class="row mb-4">
                                    <div class="col-md-12">
                                        <h3>Comments</h3>
                                        <div class="divaider"></div>
                                    </div>
                                </div>`;

        if(!comments || comments.length === 0) {
            postCommentsHtml += '<div class="col-md-12 text-center">No comments found!</div>';
        } else {
            comments.forEach(comment => {
                const commentDateTime = comment.created_at.replace('T', ' ').replace('Z', '');
    
                const commentButtons = 
                comment.user_id === loggedInUser.id
                    ? `<div style="float: right;margin-top: -16px;">
                            <div class="row py-3 ms-2">
                                <div class="btn-group">
                                    <div class="mydropdown">
                                        <a class="btn btn-light mydropdown-toggle" type="button">
                                            <i class="fa-solid fa-ellipsis"></i>
                                        </a>
                                        <ul class="mydropdown-menu dropdown-menu dropdown-menu-end" style="border: 1px solid #c2c2c270;">
                                            <li>
                                                <!-- Button trigger modal -->
                                                <a type="button" class="dropdown-item open-modal" data-target="updateCommentModal-${comment.id}">
                                                    <i class="fas fa-edit me-2"></i>Edit Comment
                                                </a>
                                                
                                            </li>
                                            <li>
                                                <a type="button" class="dropdown-item open-modal" data-target="deletCommentModal-${comment.id}"><i class="fa-solid fa-trash me-2"></i>Delete Comment</a> 
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <!-- Modal -->
                        <div class="modal fade" id="updateCommentModal-${comment.id}" tabindex="-1" aria-labelledby="editCommentModalLabel" aria-hidden="true">
                            <div class="modal-dialog">
                                <div class="modal-content">
                                    <div class="modal-header">
                                    <h1 class="modal-title fs-5" id="editCommentModalLabel">Edit comment</h1>
                                    <a type="button" class="btn-close close-modal" data-bs-dismiss="modal" aria-label="Close"></a>
                                    </div>
                                    <form id="updateCommentForm-${comment.id}" method="post" enctype="multipart/form-data">
                                    <div class="modal-body">
                                            <input type="hidden" name="post_uuid" value="${post.uuid}">
                                            <input type="hidden" name="comment_id" value="${comment.id}">
                                            <div class="mb-3">
                                            <label for="description-text" class="col-form-label">Comment:</label>
                                            <textarea class="form-control" id="description-text" name="description">${comment.description}</textarea>
                                            </div>
                                    </div>
                                    <div class="modal-footer">
                                    <a type="button" class="btn btn-secondary close-modal" data-bs-dismiss="modal">Close</a>
                                    <a onclick="updateComment(${comment.id}, ${post.id}, '${post.uuid}')" class="btn btn-success">Save changes</a>
                                    </form>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="modal fade" id="deletCommentModal-${comment.id}" tabindex="-1" aria-labelledby="deletCommentModalLabel" aria-hidden="true">
                            <form id="deleteCommentForm-${comment.id}" method="post">
                                <input type="hidden" name="post_uuid" value="${post.uuid}">
                                <input type="hidden" name="comment_id" value="${comment.id}">
                                <div class="modal-dialog modal-dialog-centered">
                                    <div class="modal-content">
                                        <div class="modal-header bg-danger text-white">
                                            <h5 class="modal-title" id="deletCommentModalLabel">Confirm Deletion</h5>
                                            <a type="button" class="btn-close btn-close-white close-modal" data-bs-dismiss="modal" aria-label="Close"></a>
                                        </div>
                                        <div class="modal-body">
                                            <p class="mb-0">Are you sure you want to delete this item? This action cannot be undone.</p>
                                        </div>
                                        <div class="modal-footer">
                                            <a type="button" class="btn btn-secondary close-modal" data-bs-dismiss="modal">Cancel</a>
                                            <a onclick="deleteComment(${comment.id}, ${post.id}, '${post.uuid}')" class="btn btn-danger" id="confirmDelete">Delete</a>
                                        </div>
                                    </div>
                                </div>
                            </form>
                        </div>`
                    : ``;
    
                const commentLikeElement = comment.liked
                    ? `<a onclick="likeComment(${comment.id}, ${post.id}, '${post.uuid}','like')" name="like" value="like" class="btn btn-success"><i class="fa-solid fa-thumbs-up"></i></a>`
                    : `<a onclick="likeComment(${comment.id}, ${post.id}, '${post.uuid}','like')" name="like" value="like" class="btn btn-outline-success"><i class="fa-regular fa-thumbs-up"></i></a>`
    
                const commentDislikeElement = comment.disliked
                    ? `<a onclick="likeComment(${comment.id}, ${post.id}, '${post.uuid}','dislike')" name="dislike" value="dislike" class="btn btn-danger"><i class="fa-solid fa-thumbs-down"></i></a>`
                    : `<a onclick="likeComment(${comment.id}, ${post.id}, '${post.uuid}','dislike')" name="dislike" value="dislike" class="btn btn-outline-danger"><i class="fa-regular fa-thumbs-down"></i></a>`
    
                postCommentsHtml += `
                    <div class="post-card bg-border-box">
                        <div class="row">
                            <h5 class="mt-2 post-user">${comment.user.username}</h5>
                            <p class="post-dateTime">${commentDateTime}</p>
                            <p class="post-description">${comment.description}</p>
                        </div>
                        <div class="mt-4">
                            <span class="like-inpost"><i class="fa-solid fa-thumbs-up"></i> ${comment.number_of_likes}</span>
                            <span class="dislike-inpost"><i class="fa-solid fa-thumbs-down"></i> ${comment.number_of_dislikes}</span>
    
                            ${commentButtons}
    
                            <div style="float: right;margin-top: -16px;">
                                <div class="row py-3 ms-2">
                                    <form id="likeCommentForm-${comment.id}">
                                        <input type="hidden" name="comment_id" value="${comment.id}">
                                    
                                        ${commentLikeElement}
    
                                        ${commentDislikeElement}
                                    
                                    </form>
                                </div>
                            </div>
                        </div>
                    </div>`;
            })
        }
        

        postCommentsHtml += `
                <div class="container p-5">
                <div class="row">
                    <div class="col-md-12">
                        <h3>Your Comment</h3>
                        <p class="text-secondary">You can share your ideas about this post here.</p>
                        <form id="commentForm-${post.id}" method="post">
                            <div class="mb-3">
                                <input type="hidden" id="post_id" name="post_id" value="${post.id}">
                            </div>

                            <div class="mb-3">
                                <textarea class="form-control" style="border-radius: 14px;" placeholder="Text" required rows="4" name="description"></textarea>
                            </div>
                            <button onclick="sumbitComment(${post.id}, '${post.uuid}')" class="btn btn-success w-100" style="border-radius: 14px;">Comment</button>
                        </form>
                    </div>
                </div>
            </div>
            `;

        postCommentsElement.innerHTML = postCommentsHtml
    }


    // laod categories in select for new post
    const selectUpdatePostCategoriesContainer = document.getElementById(`updatePostCategoriesList-${post.uuid}`);
    const updatePostCategoryOptions = categories.map(category =>
        `<label><input type="checkbox" name="update_post_categories" value="${category.id}" /> ${category.name}</label>`
    ).join('');
    selectUpdatePostCategoriesContainer.innerHTML = updatePostCategoryOptions;
}

function removePostHtml(postId) {
    const postElement = document.getElementById('post-' + postId);
    postElement.remove();
}

function showToast(res) {
    const toastLiveExample = document.getElementById('liveToast')

    var toastBody = toastLiveExample.querySelector('.toast-body')
    if(res.success) {
        toastLiveExample.classList.remove('bg-danger')
        toastLiveExample.classList.add('bg-success')
    } else {
        toastLiveExample.classList.remove('bg-success')
        toastLiveExample.classList.add('bg-danger')
    }
    toastBody.innerHTML = res.message
    showMyToast(res.message)
}


addEventListener("DOMContentLoaded", async function () {
    document.addEventListener("click", function (e) {
        if (e.target.closest(".mydropdown-toggle")) {
          e.stopPropagation();
          const toggle = e.target.closest(".mydropdown-toggle");
          const dropdown = toggle.closest(".mydropdown");
      
          // Close others
          document.querySelectorAll(".mydropdown").forEach((d) => {
            if (d !== dropdown) d.classList.remove("show");
          });
      
          dropdown.classList.toggle("show");
        } else {
          // Clicked outside
          document.querySelectorAll(".mydropdown").forEach((dropdown) => {
            dropdown.classList.remove("show");
          });
        }
      });

    const toastLiveExample = document.getElementById('liveToast')
    // toast = new bootstrap.Toast(toastLiveExample)


    await checkSession();
});

async function checkSession(){
    // Call the function to check session status
    const sessionActive = await checkSessionActive();
    if (sessionActive) {
        fetchCategories();
        fetchPosts();
        connect(); // websocket
        showAuthenticatedContainer();

        window.addEventListener("scroll", debounce(() => {
            if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 50) {
                fetchMorePosts();
            }
        }, 300));
    } else {
        showNotAuthenticatedContainer();
    }
}

function debounce(func, delay) {
    let timeoutId;
    return function (...args) {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => {
        func.apply(this, args);
      }, delay);
    };
}

let isLoading = false;
let page = 2;

async function fetchMorePosts() {
    if (isLoading) return;
    isLoading = true;
    console.log("Loading more data...");
    const response = await fetch(`/api/bunchOfPosts/${page}`);
    res = await response.json();
    if (!res.success) {
        showToast(res.message);
        return;
    }
    posts = res.data
    if(!posts || posts.length === 0) {
        console.log("No more data to load");
        return;
    }
    console.log("Page:", page, "Data count:", posts.length);
    page++;
    isLoading = false;

    fillPostsInHtml(posts, '', 'append');


}
