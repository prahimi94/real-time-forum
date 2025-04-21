// document.querySelectorAll('[id^="flush-heading-"]').forEach(header => {
//     const idPart = header.id.replace('flush-heading-', ''); // Extract the last part of the id
//     const bodyId = `flush-collapse-${idPart}`; // Construct the body id

//     header.addEventListener('click', () => {
//       const body = document.getElementById(bodyId);
  
//     // Optional: Close all others (accordion behavior)
//     document.querySelectorAll('[id^="flush-collapse-"]').forEach(body => {
//       if (body.id !== bodyId) {
//         body.style.display = 'none';
//       }
//     });
  
//     // Toggle current
//     if (body.classList.contains('show')) {
//         // Collapse it
//         body.style.height = body.scrollHeight + 'px'; // set to current height
//         requestAnimationFrame(() => {
//           body.style.height = '0px';
//           body.classList.remove('show');
//         });
//       } else {
//         // Expand it
//         body.style.height = body.scrollHeight + 'px';
//         body.classList.add('show');
    
//         // Remove fixed height after transition ends so it can adjust to dynamic content
//         body.addEventListener('transitionend', function removeHeight() {
//           body.style.height = 'auto';
//           body.removeEventListener('transitionend', removeHeight);
//         });
//       }
    
//     });

//   });


function accordionHeadedClicked(idPart) {
  const body = document.getElementById(`flush-collapse-${idPart}`);
  const bodyId = `flush-collapse-${idPart}`;

  // Close all other accordion bodies
  document.querySelectorAll('[id^="flush-collapse-"]').forEach(otherBody => {
    if (otherBody.id !== bodyId) {
      body.style.display = 'none';
      otherBody.style.height = '0px';
      otherBody.classList.remove('show');
    }
  });

  // Toggle current
  if (body.classList.contains('show')) {
    // Collapse it
    body.style.height = body.scrollHeight + 'px'; // briefly set height to scrollHeight
    requestAnimationFrame(() => {
      body.style.display = 'none';
      body.style.height = '0px';
      body.classList.remove('show');
    });
  } else {
    // Ensure it's not hidden before measuring height
    body.style.display = 'block'; // in case it was hidden
    body.style.height = 'auto';
    // const height = body.scrollHeight + 'px';
    // body.style.height = '0px'; // reset height to 0 before expanding
    // body.offsetHeight; // force reflow
    // body.style.height = height;
    body.classList.add('show');

    body.addEventListener('transitionend', function removeHeight() {
      body.style.height = 'auto';
      body.removeEventListener('transitionend', removeHeight);
    });
  }
}
