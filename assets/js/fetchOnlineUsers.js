async function fetchOnlineUsers() {
    try {
        const response = await fetch('/api/online-users');
        const data = await response.json();

        const onlineUsersContainer = document.getElementById('online-users');
        onlineUsersContainer.innerHTML = ''; // Clear existing content

        if (data.message) {
            // Display the "No users are online now" message
            onlineUsersContainer.textContent = data.message;
        } else {
            // Display the list of online users
            data.forEach(user => {
                const userElement = document.createElement('div');
                userElement.className = 'online-user';
                userElement.innerHTML = `
                    <img src="${user.profile_photo || '/default-avatar.png'}" alt="${user.username}" class="user-avatar">
                    <span class="user-name">${user.name || user.username}</span>
                `;
                onlineUsersContainer.appendChild(userElement);
            });
        }
    } catch (error) {
        console.error('Error fetching online users:', error);
    }
}

// Call the function periodically to refresh the list
setInterval(fetchOnlineUsers, 5000); // Refresh every 5 seconds