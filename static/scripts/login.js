document.getElementById('loginForm').addEventListener('submit', function (e) {
  e.preventDefault();

  const email = document.getElementById('email').value.trim();
  const password = document.getElementById('password').value.trim();
  const errorDiv = document.getElementById('error');

  if (!email || !password) {
    errorDiv.textContent = 'Email and password are required.';
    return;
  }

  // Simulate a failed login
  errorDiv.textContent = 'Login failed. Invalid credentials.';
});
