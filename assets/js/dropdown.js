document.querySelectorAll(".mydropdown-toggle").forEach((toggle) => {
    toggle.addEventListener("click", function (e) {
        console.log("clicked");
      e.stopPropagation();
      const dropdown = this.closest(".mydropdown");
      
      // Close all other dropdowns
      document.querySelectorAll(".mydropdown").forEach((d) => {
        if (d !== dropdown) d.classList.remove("show");
      });
  
      dropdown.classList.toggle("show");
    });
  });
  
  // Close dropdown if clicked outside
  document.addEventListener("click", function (e) {
    // Only close if clicked outside any .mydropdown
    if (!e.target.closest(".mydropdown")) {
      document.querySelectorAll(".mydropdown").forEach((dropdown) => {
        dropdown.classList.remove("show");
      });
    }
  });
  