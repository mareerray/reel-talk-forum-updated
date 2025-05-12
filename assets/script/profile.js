document.addEventListener('DOMContentLoaded', function () {
    const tabButtons = document.querySelectorAll('.nav-link');
    const tabContents = document.querySelectorAll('.tab-pane');

    tabButtons.forEach(btn => {
        btn.addEventListener('click', function () {
        // Remove active from all
        tabButtons.forEach(b => b.classList.remove('active'));
        tabContents.forEach(c => c.classList.remove('active'));

        // Add active to clicked button and show correct tab
        this.classList.add('active');
        const tabId = this.getAttribute('data-tab');
        document.getElementById(tabId).classList.add('active');

        // If this is the profile tab, load the profile
        if (this.id === 'my-profile-button') {
            loadUserProfile();
        }
        });
    });
});


function loadUserProfile() {
    const elements = {
        nickname: document.getElementById('profile-nickname'),
        firstname: document.getElementById('profile-firstname'),
        lastname: document.getElementById('profile-lastname'),
        age: document.getElementById('profile-age'),
        gender: document.getElementById('profile-gender'),
        email: document.getElementById('profile-email')
    };
    
    Object.values(elements).forEach(el => {
        if (el) el.textContent = `${el.id.split('-')[1].charAt(0).toUpperCase() + el.id.split('-')[1].slice(1)}: Loading...`;
    });
    
    fetch('/api/user/profile')
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to load profile: ${response.status}`);
            }
            return response.json();
        })
        .then(user => {
            if (elements.nickname) elements.nickname.textContent = `Nickname: ${user.nickname}`;
            if (elements.firstname) elements.firstname.textContent = `First Name: ${user.first_name}`;
            if (elements.lastname) elements.lastname.textContent = `Last Name: ${user.last_name}`;
            if (elements.age) elements.age.textContent = `Age: ${user.age}`;
            if (elements.gender) elements.gender.textContent = `Gender: ${user.gender}`;
            if (elements.email) elements.email.textContent = `Email: ${user.email}`;
        })
        .catch(error => {
            Object.values(elements).forEach(el => {
                if (el) el.textContent = `${el.id.split('-')[1].charAt(0).toUpperCase() + el.id.split('-')[1].slice(1)}: Error loading`;
            });
        });
}
