const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

signUpButton.addEventListener('click', () => {
	container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
	container.classList.remove("right-panel-active");
});

const nickNameRadio = document.getElementById("nickNameField");
const nickNameGroup = document.getElementById("nickNameGroup");
const emailRadio = document.getElementById("emailField");
const emailGroup = document.getElementById("emailGroup");
nickNameRadio.addEventListener("change", toggleFields);
emailRadio.addEventListener("change", toggleFields);

function toggleFields() {
	if (nickNameRadio.checked) {
		nickNameGroup.style.display = "flex";
		emailGroup.style.display = "none";
	} else if (emailRadio.checked) {
		nickNameGroup.style.display = "none";
		emailGroup.style.display = "flex";
	}
}
document.querySelectorAll('.toggle-password').forEach(icon => {
    icon.addEventListener('click', () => {
        const targetId = icon.getAttribute('data-target');
        const input = document.getElementById(targetId);

        if (input.type === 'password') {
        input.type = 'text';
        icon.textContent = 'visibility_off';
    } else {
        input.type = 'password';
        icon.textContent = 'visibility';
    }
    });
});
// Get references to the forms
const signUpForm = document.getElementById('signUpForm');
const logInForm = document.getElementById('logInForm');

// Sign Up Form Submission
signUpForm.addEventListener('submit', function(e) {
    e.preventDefault();

    // Clear previous error messages
    clearErrors();

    // Get form values
    const firstName = document.getElementById('firstName').value.trim();
    const lastName = document.getElementById('lastName').value.trim();
    const nickName = document.getElementById('nickName').value.trim();
    const gender = document.querySelector('input[name="gender"]:checked')?.value;
    const age = document.getElementById('age').value.trim();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('signUpPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    // Client-side validation
    let isValid = true;

    // Basic required field validation
    if (!firstName || !nickName || !gender || !email || !password || !confirmPassword) {
        displayError('general', 'All required fields must be filled out');
        isValid = false;
    }
    
    // Nickname validation (matches your Go validation)
    if (nickName.length < 5 || nickName.length > 15) {
        displayError('nickName', 'Nickname must be between 5 and 15 characters long');
        isValid = false;
    }

    // Check if nickname contains only valid characters
    const validNicknameRegex = /^[a-zA-Z0-9_-]+$/;
    if (!validNicknameRegex.test(nickName)) {
        displayError('nickName', 'Nickname can only contain letters, numbers, underscores, and dashes');
        isValid = false;
    }
    
    // Email validation
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(email)) {
        displayError('email', 'Please enter a valid email address');
        isValid = false;
    }


    // Age validation
    const ageValue = parseInt(age, 10);
    if (isNaN(ageValue) || ageValue < 13) {
        displayError('age', 'Age must be a number between 13 and 120');
        isValid = false;
    }
    // Password validation
    if (password.length < 8) {
        displayError('password', 'Password must be at least 8 characters long');
        isValid = false;
    }
    
    // Check for lowercase, uppercase, digit, and special character
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
    
    // Confirm password
    if (password !== confirmPassword) {
        displayError('confirmPassword', 'Passwords do not match');
        isValid = false;
    }
    
    if (isValid) {
        // Prepare data for sending to server
        const userData = {
            firstName: firstName,
            lastName: lastName,
            username: nickName,
            gender: gender,
            age: age,
            email: email,
            password: password
        };

        // Send data to server
        fetch('/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData)
        })
        .then((response) => {
            if (!response.ok) {
                return response.json().then((err) => {
                    throw new Error(err.message || 'Registration failed');
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showSuccessMessage('Registration successful! You can now log in.');
                // Switch to login view
                signInButton.click();
                signUpForm.reset();
            } else {
                if (data.field && data.message) {
                    displayError(data.field, data.message);
                } else {
                    displayError('general', data.message || 'Registration failed. Please try again.');
                }
            }
        })
        .catch((error) => {
            console.error('Error:', error);
            displayError(
                'general',
                error.message.includes('Failed to fetch')
                    ? 'Network error. Please check your connection and try again.'
                    : 'An error occurred. Please try again later.'
            );
        });
    }
});


// Login Form Submission
logInForm.addEventListener('submit', function(e) {
    e.preventDefault();
    
    // Clear previous error messages
    clearErrors();
    
    // Get login option (nickname or email)
    const loginOption = document.querySelector('input[name="login-option"]:checked').value;
    
    // Get form values based on login option
    let identifier, password;
    
    if (loginOption === 'nickName') {
        identifier = document.querySelector('#nickNameGroup input[name="nickName"]').value.trim();
    } else {
        identifier = document.querySelector('#emailGroup input[name="email"]').value.trim();
    }
    
    password = document.querySelector('#logInForm input[name="password"]').value;
    
    // Basic validation
    if (!identifier || !password) {
        displayError('login-general', 'Invalid credentials. Please try again.');
        return;
    }
    
    // Prepare login data
    const loginData = {
        loginType: loginOption,
        identifier: identifier,
        password: password
    };
    
    // Send login request to server
    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(loginData)
    })
    .then(response => {
        if(!response.ok) {
            return response.json().then(err => {
                throw new Error(err.message || 'Invalid credentials. Please try again.');
            });
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            // Store the session token in localStorage
            localStorage.setItem('sessionToken', data.token); 
    
            // Update the navigation menu
            updateNavMenu(); // Add this line
    
            // Redirect to homepage
            redirectToHomepage();      
            
            // Initialize WebSocket connection
            if (typeof initializeApp === "function") initializeApp();
        } else {
            displayError('login-general', data.message || 'Invalid credentials. Please try again.');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        displayError(
            'login-general',
            error.message.includes('Failed to fetch')
                ? 'Network error. Please check your connection and try again.'
                : 'Invalid credentials. Please try again.'
        );
    });
});

// Function to redirect to homepage
function redirectToHomepage() {
    // Hide the auth view
    document.getElementById('authView').style.display = 'none';
    // Show the main view
    document.getElementById('mainView').style.display = 'block';

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
        localStorage.removeItem("sessionToken");
        document.getElementById("authView").style.display = "block";
        document.getElementById("mainView").style.display = "none";
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

function showSuccessMessage(message) {
    const successMessage = document.createElement('div');
    successMessage.className = 'success-message';
    successMessage.textContent = message;
    successMessage.style.color = 'green';
    successMessage.style.fontSize = '0.9rem';
    successMessage.style.marginTop = '10px';
    successMessage.style.textAlign = 'center';
    
    // Add to the active form
    if (container.classList.contains('right-panel-active')) {
        signUpForm.appendChild(successMessage);
    } else {
        logInForm.appendChild(successMessage);
    }
    
    // Remove after 3 seconds
    setTimeout(() => {
        successMessage.remove();
    }, 3000);
}

document.addEventListener("DOMContentLoaded", () => {
    updateNavMenu(); // Populate navigation menu based on user state
});

// Uncomment this in auth.js
document.addEventListener("DOMContentLoaded", () => {
    const sessionToken = localStorage.getItem("sessionToken");

    if (sessionToken) {
        document.getElementById("authView").style.display = "none";
        document.getElementById("mainView").style.display = "block";
        fetchUserProfile();
        // loadCategories();
        // loadPosts();
        if (typeof initializeApp === "function") initializeApp();
    } else {
        document.getElementById("authView").style.display = "block";
        document.getElementById("mainView").style.display = "none";
    }

    updateNavMenu(); // This ensures nav menu updates on page reload
});


function updateNavMenu() {
    const navMenu = document.getElementById("nav-menu");

    if (localStorage.getItem("sessionToken")) {
        navMenu.innerHTML = `
            <li><a href="#" id="logout-button">[Logout]</a></li>
        `;

        const logoutBtn = document.getElementById("logout-button");
        logoutBtn.addEventListener("click", handleLogout);
    } else {
        navMenu.innerHTML = `
            <li><a href="#" id="login-link">[Login]</a></li>
            <li><a href="#" id="signup-link">[Sign Up]</a></li>
        `;

        const loginLink = document.getElementById("login-link");
        const signupLink = document.getElementById("signup-link");
        const container = document.getElementById('container');

        if (loginLink) {
            loginLink.addEventListener("click", function(event) {
                event.preventDefault();
                container.classList.remove("right-panel-active");
            });
        }
        if (signupLink) {
            signupLink.addEventListener("click", function(event) {
                event.preventDefault();
                container.classList.add("right-panel-active");
            });
        }
    }
}

/* function handleLogout(event) {
    event.preventDefault();

    // Clear client-side storage
    localStorage.clear();
    sessionStorage.clear();

    // Notify server
    fetch('/api/logout', { method: 'POST', credentials: 'include' })
        .then(() => {
            document.getElementById("authView").style.display = "block";
            document.getElementById("mainView").style.display = "none";
            updateNavMenu(); // Refresh navigation
        })
        .catch(error => console.error('Logout error:', error));
} */

        function handleLogout(event) {
    event.preventDefault();

    // Clear client-side storage
    localStorage.clear();
    sessionStorage.clear();

    // Notify server
    fetch('/api/logout', { method: 'POST', credentials: 'include' })
        .finally(() => {
            window.location.reload(); // <-- This ensures a full reset
        });
}
