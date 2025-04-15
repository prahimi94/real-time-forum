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
      }
    }
  });

//todo
// When the user clicks on the button, open the modal
// btn.onclick = function() {
//   modal.style.display = "block";
// }

// When the user clicks on <span> (x), close the modal
// span.onclick = function() {
//   modal.style.display = "none";
// }

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