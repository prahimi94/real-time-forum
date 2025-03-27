let categories = [];
let posts = [];

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
    console.log('fetchMyCreatedPosts');
    const response = await fetch('/api/myCreatedPosts/');
    posts = await response.json();

    fillPostsInHtml(posts);
}

async function fetchMyLikedPosts() {
    console.log('fetchMyLikedPosts');
    const response = await fetch('/api/myLikedPosts/');
    posts = await response.json();

    fillPostsInHtml(posts);
}

function fillPostsInHtml(posts) {
    // load posts for home page
    const postsContainer = document.getElementById('postsContainer');
    postsContainer.innerHTML = "";
    
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
                                <span class="d-block post-dateTime">${post.created_at}</span>
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

addEventListener("DOMContentLoaded", function () {
    console.log('DOMContentLoaded');
    fetchCategories();
    fetchPosts();
});