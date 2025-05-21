document.addEventListener('DOMContentLoaded', function() {
    // Load posts when the page loads
    loadPosts();
});

let allPost = [];

function loadPosts() {
    fetch('/api/posts/list')
    .then(response => response.json())
    .then(posts => {
        allPosts = posts;
    displayPosts(posts, document.querySelector('#post-feed-content .post-grid'));
    });
}

function filterPostsByCategory(categoryName) {
    let filtered;
    if (categoryName === 'All') {
    filtered = allPosts;
} else {
    filtered = allPosts.filter(post =>
    post.categories.split(',').some(cat => cat.trim().includes(categoryName))
    );
}
displayPosts(filtered, document.querySelector('#post-feed-content .post-grid'));
}

const formattedDateTime = (() => {
    const date = new Date();
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    return `${day}-${month}-${year} ${hours}:${minutes}:${seconds}`;
})();

function displayPosts(posts, postContainer) {
    // Clear existing content
    postContainer.innerHTML = '';
    
    if (posts.length === 0) {
        postContainer.innerHTML = '<div class="col"><div class="post-container" style="grid-column: span 2; font-size: 1.2rem;">No posts yet.</div></div>';
        return;
    }
    
    // Add each post - USING CLASSES INSTEAD OF IDs
    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.className = 'post-container';
        postElement.innerHTML = `
            <h1 class="post-title" data-post-id="${post.id}" style="cursor: pointer;">${post.title}</h1>                
            <p class="post-info">Posted by: ${post.nickname}  on: ${formattedDateTime}</p>
            <p class="post-categories">Categories: ${post.categories}</p>
            <a href="#" class="read-button" data-post-id="${post.id}">Read More</a>
        `;
        postContainer.appendChild(postElement);
    });

    // Add event listeners to post titles and read more buttons - SCOPE TO POST FEED TAB
    document.querySelectorAll('#post-feed-content .post-title').forEach(title => {
        title.addEventListener('click', function() {
            const postId = this.getAttribute('data-post-id');
            const post = posts.find(p => p.id == postId);
            showSinglePost(postId, post);
        });
    });
    
    // Add event listeners to read more buttons - SCOPE TO POST FEED TAB
    document.querySelectorAll('#post-feed-content .read-button').forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const postId = this.getAttribute('data-post-id');
            const post = posts.find(p => p.id == postId);
            showSinglePost(postId, post);
        });
    });
}

function showSinglePost(postId, post) {
    // Get the post feed tab
    const postFeedButton = document.getElementById('post-feed-button');

    // Change the tab text to indicate we're viewing a single post
    postFeedButton.textContent = 'Post';

    // Find the container where posts are displayed (below the tab)
    const postContainer = document.querySelector('#post-feed-content .post-grid');

    // Update the container with the single post content
    postContainer.innerHTML = `
        <div class="single-post-container">
            <button id="back-to-posts" class="back-to-post-button">
                ← Back to Posts
            </button>
            <h1 class="post-title" data-post-id="${post.id}" style="cursor: pointer;">${post.title}</h1>                
            <p class="post-info">Posted by: ${post.nickname} on: ${formattedDateTime}</p>
            <p class="post-categories">Categories: ${post.categories}</p>
            <p class="post-content">${post.content}</p>
            
            <div class="comments-section">
                <h3>Comments</h3>
                <div class="comments-container">
                    <!-- Comments will be dynamically loaded here -->
                    <p>Loading comments...</p>
                </div>
                
                <div class="add-comment">
                    <h4>Add a Comment</h4>
                    <textarea 
                        class="form-control comment-text" 
                        rows="5" 
                        cols="40" 
                        maxlength="200" 
                        placeholder="Write your comment here..."></textarea>
                    <button 
                        class="submit-comment" 
                        data-post-id="${postId}">Submit Comment</button>
                </div>
            </div>
        </div>
    `;

    // Add event listener to back button
    document.getElementById('back-to-posts').addEventListener('click', function () {
        // Change the tab text back to 'Post Feed'
        postFeedButton.textContent = 'Post Feed';

        // Reload all posts
        loadPosts();
    });

    // Add event listener to comment submission
    const submitButton = document.querySelector('.submit-comment'); // Correctly select the button
    submitButton.addEventListener('click', function () {
        const commentText = document.querySelector('.comment-text').value.trim();
        if (!commentText) {
            alert("Comment cannot be empty!");
            return;
        }

        if (commentText.length > 200) {
            alert("Comment must be less than 200 characters.");
            return;
        }

        // Get the authentication token
        const token = localStorage.getItem('sessionToken');
        if (!token) {
            alert("You must be logged in to comment.");
            return;
        }
        
        // Ensure postId is a number
        const postId = parseInt(this.getAttribute('data-post-id'), 10);
        // Check for NaN and invalid postId
        if (isNaN(postId) || postId <= 0) {
            console.error("Invalid post ID:", this.getAttribute('data-post-id'));
            alert("Invalid post ID. Please try again.");
            return;
        }

        // Prepare request body
        const requestBody = {
            post_id: postId,
            content: commentText
        };
        
        fetch('/api/comments', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(requestBody)
        })
        
        .then(response => {
            if (response.ok) {
                document.querySelector('.comment-text').value = '';
                loadComments(postId); // Reload comments after successful submission
                return {status: "success"};
            } else {
                // Get the error message from the response
                return response.text().then(errorText => {
                    console.error(`Server error (${response.status}):`, errorText);
                    throw new Error(`${response.status} ${errorText}`);
                });
            }
        })
        .then(data => {
            if (data.status === "success") {
                console.log("Comment posted successfully");
            }
        })
        .catch(error => {
            console.error('Error:', error);
            if (error.message.includes('401')) {
                alert("Your session has expired. Please log in again.");
                window.location.href = '/login';
            } else {
                alert("Failed to submit comment: " + error.message);
            }        
        });
    });
    // Load comments for this post
    loadComments(postId);
}

function loadComments(postId) {
    const container = document.querySelector('.comments-container');
    container.innerHTML = '<p>Loading comments...</p>';

    fetch(`/api/comments/${postId}`)
    .then(response => {
        if (!response.ok) {
            // Parse JSON error message from backend
            return response.json().then(error => {
                throw new Error(error.error || "Failed to load comments");
            });
        }
        return response.json();
    })
    .then(comments => {
        comments = Array.isArray(comments) ? comments : [];
        container.innerHTML = comments.length > 0 
            ? comments.map(comment => `
                <div class="comment">
                    <p><strong>${comment.user_name}</strong> (${new Date(comment.created_at).toLocaleString()}):</p>
                    <p>${comment.content}</p>
                </div>
            `).join('')
            : '<p>No comments yet.</p>'; // Show this for empty arrays
    })
    .catch(error => {
        container.innerHTML = `<p>${error.message}</p>`; // Display backend error
        console.error("Comments error:", error);
    });
}


