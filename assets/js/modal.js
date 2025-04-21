// Get the modal
var modal = document.getElementById("myModal");

// Get the button that opens the modal
var btn = document.getElementById("myBtn");

// Get the <span> element that closes the modal
var span = document.getElementsByClassName("close")[0];

document.addEventListener("click", function (e) {
    const trigger = e.target.closest(".open-modal");
    if (trigger) {
      e.preventDefault();
      const modalId = trigger.getAttribute("data-target");
      const modal = document.getElementById(modalId);
      
      if (modal) {
        modal.classList.add("show");
        modal.style.display = "block";

        console.log('opened')
        if (modal.id.startsWith('updatePostModal-')) {
          const postId = modal.id.split('-')[1]; // Extract the post ID from the modal ID
          
          // Fetch and populate the form with post data here
          const postElement = document.getElementById(`post-${postId}`);
          const title = postElement.querySelector('.post-title').textContent.trim();
          const description = postElement.querySelector('.post-description').textContent.trim();
          const selectedCategories = Array.from(postElement.querySelectorAll('.m-posts-ctg a')).map(category => {
              const categoryName = category.textContent.trim();
              const categoryObj = categories.find(cat => cat.name === categoryName);
              return categoryObj ? categoryObj.id : null;
          }).filter(id => id !== null);

          const form = document.getElementById(`updatePostForm-${postId}`);
          form.querySelector('input[name="title"]').value = title;
          form.querySelector('textarea[name="description"]').value = description;

          if (selectedCategories && selectedCategories.length > 0) {
            const checkboxes = form.querySelectorAll('input[name="update_post_categories"]');
            checkboxes.forEach((checkbox) => {
              const isChecked = selectedCategories.includes(parseInt(checkbox.value));
              checkbox.checked = isChecked;

              if (isChecked) {
                // Manually trigger the change event
                // Use { bubbles: true } if any listeners are attached higher up the DOM (like on a form or container).
                checkbox.dispatchEvent(new Event('change', { bubbles: true }));
              }
            });

          }
        }
      }
    }
  });

document.addEventListener("click", function (e) {
    if (e.target.classList.contains("close-modal")) {
      const modal = e.target.closest(".modal");
      if (modal) {
        modal.classList.remove("show");
        modal.style.display = "none";
      }
    }
  });

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}