// Show the "other" nationality text field if "Other" is selected
// document.getElementById('nationality').addEventListener('change', function() {
// 	if (this.value === 'other') {
// 		document.getElementById('otherNationality').style.display = 'inline-block';
// 	} else {
// 		document.getElementById('otherNationality').style.display = 'none';
// 	}
// });

function submitFormOnEnter(event) {
    if (event.keyCode === 13) { // Check if Enter key was pressed
      event.preventDefault(); // Prevent default behavior of form submission
      const form = event.target.closest('form'); // Find the closest form element
      const formData = new FormData(form); // Create a new FormData object with form data
      const fullName = formData.get('full-name'); // Get the value of the full-name input
      console.log('Affiant\'s Full Name:', fullName); // Log the value to the console
    }
  }
  
  
  document.addEventListener('DOMContentLoaded', function() {
    const inputElement = document.getElementById('full-name');
    inputElement.addEventListener('keydown', submitFormOnEnter);
  });
  