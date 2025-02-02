document.addEventListener('DOMContentLoaded', function () {
    new ShowHomePage();

    document.body.addEventListener('click', function (event) {
        if (event.target.matches('#showRegister')) {
            new RegisterForm();
        } else if (event.target.matches('#showLogin')) {
            new LoginForm();
        } else if (event.target.matches('#showHome')) {
            new ShowHomePage();
        }
    });
});

class ShowHomePage {
    constructor() {
        this.element = document.createElement('div');
        this.element.className = 'home';
        document.getElementById('formContainer').appendChild(this.element);
        this.showHomePage();
    }

    showHomePage() {
        this.element.innerHTML = `
            <h1>WELCOME TO REAL-TIME-FORUM</h1>
            <p class="subtext">Join our community and explore new ideas</p>
            <p class="subtext">Connect, share, and learn with fellow enthusiasts</p>
            <div class="buttons">
                <button class="btn" id="showLogin">Sign In</button>
                <button class="btn" id="showRegister">Sign Up</button>
                <button class="btn" onclick="window.location.href = '/guest'">Guest</button>
            </div>
        `;
    }
}

class LoginForm {
    constructor() {
        this.render();
    }

    render() {
        const formContainer = document.getElementById('formContainer');
        formContainer.innerHTML = `
            <div class="login-form">
                <h1>Login</h1>
                <div class="container">
                    <div class="main">
                        <div class="content">
                            <h2>Log In</h2>
                            <form id="loginForm">
                                <input type="text" name = "NicknameOREmail"  placeholder="Nickname or Email" required autofocus>
                                <input type="password" name="password" placeholder="Password" required>
                                <button class="btn" type="submit">Login</button>
                            </form>
                            <p class="account">Don't Have An Account? <a href="#" id="showRegister">Register</a></p>
                        </div>
                        <div class="form-img">
                            <img src="../styles/bg.png" alt="">
                        </div>
                    </div>
                </div>
            </div>
        `;

        document.getElementById('loginForm').addEventListener('submit', this.handleSubmit.bind(this));
    }


    async handleSubmit(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = Object.fromEntries(formData);
        let user = {}
        const isEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.NicknameOREmail);
        if (isEmail) {
            user = {
                email: data.NicknameOREmail,
                password : data.password,
            }
        }else {
            user = {
              nickname : data.NicknameOREmail,
              password : data.password,
            }
        }


        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(user),
            });

            const result =await response.json()
            if (!response.ok) {
                throw new Error(result.message || 'login failed');
            }
            

            alert('Login successful!');
            sessionStorage.setItem('username', result.username); 
            new ForumPage();
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }
}

class RegisterForm {
    constructor() {
        this.render();
    }

    render() {
        const formContainer = document.getElementById('formContainer');
        formContainer.innerHTML = `
            <div class="register-form">
                <h1>Inscription</h1>
                <div class="container">
                    <div class="main">
                        <div class="content">
                            <h2>Register</h2>
                            <form id="registerForm">
                                <input type="text" name="nickname" placeholder="Nickname" required>
                                <input type="number" name="age" min="0" placeholder="Age" required>
                                <select name="gender" required>
                                    <option value="" disabled selected>Gender</option>
                                    <option value="male">Male</option>
                                    <option value="female">Female</option>
                                    <option value="other">Other</option>
                                </select>
                                <input type="text" name="firstName" placeholder="First Name" required>
                                <input type="text" name="lastName" placeholder="Last Name" required>
                                <input type="email" name="email" placeholder="E-mail" required>
                                <input type="password" name="password" placeholder="Password" required>
                                <button class="btn" type="submit">Register</button>
                            </form>
                            <p class="account">Already have an account? <a href="#" id="showLogin">Log in</a></p>
                        </div>
                        <div class="form-img">
                            <img src="../styles/bg.png" alt="">
                        </div>
                    </div>
                </div>
            </div>
        `;
        document.getElementById('registerForm').addEventListener('submit', this.handleSubmit.bind(this));
    }

    async handleSubmit(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = Object.fromEntries(formData)


        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data),
            });

            const result = await response.json();
            console.log(result.message)
            if (!response.ok) {
                throw new Error(result.message || 'Registration failed');
            }



            alert('Registration successful!');
            new LoginForm();
        } catch (error) {
            console.log("Error:", error.message);
            alert('Error: ' + error.message);
        }
    }
}

class ForumPage {
    constructor() {
        this.category = [
            { name: "Tech", color: "rgb(34, 193, 195)" }, // Teal
            { name: "Finance", color: "rgb(255, 99, 71)" }, // Tomato
            { name: "Health", color: "rgb(70, 130, 180)" }, // Steel Blue
            { name: "Startup", color: "rgb(255, 223, 186)" }, // Light Goldenrod Yellow
            { name: "Innovation", color: "rgb(186, 85, 211)" }, // Medium Orchid
        ];
        this.posts = [];
        this.selectedCategories = [];
        this.maxCategories = 3; 
        this.likePost = this.likePost.bind(this);
        this.dislikePost = this.dislikePost.bind(this);
        this.render();
        this.fetchPosts();
    }

    render() {
        const forumContainer = document.getElementById('formContainer');
        forumContainer.innerHTML = `

            <div class="user">
                <h1>Real-Time-Forum</h1>
                <span id="logged-in-label">${this.getUsername()}<span>
                <button id="logoutButton">❌</button>

            </div>
            
<div class="filter-container">
    <div class="category-container">
        <label class="category-filter-label" for="categoryFilter">Filter by Category:</label>
        <select id="categoryFilter">
            <option value="">All Categories</option>
            ${this.category.map(cat => `<option value="${cat.name}">${cat.name}</option>`).join('')}
        </select>
    </div>
</div>
                <form id="postForm">
                    <input type="text" name="title" placeholder="Title:" required>
                    <textarea name="content" placeholder="What is happening?!" required></textarea>
                    <div id="categoriesContainer"></div> <!-- Category selection here -->
                    <button id="submit" type="submit">Post</button>
                </form>
            <div id="postsContainer">
            </div>
        `;
       

        this.renderCategories(); 
        document.getElementById('categoryFilter').addEventListener('change', this.filterPosts.bind(this));
        document.getElementById('postForm').addEventListener('submit', this.handlePostSubmit.bind(this));
        document.getElementById('logoutButton').addEventListener('click', this.handleLogout.bind(this));
    }

    getUsername() {
        // Retrieve the logged-in username from session or context
        return sessionStorage.getItem('username') || "Guest" // Example placeholder
    }

    async fetchPosts() {
        try {
            const response = await fetch('/api/posts', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
    
            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.message || 'Failed to fetch posts');
            }
    
            this.posts = result.map(post => ({
                username:post.username,
                title: post.title,
                content: post.content,
                category: post.category, 
                created_at: post.created_at,
                likes: post.likes || 0,
                dislikes: post.dislikes || 0,
                id: post.id,
                isLiked: false,
                isDisliked: false,
            }));

            console.log("Posts after adding new post:", this.posts); 
    
            this.displayPosts(); 
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }

    renderCategories() {
        const categoriesContainer = document.getElementById('categoriesContainer');
        categoriesContainer.innerHTML = '';

        this.category.forEach((category) => {
            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.name = "category";
            checkbox.value = category.name;
            checkbox.id = category.name;

            const label = document.createElement("label");
            label.append(document.createTextNode(category.name));
            label.style.backgroundColor = category.color;
            label.classList.add("category-label");
            
            checkbox.addEventListener("change", () => {
                this.handleCategoryChange(checkbox);
            });

            categoriesContainer.append(checkbox, label);
        });
    }

    handleCategoryChange(checkbox) {
        if (checkbox.checked) {
            if (this.selectedCategories.length < this.maxCategories) {
                this.selectedCategories.push(checkbox.value);
            } else {
                checkbox.checked = false; 
            }
        } else {
            const index = this.selectedCategories.indexOf(checkbox.value);
            if (index > -1) {
                this.selectedCategories.splice(index, 1);
            }
        }

        const checkboxes = document.querySelectorAll('#categoriesContainer input[type="checkbox"]');
        checkboxes.forEach(box => {
            if (!box.checked && this.selectedCategories.length >= this.maxCategories) {
                box.disabled = true;
            } else {
                box.disabled = false;
            }
        });
    }

    async handlePostSubmit(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = {
            username: this.getUsername(),
            title: formData.get('title'),
            content: formData.get('content'),
            category: [...this.selectedCategories], 
            created_at: new Date().toISOString(),
        };

       

        try {
            const response = await fetch('/api/post', {
                method: 'post',
                headers: {
                    'Content-Type': 'application/json',
                },
                body : JSON.stringify(data),
            });

            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.message || 'Logout failed');
            }

            if (!result.id) {
                throw new Error('Post ID is undefined');
            }

            alert("post secced");
            this.posts.unshift({
                id: result.id,
                username: result.username,
                title: result.title,
                content: result.content,
                category: result.category,
                created_at: result.created_at,
                likes: result.likes,
                dislikes: result.dislikes,
                isLiked: false,
                isDisliked: false,
            });
     // Add the new post to the beginning
            console.log("Posts after adding new post:", this.posts); 
            this.displayPosts();
            event.target.reset();
            this.selectedCategories = []; 
            this.renderCategories();
        } catch (error) {
            alert('Error: ' + error.message);
        }
       
    }

    filterPosts() {
        const categoryFilter = document.getElementById('categoryFilter').value;
        const filteredPosts = this.posts.filter(post => {
            const matchesCategory = categoryFilter ? post.category.includes(categoryFilter) : true;
            return matchesCategory;
        });
    
        this.displayPosts(filteredPosts);
    }

    displayPosts(postsToDisplay = this.posts) {
        const postsContainer = document.getElementById('postsContainer');
        postsContainer.innerHTML = '';
       
           postsToDisplay.forEach((post) => {
            const postElement = document.createElement('div');
            postElement.className = 'post';
            const formattedDate = formatDate(new Date(post.created_at));
            const TimeAgo = timeAgo(new Date(post.created_at));
            const likeCount = post.likes || 0;

            const dislikeCount = post.dislikes || 0;
    
            postElement.innerHTML = `
                <h4 style="font-weight: bold; color: #007bff;">${post.username}</h4>
                <p style="font-size: 12px; color: gray;">${TimeAgo}</p>
                <h3 style="font-size: 25px">${post.title}</h3>
                <p>${post.content}</p>
                <p>${this.getCategoryElements(post.category)}</p>
                <p style="font-size: 12px; color: gray;">${formattedDate}</p>
    
                <div class="reaction-buttons">
                    <button class="like-button" style="${post.isLiked ? 'color: green;' : ''}">👍 ${likeCount}</button>
                    <button class="dislike-button"  style="${post.isDisliked ? 'color: red;' : ''}">👎 ${dislikeCount}</button>
                    <button class="comment-button">💬 Comment</button>
                </div>
            `;

            postElement.querySelector('.like-button').addEventListener('click', () => {
                this.likePost(post.id); 
            });

            postElement.querySelector('.dislike-button').addEventListener('click', () => {
                this.dislikePost(post.id); 
            });
             
            postsContainer.appendChild(postElement);
        });
    }
    

  async likePost(postId) {
    
        try {
            const response = await fetch('/api/like', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ post_id: postId }),
            });

            if (!response.ok) {
                throw new Error('Failed to like the post');
            }

        // Récupérer le post mis à jour
        const result = await response.json();
        const post = this.posts.find(p => p.id === postId);

        if (post) {
            post.likes = result.likes; // Mettre à jour le nombre de likes
            post.dislikes = result.dislikes; // Mettre à jour le nombre de dislikes
            post.isLiked = !post.isLiked; // Inverser l'état de like

            // Si l'utilisateur a aimé, réinitialiser l'état de dislike
            if (post.isLiked) {
                post.isDisliked = false; 
            }
        }
            this.displayPosts();
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }

    async dislikePost(postId) {
        
        try {
            const response = await fetch('/api/dislike', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ post_id: postId }),
            });

            if (!response.ok) {
                throw new Error('Failed to dislike the post');
            }

        // Récupérer le post mis à jour
        const result = await response.json();
        const post = this.posts.find(p => p.id === postId);

        if (post) {
            post.dislikes = result.dislikes; // Mettre à jour le nombre de dislikes
            post.likes = result.likes; // Mettre à jour le nombre de likes
            post.isDisliked = !post.isDisliked; // Inverser l'état de dislike

            // Si l'utilisateur a disliké, réinitialiser l'état de like
            if (post.isDisliked) {
                post.isLiked = false; 
            }
        }

            this.displayPosts();
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }

    getCategoryElements(categoriesArray) {
        return categoriesArray.map(categoryName => {
            const category = this.category.find(cat => cat.name === categoryName);
            const color = category ? category.color : 'gray'; 
            return `
                <span class="category-label"   style="background-color: ${color}; font-size: 12px; ">
                    ${categoryName}
                </span>
            `;
        }).join('');
    }

    async handleLogout() {
        try {
            const response = await fetch('/api/logout', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            const result = await response.json();
            if (!response.ok) {
                throw new Error(result || 'Logout failed');
            }

            alert("You have been logged out.");
            const forumContainer = document.getElementById('formContainer');
            forumContainer.innerHTML = '';

            new ShowHomePage();
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }
}

function formatDate(date) {
    const options = {
        hour: 'numeric',
        minute: 'numeric',
        hour12: true,
        month: 'long',
        day: 'numeric',
        year: 'numeric'
    };
    return date.toLocaleString('en-US', options);
}

function timeAgo(date) {
    const now = new Date();
    const seconds = Math.floor((now - date) / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(seconds / 3600);
    const days = Math.floor(seconds / 86400);
    const months = Math.floor(seconds / 2592000); 
    const years = Math.floor(seconds / 31536000);

    if (seconds < 60) return `${seconds} second${seconds !== 1 ? 's' : ''} ago`;
    if (minutes < 60) return `${minutes} minute${minutes !== 1 ? 's' : ''} ago`;
    if (hours < 24) return `${hours} hour${hours !== 1 ? 's' : ''} ago`;
    if (days < 30) return `${days} day${days !== 1 ? 's' : ''} ago`;
    if (months < 12) return `${months} month${months !== 1 ? 's' : ''} ago`;
    return `${years} year${years !== 1 ? 's' : ''} ago`;
}

