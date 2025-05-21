document.addEventListener('DOMContentLoaded', function() {
    // Check if user is logged in
    const token = localStorage.getItem('sessionToken');

    // Load categories when the page loads
    loadCategories();
    
    // Add event listener to the form
    const form = document.getElementById('create-post-form');
    if (form) {
        form.addEventListener('submit', createPost);
    } else {
        console.error("Form element 'create-post-form' not found");
    }
});

// Function to load all available categories
function loadCategories() {
    const token = localStorage.getItem('sessionToken');
    
    fetch('/api/categories', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to load categories');
        }
        return response.json();
    })
    .then(categories => {
        const container = document.getElementById('categories-container');
        container.innerHTML = '';
        
        categories.forEach(category => {
            const div = document.createElement('div');
            div.className = 'category-checkbox';
            
            const checkbox = document.createElement('input');
            checkbox.type = 'checkbox'; // Change from 'radio' to 'checkbox'
            checkbox.id = `category-${category.id}`;
            checkbox.name = 'categories';
            checkbox.value = category.name;
            checkbox.dataset.emoji = category.emoji;
                        
            const label = document.createElement('label');
            label.htmlFor = `category-${category.id}`;
            label.textContent = `${category.name} ${category.emoji}`;
            label.title = `${category.name} ${category.emoji}`;
            
            div.appendChild(checkbox);
            div.appendChild(label);
            container.appendChild(div);
        });

        // Add event listeners to checkboxes to enforce max 3 categories
        const checkboxes = document.querySelectorAll('input[name="categories"]');
        checkboxes.forEach(checkbox => {
            checkbox.addEventListener('change', function() {
                const checked = document.querySelectorAll('input[name="categories"]:checked');
                const errorElement = document.getElementById('categories-error');
                
                if (checked.length > 3) {
                    this.checked = false;
                    if (errorElement) {
                        errorElement.textContent = 'You can select maximum 3 categories';
                    }
                } else if (errorElement) {
                    errorElement.textContent = '';
                }
            });
        });
    })
    
    .catch(error => {
        console.error('Error loading categories:', error);
        const container = document.getElementById('categories-container');
        if (container) {
            container.innerHTML = '<p class="error">Failed to load categories. Please try again later.</p>';
        }
    });
}

// Function to create a new post
function createPost(event) {
    event.preventDefault();
    
    const titleElement = document.getElementById('input-post-title');
    const contentElement = document.getElementById('input-post-content');
    
    if (!titleElement || !contentElement) {
        console.error("Title or content input elements not found");
        alert("Form elements not found. Please refresh the page and try again.");
        return;
    }
    
    const title = titleElement.value.trim();
    const content = contentElement.value.trim();
    const checkedCategories = document.querySelectorAll('input[name="categories"]:checked');
    const errorElement = document.getElementById('categories-error');
    
    // Validate input
    if (!title) {
        alert('Please enter a title');
        return;
    }
    
    if (!content) {
        alert('Please enter content');
        return;
    }
    
    if (checkedCategories.length < 1) {
        if (errorElement) {
            errorElement.textContent = 'Please select at least one category';
        }
        return;
    }
    
    if (checkedCategories.length > 3) {
        if (errorElement) {
            errorElement.textContent = 'You can select maximum 3 categories';
        }
        return;
    }
    
    // Format categories as they are stored in the database: "Category1 Emoji1,Category2 Emoji2"
    const formattedCategories = Array.from(checkedCategories).map(checkbox => 
        `${checkbox.value} ${checkbox.dataset.emoji}`
    ).join(',');
        
    // Create request body
    const requestBody = {
        title: title,
        content: content,
        categories: formattedCategories
    };
        
    // Get auth token
    const token = localStorage.getItem('sessionToken');
    
    // Send request to create post
    fetch('/api/posts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(requestBody)
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => {
                throw new Error(text || 'Failed to create post');
            });
        }
        return response.json();
    })
    .then(data => {
        // alert('Post created successfully!');
        window.location.href = `/post/${data.post_id}`;
    })
    .catch(error => {
        console.error('Error creating post:', error);
        // alert(`Error: ${error.message}`);
    });
}