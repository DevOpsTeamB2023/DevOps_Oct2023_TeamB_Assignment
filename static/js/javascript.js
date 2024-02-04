function signup(){
    var request = new XMLHttpRequest();
    const form = document.getElementById('signupForm');

    const curl = 'http://localhost:5001/api/v1/accounts';

    const username = form.elements['signup_username'].value;
    const password = form.elements['signup_password'].value;
    console.log(username);
    console.log(password);

    request.open("POST", curl);
    request.send(JSON.stringify({
        "username": username,
        "password": password, 
        "accType": "User", 
        "accStatus": "Pending"
        
    }));
    form.reset();
    alert("Account request sent. Please wait for admin approval.");
    return false //prevent default submission
}

function login(){
    var request = new XMLHttpRequest();
    const form = document.getElementById('loginForm');
    const username = form.elements['login_username'].value;
    const password = form.elements['login_password'].value;
    console.log(username);
    console.log(password);

    const curl = 'http://localhost:5001/api/v1/accounts?username=' + encodeURIComponent(username) + '&password=' + encodeURIComponent(password);
    console.log(curl);

    request.open("GET", curl);
    request.onreadystatechange = function() {
      if (request.readyState === 4) {
        if (request.status === 200) {
          // Successful login, redirect to main page
          location.href = "/static/templates/user_details.html";
        } else if (request.status === 401) {
          // Login failed, handle error
          form.reset();
          document.getElementById('error-message').innerHTML = 'Incorrect Phone Number or Password.';
        } else {
          // Handle other status codes or network errors
          document.getElementById('error-message').innerHTML = 'An error occurred. Please try again later.';
        }
      }
    };
    request.send();
    return false
}

function listUsers() {
  // Make a GET request to the server endpoint
  const url = `http://localhost:5001/api/v1/accounts/all`;
  fetch(url)
    .then(response => {
      if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then(data => {
      console.log("Data from server:", data);

      // Get the table body element
      var tableBody = document.getElementById('user_details_table').getElementsByTagName('tbody')[0];

      // Clear existing rows
      tableBody.innerHTML = '';

      // Iterate through the received data and append rows to the table
      data.forEach(user => {
        var row = tableBody.insertRow();
        row.innerHTML = `<td>${user.accId}</td>
                        <td>${user.username}</td>
                        <td>${user.accType}</td>
                        <td>
                          ${user.accStatus}
                          ${user.accStatus === 'Pending' ? '<button class="btn btn-outline-secondary" onclick="return approveUser(' + user.accId + ')">approve</button>' : ''}
                        </td>
                        <td>
                          <button class="btn btn-outline-secondary" onclick="return modifyUser(${user.accId})">modify</button>
                          <button class="btn btn-outline-secondary" onclick="return deleteUser(${user.accId})">delete</button>
                        </td>`;
      });
    })
    .catch(error => console.error('Error fetching user details:', error));
}

// Function to delete a user (replace this with your actual delete logic)
function deleteUser(userId) {
  console.log('Deleting user with ID:', userId);
  const url = `http://localhost:5001/api/v1/accounts/delete?accID=${userId}`;

  // Confirm deletion with the user (you can customize this)
  if (confirm("Are you sure you want to delete this user?")) {
  //     // Make a DELETE request to the server endpoint
      fetch(url, {
        method: 'DELETE',
      })
      .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.text(); // assuming the server returns a text response
      })
      .then(data => {
        console.log('Server response:', data);
        
        // Optionally, you can call listUsers() again to refresh the user list
        listUsers();
      })
      .catch(error => console.error('Error deleting user:', error));
  }
}


// Function to modify a user (replace this with your actual modify logic)
function modifyUser(userId) {
  console.log('Modifying user with ID:', userId);
  const url = `http://localhost:5001/api/v1/accounts/get?accID=${userId}`
  // Fetch user details by userId
  fetch(url)
  .then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.json();
  })
  .then(user => {
    localStorage.setItem('modifyUserData', JSON.stringify(user));
    location.href = "/static/templates/modify_user.html";
  })

}

// Update user details to DB
async function submitModification() {
  event.preventDefault();

  const form = document.getElementById('modifyForm');
  const accID = document.getElementById('modify_accID').innerHTML;

  const url = `http://localhost:5001/api/v1/accounts/${accID}`;

  const username = form.elements['modify_username'].value;
  const accType = document.getElementById('modify_user').checked ? 'User' : 'Admin';

  try {
    const response = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        "username": username,
        "accType": accType,
        //"accStatus": accStatus
      }),
    });

    if (response.ok) {
      console.log("Update successful");
      window.location.href = "/static/templates/user_details.html";
    } else {
      const errorText = await response.text();
      alert("Error updating the Account Details. Status: " + response.status + "\n" + errorText);
    }
  } catch (error) {
    console.error("Error updating the Account Details:", error);
    alert("An unexpected error occurred. Please try again later.");
  }

  // Prevent default form submission
  // return false;
}


async function approveUser(userId){
  event.preventDefault();

  const url = `http://localhost:5001/api/v1/accounts/approve?accID=${userId}`

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        "accStatus": "Created"
      }),
    });

    if (response.ok) {
      console.log("Update successful");
      window.location.href = "/static/templates/user_details.html";
    } else {
      const errorText = await response.text();
      alert("Error approving the Account. Status: " + response.status + "\n" + errorText);
    }
  } catch (error) {
    console.error("Error approving the Account: ", error);
    alert("An unexpected error occurred. Please try again later.");
  }

}
