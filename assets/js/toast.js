function showMyToast() {
  // Get the snackbar DIV
  var toast = document.getElementById("liveToast");

  // Add the "show" class to DIV
  toast.classList.add("show");

  // close toast by clicking
  toast.onclick = function () {
    toast.classList.remove("show");
  };

  // After 3 seconds, remove the show class from DIV
  setTimeout(function () { toast.classList.remove("show"); }, 2000);
}