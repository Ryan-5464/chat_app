document.getElementById('login-btn').addEventListener('click', function (e) {
    e.preventDefault();

    const email = document.getElementById(APP.ID.GEN.EMAIL).value.trim();
    const password = document.getElementById(APP.ID.GEN.PASSWORD).value.trim();

    attemptLogin(email, password)
});

function attemptLogin(email, password) {
    fetch(APP.URL.BASE.concat(APP.ENDPOINT.LOGIN), loginRequestBody(email, password))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[AttemptLogin]Received:', responsePayload);
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

function loginRequestBody(email, password) {
    return {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            Email: email,
            Password: password,
        })
    }
}

