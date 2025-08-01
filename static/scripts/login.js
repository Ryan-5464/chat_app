document.getElementById('login-btn').addEventListener('click', function (e) {
    e.preventDefault();

    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();

    attemptLogin(email, password)
});

function attemptLogin(email, password) {
    fetch(BASEURL + '/api/login', loginRequestBody(email, password))
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
            window.location.href = '/chat';
        }

        if (responsePayload.NoError == false) {
            const errorElement = document.getElementById('error')
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
