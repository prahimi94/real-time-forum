document.querySelectorAll('[id^="flush-heading-"]').forEach(header => {
    const idPart = header.id.replace('flush-heading-', ''); // Extract the last part of the id
    const bodyId = `flush-collapse-${idPart}`; // Construct the body id

    header.addEventListener('click', () => {
      const body = document.getElementById(bodyId);
  
    // Optional: Close all others (accordion behavior)
    document.querySelectorAll('[id^="flush-collapse-"]').forEach(body => {
      if (body.id !== bodyId) {
        body.style.display = 'none';
      }
    });
  
    // Toggle current
    if (body.classList.contains('show')) {
        // Collapse it
        body.style.height = body.scrollHeight + 'px'; // set to current height
        requestAnimationFrame(() => {
          body.style.height = '0px';
          body.classList.remove('show');
        });
      } else {
        // Expand it
        body.style.height = body.scrollHeight + 'px';
        body.classList.add('show');
    
        // Remove fixed height after transition ends so it can adjust to dynamic content
        body.addEventListener('transitionend', function removeHeight() {
          body.style.height = 'auto';
          body.removeEventListener('transitionend', removeHeight);
        });
      }
    
    });

  });