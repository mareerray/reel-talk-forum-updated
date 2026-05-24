document.addEventListener('DOMContentLoaded', () => {
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const signUpFormContainer = document.getElementById('signUpFormContainer');
  const logInFormContainer = document.getElementById('logInFormContainer');

  const nickNameRadio = document.getElementById("nickNameField");
  const nickNameGroup = document.getElementById("nickNameGroup");
  const emailRadio = document.getElementById("emailField");
  const emailGroup = document.getElementById("emailGroup");

  const signUpForm = document.getElementById('signUpForm');
  const logInForm = document.getElementById('logInForm');

  function showSignIn() {
    if (signUpFormContainer) signUpFormContainer.style.display = 'none';
    if (logInFormContainer) logInFormContainer.style.display = 'flex';
    const switchToSignup = document.getElementById('switchToSignup');
    const switchToLogin = document.getElementById('switchToLogin');
    if (switchToSignup) switchToSignup.style.display = 'block';
    if (switchToLogin) switchToLogin.style.display = 'none';
  }

  function showSignUp() {
    if (logInFormContainer) logInFormContainer.style.display = 'none';
    if (signUpFormContainer) signUpFormContainer.style.display = 'flex';
    const switchToSignup = document.getElementById('switchToSignup');
    const switchToLogin = document.getElementById('switchToLogin');
    if (switchToSignup) switchToSignup.style.display = 'none';
    if (switchToLogin) switchToLogin.style.display = 'block';
  }

  if (signUpButton) signUpButton.addEventListener('click', showSignUp);
  if (signInButton) signInButton.addEventListener('click', showSignIn);

  function updateNavMenu() {
    const navMenu = document.getElementById("nav-menu");
    if (!navMenu) return;

    if (sessionStorage.getItem("sessionToken")) {
        const userNickname = sessionStorage.getItem("userNickname") || "User";
        navMenu.innerHTML = `
        <li><span class="user-display">User : ${userNickname}</span></li>
        <li><a href="#" id="logout-button" aria-label="Log out" title="Log out">
            <i class="material-icons" aria-hidden="true">logout</i>
        </a></li>
        `;

        const logoutBtn = document.getElementById("logout-button");
        if (logoutBtn) {
        logoutBtn.addEventListener("click", handleLogout);
        }
    } else {
        navMenu.innerHTML = `
        <li><a href="#" id="login-link" aria-label="Log in" title="Log in">
            <i class="material-icons" aria-hidden="true">login</i>
        </a></li>
        <li><a href="#" id="signup-link" aria-label="Sign up" title="Sign up">
            <i class="material-icons" aria-hidden="true">person_add</i>
        </a></li>
        `;

        const loginLink = document.getElementById("login-link");
        const signupLink = document.getElementById("signup-link");
        const container = document.getElementById("container");

        if (loginLink) {
            loginLink.addEventListener("click", function (event) {
                event.preventDefault();
                showSignIn();
            });
        }

        if (signupLink) {
            signupLink.addEventListener("click", function (event) {
                event.preventDefault();
                showSignUp();
            });
        }  
    }
  }


  function toggleFields() {
    if (nickNameRadio && nickNameRadio.checked) {
      if (nickNameGroup) nickNameGroup.style.display = "flex";
      if (emailGroup) emailGroup.style.display = "none";
    } else if (emailRadio && emailRadio.checked) {
      if (nickNameGroup) nickNameGroup.style.display = "none";
      if (emailGroup) emailGroup.style.display = "flex";
    }
  }

  if (nickNameRadio) nickNameRadio.addEventListener("change", toggleFields);
  if (emailRadio) emailRadio.addEventListener("change", toggleFields);

  if (signUpForm) {
    signUpForm.addEventListener('submit', function(e) {
        e.preventDefault();
        clearErrors();

        const firstName = document.getElementById('firstName').value.trim();
        const lastName = document.getElementById('lastName').value.trim();
        const nickName = document.getElementById('nickName').value.trim();
        const gender = document.querySelector('input[name="gender"]:checked')?.value;
        const age = document.getElementById('age').value.trim();
        const email = document.getElementById('email').value.trim();
        const password = document.getElementById('signUpPassword').value;
        const confirmPassword = document.getElementById('confirmPassword').value;

        let isValid = true;

        if (!firstName || !nickName || !gender || !email || !password || !confirmPassword) {
        displayError('general', 'All required fields must be filled out');
        isValid = false;
        }

        if (nickName.length < 5 || nickName.length > 15) {
        displayError('nickName', 'Nickname must be between 5 and 15 characters long');
        isValid = false;
        }

        const validNicknameRegex = /^[a-zA-Z0-9_-]+$/;
        if (!validNicknameRegex.test(nickName)) {
        displayError('nickName', 'Nickname can only contain letters, numbers, underscores, and dashes');
        isValid = false;
        }

        const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
        if (!emailRegex.test(email)) {
        displayError('email', 'Please enter a valid email address');
        isValid = false;
        }

        const ageValue = parseInt(age, 10);
        if (isNaN(ageValue) || ageValue < 13) {
        displayError('age', 'Age must be a number between 13 and 120');
        isValid = false;
        }

        if (password.length < 8) {
        displayError('password', 'Password must be at least 8 characters long');
        isValid = false;
        }

        if (!/[a-z]/.test(password)) {
        displayError('password', 'Password must contain at least one lowercase letter');
        isValid = false;
        }

        if (!/[A-Z]/.test(password)) {
        displayError('password', 'Password must contain at least one uppercase letter');
        isValid = false;
        }

        if (!/[0-9]/.test(password)) {
        displayError('password', 'Password must contain at least one digit');
        isValid = false;
        }

        if (!/[@$!%*?&]/.test(password)) {
        displayError('password', 'Password must contain at least one special character (@, $, !, %, *, ?, &)');
        isValid = false;
        }

        if (password !== confirmPassword) {
        displayError('confirmPassword', 'Passwords do not match');
        isValid = false;
        }

        if (!isValid) return;

        const userData = {
        firstName,
        lastName,
        username: nickName,
        gender,
        age,
        email,
        password
        };

        fetch('/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(userData)
        })
        .then(response => {
        if (!response.ok) {
            return response.json().then(err => {
            throw new Error(err.message || 'Registration failed');
            });
        }
        return response.json();
        })
        .then(data => {
        if (data.success) {
            signInButton?.click();
            signUpForm.reset();
        } else {
            displayError(data.field || 'general', data.message || 'Registration failed. Please try again.');
        }
        })
        .catch(error => {
        console.error('Error:', error);
        displayError('general',
            error.message.includes('Failed to fetch')
            ? 'Network error. Please check your connection and try again.'
            : 'An error occurred. Please try again later.'
        );
        });
    });
 }
    if (logInForm) {
    logInForm.addEventListener('submit', function(e) {
        e.preventDefault();
        clearErrors();

        const loginOption = document.querySelector('input[name="login-option"]:checked')?.value;
        let identifier, password;

        if (loginOption === 'nickName') {
        identifier = document.querySelector('#nickNameGroup input[name="nickName"]')?.value.trim();
        } else {
        identifier = document.querySelector('#emailGroup input[name="email"]')?.value.trim();
        }

        password = document.querySelector('#logInForm input[name="password"]')?.value;

        if (!identifier || !password) {
        displayError('login-general', 'Invalid credentials. Please try again.');
        return;
        }

        const loginData = { loginType: loginOption, identifier, password };

        fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(loginData)
        })
        .then(response => {
        if (!response.ok) {
            return response.json().then(err => {
            throw new Error(err.message || 'Invalid credentials. Please try again.');
            });
        }
        return response.json();
        })
        .then(data => {
        if (data.success) {
            sessionStorage.setItem('sessionToken', data.token);
            sessionStorage.setItem('userNickname', data.nickname || data.email);
            updateNavMenu();
            redirectToHomepage();
            if (typeof initializeApp === "function") initializeApp();
        } else {
            displayError('login-general', data.message || 'Invalid credentials. Please try again.');
        }
        })
        .catch(error => {
        console.error('Error:', error);
        displayError('login-general',
            error.message.includes('Failed to fetch')
            ? 'Network error. Please check your connection and try again.'
            : 'Invalid credentials. Please try again.'
        );
        });
    });
    }
  showSignIn();
  updateNavMenu();
});


// Function to redirect to homepage
function redirectToHomepage() {
    // Hide the auth view
    document.getElementById('authView').style.display = 'none';
    // Show the main view
    document.getElementById('mainView').style.display = 'block';
    // Remove auth-view class from body
    document.body.classList.remove('auth-view');
    // Fetch user profile data
    fetchUserProfile();
    // Load other necessary data for the main view
    // loadCategories();
    loadPosts();
}

// Function to fetch user profile data
function fetchUserProfile() {
    fetch('/api/user/profile')
    .then(response => {
        if (!response.ok) {
            throw new Error('Invalid token');
        }
        return response.json();
    })
    .then(data => {
        sessionStorage.setItem("userNickname", data.nickname);
        updateNavMenu();
        // Update the profile information in the UI
        document.getElementById('profile-nickname').textContent = 'Nickname: ' + data.nickname;
        document.getElementById('profile-firstname').textContent = 'First Name: ' + data.first_name;
        document.getElementById('profile-lastname').textContent = 'Last Name: ' + data.last_name;
        document.getElementById('profile-age').textContent = 'Age: ' + data.age;
        document.getElementById('profile-gender').textContent = 'Gender: ' + data.gender;
        document.getElementById('profile-email').textContent = 'Email: ' + data.email;
        
        // Update the greeting in the header
        document.querySelector('.page-title').textContent = 'Welcome, ' + data.nickname + '!';
    })
    .catch(error => {
        console.error('Error fetching profile:', error);
        // Clear invalid token and redirect to login
        sessionStorage.removeItem("sessionToken");
        document.getElementById("authView").style.display = "block";
        document.getElementById("mainView").style.display = "none";
        document.body.classList.add('auth-view');
    });
}

document.addEventListener("DOMContentLoaded", () => {
    const sessionToken = sessionStorage.getItem("sessionToken");

    if (sessionToken) {
        document.getElementById("authView").style.display = "none";
        document.getElementById("mainView").style.display = "block";
        document.body.classList.remove('auth-view');
        fetchUserProfile();
        // loadCategories();
        // loadPosts();
        if (typeof initializeApp === "function") initializeApp();
    } else {
        document.getElementById("authView").style.display = "block";
        document.getElementById("mainView").style.display = "none";
        document.body.classList.add('auth-view');
    }

    updateNavMenu(); // This ensures nav menu updates on page reload
});


function handleLogout(event) {
    event.preventDefault();

    if (typeof socket !== 'undefined' && socket) {
        socket.close();
    }

    // Clear client-side storage
    sessionStorage.clear();

    // Notify server
    fetch('/api/logout', { method: 'POST', credentials: 'include' })
        .finally(() => {
            window.location.reload(); // <-- This ensures a full reset
        });
}

// Helper functions for displaying errors and success messages
function displayError(field, message) {
    // Create error element if it doesn't exist
    let errorElement = document.getElementById(`${field}-error`);
    
    if (!errorElement) {
        errorElement = document.createElement('div');
        errorElement.id = `${field}-error`;
        errorElement.className = 'error-message';
        
        // Find the appropriate input field to place the error after
        const inputField = document.getElementById(field) || 
                        document.querySelector(`[name="${field}"]`);
        
        if (inputField) {
            // For radio buttons, find the parent container
            if (inputField.type === 'radio') {
                const radioGroup = inputField.closest('.radio-group');
                if (radioGroup) {
                    radioGroup.appendChild(errorElement);
                }
            } else {
                // For other inputs, find the parent input-group
                const inputGroup = inputField.closest('.input-group');
                if (inputGroup) {
                    inputGroup.appendChild(errorElement);
                } else {
                    // Fallback for general errors
                    if (field === 'general') {
                        signUpForm.prepend(errorElement);
                    } else if (field === 'login-general') {
                        logInForm.prepend(errorElement);
                    }
                }
            }
        } else if (field === 'general') {
            signUpForm.prepend(errorElement);
        } else if (field === 'login-general') {
            logInForm.prepend(errorElement);
        }
    }
    
    errorElement.textContent = message;
    errorElement.style.color = 'red';
    errorElement.style.fontSize = '0.8rem';
    errorElement.style.marginTop = '5px';
}

function clearErrors() {
    const errorMessages = document.querySelectorAll('.error-message');
    errorMessages.forEach(error => error.remove());
}

