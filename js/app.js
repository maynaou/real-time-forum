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
            console.log(result.message)
            if (!response.ok) {
                throw new Error(result.message || 'login failed');
            }


            alert('Login successful!');
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
        this.posts = [];
        this.selectedCategories = [];
        this.maxCategories = 3; 
        this.render();
    }

    render() {
        const forumContainer = document.getElementById('formContainer');
        forumContainer.innerHTML = `
            <div class="user">
                <h1>Real-Time-Forum</h1>
                <button id="logoutButton">Logout</button>
            </div>
            <div class="post-form">
                <form id="postForm">
                    <input type="text" name="title" placeholder="Title:" required>
                    <textarea name="content" placeholder="What is happening?!" required></textarea>
                    <div id="categoriesContainer"></div> <!-- Category selection here -->
                    <button id="submit" type="submit">Post</button>
                </form>
            </div>
            <div class="filter-section">
                <h2>Filter Posts</h2>
                <button id="filterByCategory">Filter by Category</button>
                <button id="filterByLikes">Filter by Liked Posts</button>
            </div>
            <div id="postsContainer">
                <h2>Posts</h2>
            </div>
        `;

        this.renderCategories(); // Call to render category checkboxes

        document.getElementById('postForm').addEventListener('submit', this.handlePostSubmit.bind(this));
        document.getElementById('filterByCategory').addEventListener('click', this.filterByCategory.bind(this));
        document.getElementById('filterByLikes').addEventListener('click', this.filterByLikes.bind(this));
        document.getElementById('logoutButton').addEventListener('click', this.handleLogout.bind(this));
    }

    renderCategories() {
        const categories = [
            { name: "Tech", color: "rgb(25, 195, 125)" },
            { name: "Finance", color: "rgb(255, 85, 85)" },
            { name: "Health", color: "rgb(30, 144, 255)" },
            { name: "Startup", color: "rgb(255, 215, 0)" },
            { name: "Innovation", color: "rgb(148,0,211)" },
        ];

        const categoriesContainer = document.getElementById('categoriesContainer');
        categoriesContainer.innerHTML = '';

        categories.forEach((category) => {
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
            title: formData.get('title'),
            content: formData.get('content'),
            categories: [...this.selectedCategories], // Use selected categories
            likes: 0,
            dislikes: 0,
        };

        // Simulate API call to save the post
        this.posts.push(data);
        this.displayPosts();
        event.target.reset();
        this.selectedCategories = []; // Reset categories after submission
        this.renderCategories(); // Re-render categories to reset checkboxes
    }

    displayPosts() {
        const postsContainer = document.getElementById('postsContainer');
        postsContainer.innerHTML = '';

        this.posts.forEach((post, index) => {
            const postElement = document.createElement('div');
            postElement.className = 'post';
            postElement.innerHTML = `
                <h3>${post.title}</h3>
                <p>${post.content}</p>
                <p>Categories: ${post.categories.join(', ')}</p>
                <p>Likes: ${post.likes} | Dislikes: ${post.dislikes}</p>
                <button onclick="forumPage.likePost(${index})">Like</button>
                <button onclick="forumPage.dislikePost(${index})">Dislike</button>
            `;
            postsContainer.appendChild(postElement);
        });
    }

    likePost(index) {
        this.posts[index].likes += 1;
        this.displayPosts();
    }

    dislikePost(index) {
        this.posts[index].dislikes += 1;
        this.displayPosts();
    }

    filterByCategory() {
        const category = prompt("Enter category to filter:");
        const filteredPosts = this.posts.filter(post => post.categories.includes(category));
        this.displayFilteredPosts(filteredPosts);
    }

    filterByLikes() {
        const likedPosts = this.posts.filter(post => post.likes > 0);
        this.displayFilteredPosts(likedPosts);
    }

    displayFilteredPosts(filteredPosts) {
        const postsContainer = document.getElementById('postsContainer');
        postsContainer.innerHTML = '';

        filteredPosts.forEach((post, index) => {
            const postElement = document.createElement('div');
            postElement.className = 'post';
            postElement.innerHTML = `
                <h3>${post.title}</h3>
                <p>${post.content}</p>
                <p>Categories: ${post.categories.join(', ')}</p>
                <p>Likes: ${post.likes} | Dislikes: ${post.dislikes}</p>
                <button onclick="forumPage.likePost(${index})">Like</button>
                <button onclick="forumPage.dislikePost(${index})">Dislike</button>
            `;
            postsContainer.appendChild(postElement);
        });
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