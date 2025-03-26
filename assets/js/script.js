let categories = [];
let posts = [];

async function fetchCategories() {
    const response = await fetch('/api/categories/');
    categories = await response.json();

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
}

async function fetchPosts() {
    const response = await fetch('/api/posts/');
    posts = await response.json();

    const postsContainer = document.getElementById('postsContainer');
    postsContainer.innerHTML = "";
    posts.forEach(post => {
        const postDivTag = document.createElement('div');
        postDivTag.setAttribute('class', 'col-sm-12 col-md-12 mb-3');
        postDivTag.innerHTML = `<div class="post-card">
        <h5 class="mt-2 post-title"><a href="/post/${post.uuid}" class="nav-link" >${post.title}</a></h5>
        <div>
            <div class="d-flex text-body-secondary pt-3 m-posts">
                <div style="padding: 7px;"><i class="fa-solid fa-user" style="font-size: 2rem;"></i></div>

                <div class="pb-3 mb-0 small lh-sm w-100 mb-3 ms-2 mt-1">
                    <div class="d-flex justify-content-between mb-1 m-posts-userInfo">
                        <span class="post-user">${post.user.username}</span>
                        <span class="text-right m-posts-ctg">
                            
                                <span class="badge-p text-dark"><a href="/posts/Name">.Name</a></span>
                            
                        </span>
                    </div>
                    <span class="d-block post-dateTime">${post.created_at}</span>
                </div>
            </div>
        </div>
        <p class="post-description">${post.description}</p>
    </div>`
        postsContainer.appendChild(postDivTag);
    });
}


addEventListener("DOMContentLoaded", function () {
    console.log('DOMContentLoaded');
    fetchCategories();
    fetchPosts();
});