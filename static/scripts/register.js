document.getElementById("register-btn").addEventListener("click", function (e) {
  e.preventDefault(); 

  const username = document.getElementById(APP.ID.GEN.PFL_USERNAME).value.trim();
  const password = document.getElementById(APP.ID.GEN.PASSWORD).value.trim();
  const email = document.getElementById(APP.ID.GEN.EMAIL).value.trim();

  attemptRegister(email, password, username);
});

function attemptRegister(email, password, username) {
    fetch(APP.URL.BASE.concat(APP.ENDPOINT.REGISTER), registerRequestBody(email, password, username))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[AttemptRegister]Received:', responsePayload);
        console.log(responsePayload.NoError)
        if (responsePayload.NoError ==  true) {
            window.location.href = APP.ENDPOINT.CHAT;
        }

        if (responsePayload.NoError == false) {
            const errorElement = document.getElementById(APP.ID.GEN.ERROR)
            errorElement.textContent = responsePayload.ErrorMessage
        }
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}

function registerRequestBody(email, password, username) {
    return {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            Email: email,
            Password: password,
            Name: username,
        })
    }
}
