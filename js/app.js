document.addEventListener('DOMContentLoaded', function () {
    const currentPage = sessionStorage.getItem('currentPage')

    navigateToPage(currentPage)

    document.body.addEventListener('click', function (event) {
        if (event.target.matches('#showRegister')) {
            new RegisterForm();
        } else if (event.target.matches('#showLogin')) {
            new LoginForm();
        } 
    });
});

function navigateToPage(page) {
    const postId = sessionStorage.getItem('post_id')
    const posts = JSON.parse(sessionStorage.getItem('posts'));
    const user = sessionStorage.getItem('user')
    switch (page) {
        case 'login':
            new LoginForm();
            break;
        case 'register':
            new RegisterForm();
            break;
        case 'forum':
            new ForumPage();
            break;
        case 'commentPage':
            new CommentPage(postId, posts);
            break
        case 'messagePage':
            new Message(user);
            break
        default:
            new RegisterForm();
            break;
    }
}


class LoginForm {
    constructor() {
        this.render();

        sessionStorage.setItem('currentPage', 'login');
        this.checkAuth();
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
                            <div id="loginMessage"></div>
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
    async checkAuth() {
        try {
            const response = await fetch('/api/login', { method: 'GET' });
            if (response.ok) {
                const result = await response.json();
                console.log(result);
                
                if (result.authenticated === "true") {
                    new ForumPage(); 
                    return;
                }
            }
        } catch (error) {
            console.error('Error checking authentication:', error);
        }
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
                password: data.password,
            }
        } else {
            user = {
                nickname: data.NicknameOREmail,
                password: data.password,
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

            const messageElement = document.getElementById('loginMessage'); 

            if (response.status === 401) {
                messageElement.textContent = 'Unauthorized: Please check your credentials.'; // Display message
                messageElement.style.color = 'red';
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm()
                return; 
            }

            const result = await response.json()
            if (!response.ok) {
                throw new Error(result.message || 'login failed');
            }
            new ForumPage();
        } catch (error) {
            const messageElement = document.getElementById('loginMessage'); 
            messageElement.textContent = 'Error: ' + "password or nickname invalid"; 
            messageElement.style.color = 'red'; 
        }
    }
}

class RegisterForm {
    constructor() {
        this.render();
        sessionStorage.setItem('currentPage', 'register');
        this.checkAuth();
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
    
    async checkAuth() {
        try {
            const response = await fetch('/api/register', { method: 'GET' });
            if (response.ok) {
                const result = await response.json();
                console.log(result);
                
                if (result.authenticated === "true") {
                    new ForumPage(); 
                    return;
                }
            }
        } catch (error) {
            console.error('Error checking authentication:', error);
        }
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

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm()
                return; 
            }

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
        this.comments = []
        this.selectedCategories = [];
        this.maxCategories = 3;
        this.likePost = this.likePost.bind(this);
        this.dislikePost = this.dislikePost.bind(this);
        this.init()
        
    }

    async init() {
        await this.fetchUser();
        this.render();
        this.selectedFilter = sessionStorage.getItem('categoryFilter') || '';
        sessionStorage.setItem('currentPage', 'forum');
        await this.fetchPosts(); 
        this.hasHighlightedUser = false; 
        this.connectWebSocket();
    }

    connectWebSocket() {
        this.ws = new WebSocket("ws://localhost:8084/ws"); 
        this.ws.onopen = () => {
            console.log('WebSocket connection established');
        };

        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);  
            if (data.users) {
                this.displayUsers(data.users,data.sender,data.receiver);   
            } else {
                this.hasHighlightedUser = true;  
            }      
         
        };

        this.ws.onerror = (error) => {
            console.error(`WebSocket error: ${error}`);
        };

        this.ws.onclose = () => {
            console.log("WebSocket connection closed.");
        };
        
    }


    displayUsers(users,sender,receiver) {
        const userList = document.getElementById('userList');
        userList.innerHTML = '';

        console.log(sessionStorage.getItem('username'));
       
        console.log(users);

        if (users === undefined) {
            return
        }
        
        users.forEach(user => {
            if (user.nickname !== sessionStorage.getItem('username')) {
                const userItem = document.createElement('div');
                userItem.className = 'user-item';
                userItem.classList.add(user.online ? 'online' : 'offline');
                userItem.innerText = user.nickname;
           
                if ((user.nickname === sender) &&  this.hasHighlightedUser && receiver != "") {
                    userItem.style.backgroundColor = '#4e34b6'
                }

                if (user.online) {
                    sessionStorage.setItem('user', user.nickname);
                    userItem.addEventListener('click', function () {
                        userItem.style.backgroundColor = ''
                        new Message(user.nickname,users);

                    });
                }

                userList.appendChild(userItem);
            }
        });
    }
    
    

    render() {
        const forumContainer = document.getElementById('formContainer');
        forumContainer.innerHTML = `

            <div class="user">
                <h1>Real-Time-Forum</h1>
                <span id="logged-in-label">${this.getUsername()}<span>
                <button id="logoutButton">❌</button>
            </div>

            <div id="userListContainer">
            <h2>Users</h2>
            <div class="user-list" id="userList">
            </div>
            </div>
            
            <div id="postsContainer">
            </div>
        `;


        this.resetView()

        document.getElementById('logoutButton').addEventListener('click', this.handleLogout.bind(this));
    }

    getUsername() {
        // Retrieve the logged-in username from session or context
        return sessionStorage.getItem('username')  // Example placeholder
    }

    async fetchUser() {
        try {
            const response = await fetch('/api/user', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',

                },
            }); // Assurez-vous que cet endpoint existe
            if (!response.ok) throw new Error('Failed to fetch users');

            const user = await response.json();
            console.log("user",user);
            
            sessionStorage.setItem('username', user.nickname);
        } catch (error) {
            console.error('Error fetching users:', error);
        }
    }







    resetView() {
        const forumContainer = document.getElementById('formContainer');

        if (!document.getElementById('postForm')) {
            forumContainer.insertAdjacentHTML('afterbegin', `
                <form id="postForm">
                    <input type="text" name="title" placeholder="Title:" required>
                    <textarea name="content" placeholder="What is happening?!" required></textarea>
                    <div id="categoriesContainer"></div> <!-- Category selection here -->
                    <button id="submit" type="submit">Post</button>
                </form>
            `);
        }

        if (!document.getElementById('filter-container')) {
            forumContainer.insertAdjacentHTML('afterbegin', `
                <div id="filter-container" class="filter-container">
                    <div class="category-container">
                        <label class="category-filter-label" for="categoryFilter">Filter by Category:</label>
                        <select id="categoryFilter">
                            <option value="">All Categories</option>
                            ${this.category.map(cat => `<option value="${cat.name}">${cat.name}</option>`).join('')}
                            <option value="myposts">My Own Posts</option>
                            <option value="likedposts">The Posts That I Liked</option>
                        </select>
                    </div>
                </div>
            `);
        }
        this.renderCategories()


        document.getElementById('categoryFilter').addEventListener('change', this.filterPosts.bind(this));
        document.getElementById('postForm').addEventListener('submit', this.handlePostSubmit.bind(this));
    }

    async fetchPosts() {
        try {
            const response = await fetch('/api/post', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',

                },
            });

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm()
                return; // Sortir de la fonction
            }

            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.message || 'Failed to fetch posts');
            }
            this.posts = result.map(post => ({
                username: post.username,
                title: post.title,
                content: post.content,
                category: post.category,
                created_at: post.created_at,
                likes: post.likes || 0,
                dislikes: post.dislikes || 0,
                comments: post.comments || 0,
                id: post.id,
                isDisliked: false,
                isLiked: post.isLiked,
            }));


            if (this.selectedFilter) {
                document.getElementById('categoryFilter').value = this.selectedFilter;
                this.filterPosts();
            } else {
                this.displayPosts();
            }

        } catch (error) {
            alert('Error: ' + error.message);
        }
    }

    renderCategories() {
        const categoriesContainer = document.getElementById('categoriesContainer');
        if (!categoriesContainer) {
            console.error('categoriesContainer is not found in the DOM.');
            return;
        }
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
            const response = await fetch('/api/post/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm()
                return;
            }

            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.message || 'Logout failed');
            }

            if (!result.id) {
                throw new Error('Post ID is undefined');
            }

            alert("post secced");
            this.posts.unshift({ ...result, isLiked: false, isDisliked: false });
            document.getElementById('categoryFilter').value = '';
            console.log("Posts after adding new post:", this.posts);
            this.fetchPosts();
            event.target.reset();
            this.selectedCategories = [];
            this.renderCategories();
        } catch (error) {
            alert('Error: ' + error.message);
        }

    }

    filterPosts() {
        const categoryFilter = document.getElementById('categoryFilter').value;
        sessionStorage.setItem('categoryFilter', categoryFilter);
        let filteredPosts;

        if (categoryFilter === "myposts") {
            const username = this.getUsername();
            filteredPosts = this.posts.filter(post => post.username === username);
        } else if (categoryFilter === "likedposts") {
            filteredPosts = this.posts.filter(post => post.isLiked);
        } else if (categoryFilter) {
            filteredPosts = this.posts.filter(post => post.category && post.category.includes(categoryFilter));
        } else {
            filteredPosts = this.posts;
        }

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
            const commentCount = post.comments || 0;
            const dislikeCount = post.dislikes || 0;

            postElement.innerHTML = `
               <div class="user-info">
    <div class="avatar">
        <img src="../styles/user.png" alt="User Avatar">
    </div>
    <h4>${post.username}</h4>
    <p>${TimeAgo}</p>
</div>
  
                <h3>${post.title}</h3>
                <p>${post.content}</p>
                <p>${this.getCategoryElements(post.category)}</p>
                <p style="font-size: 12px; color: gray;">${formattedDate}</p>
    
                <div class="reaction-buttons">
                    <button class="like-button" style="${post.isLiked ? 'color: blue;' : ''}">👍 ${likeCount}</button>
                    <button class="dislike-button"  style="${post.isDisliked ? 'color: red;' : ''}">👎 ${dislikeCount}</button>
                    <button class="comment-button">💬 ${commentCount}</button>
                </div>
            `;
            postElement.querySelector('.comment-button').addEventListener('click', () => {
                new CommentPage(post.id, this.posts);
            });

            postElement.querySelector('.like-button').addEventListener('click', () => {
                this.likePost(post.id, this.comments);
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

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm()
                return
            }

            const result = await response.json();

            if (!response.ok) {
                throw new Error('Failed to like the post');
            }
            const post = this.posts.find(p => p.id === postId);
            if (post) {
                post.likes = result.likes;
                post.dislikes = result.dislikes;
                post.isLiked = !post.isLiked;

                if (post.isLiked) {
                    post.isDisliked = false;
                }
                this.displayPosts();
            }

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

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm()
                return;
            }

            // Récupérer le post mis à jour
            const result = await response.json();

            if (!response.ok) {
                throw new Error('Failed to dislike the post');
            }
            const post = this.posts.find(p => p.id === postId);
            if (post) {
                post.dislikes = result.dislikes;
                post.likes = result.likes;
                post.isDisliked = !post.isDisliked;
                if (post.isDisliked) {
                    post.isLiked = false;
                }
                this.displayPosts();
            }
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
        this.ws.send(JSON.stringify({ content: "logout" }));
        await new Promise(resolve => setTimeout(resolve, 100));
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
            sessionStorage.clear();
            const forumContainer = document.getElementById('formContainer');
            forumContainer.innerHTML = '';
            this.ws.close();
            new LoginForm();
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }
}

class CommentPage {
    constructor(postId, posts) {
        this.category = [
            { name: "Tech", color: "rgb(34, 193, 195)" }, // Teal
            { name: "Finance", color: "rgb(255, 99, 71)" }, // Tomato
            { name: "Health", color: "rgb(70, 130, 180)" }, // Steel Blue
            { name: "Startup", color: "rgb(255, 223, 186)" }, // Light Goldenrod Yellow
            { name: "Innovation", color: "rgb(186, 85, 211)" }, // Medium Orchid
        ];
        this.postId = postId;
        this.posts = posts;
        this.comments = [];
        this.render();
        sessionStorage.setItem('currentPage', 'commentPage');
        sessionStorage.setItem('post_id', postId)
        sessionStorage.setItem('posts', JSON.stringify(posts));
        this.fetchComments()
    }

    render() {
        const categoryContainer = document.querySelector('.filter-container');
        const postForm = document.getElementById('postForm');

        if (categoryContainer) {
            categoryContainer.remove();
        }

        if (postForm) {
            postForm.remove();
        }

        const post = this.posts.find(p => p.id === this.postId);
        if (!post) return;

        let postsContainer = document.getElementById('postsContainer')
        if (!postsContainer) {
            const forumContainer = document.getElementById('formContainer');
            forumContainer.innerHTML = `
           <div class="user">
               <h1>Real-Time-Forum</h1>
               <span id="logged-in-label">${sessionStorage.getItem('username')}<span>
               <button id="logoutButton">❌</button>
           </div>
       `;
            postsContainer = document.createElement('div')
            postsContainer.id = 'postContainer'
            forumContainer.appendChild(postsContainer)
        }
        postsContainer.innerHTML = `
            <div class="post">
                <div class="user-info">
                    <div class="avatar">
                        <img src="../styles/user.png" alt="User Avatar">
                    </div>
                    <h4>${post.username}</h4>
                    <p>${timeAgo(new Date(post.created_at))}</p>
                </div>
                <h3>${post.title}</h3>
                <p>${post.content}</p>
                <p>${this.getCategoryElements(post.category)}</p>
                <p style="font-size: 12px; color: gray;">${formatDate(new Date(post.created_at))}</p>
            </div>
                 <div  id="postForm">
                <input type="text" id="comment-input-${this.postId}" placeholder="Votre commentaire" />
                    <div class="reaction-buttons">
                    <button class="like-button" id="add-comment-button-${this.postId}">Ajouter un commentaire</button>
                    <button class="dislike-button" id="back-button">Retour</button>
                </div>
                </div>
            <div id="comments-list">
            </div>`;

        document.getElementById(`add-comment-button-${this.postId}`).addEventListener('click', () => {
            const commentInput = document.getElementById(`comment-input-${this.postId}`);
            this.addComment(commentInput.value);
            commentInput.value = '';
        });

        document.getElementById('back-button').addEventListener('click', () => {
            new ForumPage()
        });
    }

    renderComments() {
        const commentsList = document.getElementById("comments-list");
        commentsList.innerHTML = ''

        this.comments.forEach((comment) => {
            const commentElement = document.createElement('div');
            commentElement.className = 'post';
            const formattedDate = formatDate(new Date(comment.created_at));
            const TimeAgo = timeAgo(new Date(comment.created_at));
            const likeCount = comment.likes || 0;
            const dislikeCount = comment.dislikes || 0;
            commentElement.innerHTML = `
                <div class="user-info">
                <div class="avatar">
                    <img src="../styles/user.png" alt="User Avatar" style="width: 30px; height: 30px;">
                </div>
                    <h4>${comment.username}</h4> <!-- Utiliser le nom d'utilisateur du post si non spécifié -->
                    <p>${TimeAgo}</p>
                </div>
                <p>${comment.content}</p>
                <p style="font-size: 12px; color: gray;">${formattedDate}</p>
                 <div class="reaction-buttons">
                    <button class="like-button" >👍 ${likeCount}</button>
                    <button class="dislike-button">👎 ${dislikeCount}</button>
                </div>
            `
            commentElement.querySelector('.like-button').addEventListener('click', () => {
                this.likeComment(comment.id);
            });

            commentElement.querySelector('.dislike-button').addEventListener('click', () => {
                this.dislikeComment(comment.id);
            });

            commentsList.appendChild(commentElement)
        })
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

    async fetchComments() {
        try {
            const response = await fetch('/api/comment', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    "X-Requested-With": this.postId,
                },
            });

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm();
                return;
            }
            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.message || 'Failed to fetch comments');
            }

            this.comments = result.map(comment => ({
                username: comment.username,
                content: comment.content,
                created_at: comment.created_at,
                id: comment.id,
                likes: comment.likes || 0,
                dislikes: comment.dislikes || 0,
            }));


            this.renderComments();
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }


    async addComment(comment) {
        if (!comment) return;
        try {
            const response = await fetch('/api/comment/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },

                body: JSON.stringify({ post_id: this.postId, content: comment, created_at: new Date().toISOString(), }),
            });

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm();
                return;
            }
            const newComment = await response.json();
            if (!response.ok) {
                throw new Error('Erreur lors de l\'ajout du commentaire');
            }

            this.comments.unshift({ ...newComment, isLiked: false, isDisliked: false });


            this.fetchComments();
        } catch (error) {

            alert('Erreur: ' + error.message);
        }
    }

    async likeComment(commentId) {
        try {
            const response = await fetch('/api/like', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ post_id: commentId }),
            });

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm();
                return;
            }



            // Récupérer le post mis à jour
            const result = await response.json();
            if (!response.ok) throw new Error('Failed to like the comment');

            const comment = this.comments.find(c => c.id === commentId);
            if (comment) {
                comment.likes = result.likes;
                comment.dislikes = result.dislikes;
                this.renderComments();
            }

        } catch (error) {
            alert('Error: ' + error.message);
        }
    }

    async dislikeComment(commentId) {
        try {
            const response = await fetch('/api/dislike', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ post_id: commentId }),
            });
            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm();
                return;
            }
            const result = await response.json();
            if (!response.ok) {
                throw new Error(result || 'Logout failed');
            }

            const comment = this.comments.find(c => c.id === commentId);
            if (comment) {
                comment.dislikes = result.dislikes;
                comment.likes = result.likes;
                this.renderComments();
            }
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }
}

let lastLoadedTimestamp = null;

let isFetching = false

class Message {
    constructor(username) {
        this.username = username;
        this.older = false;
        this.b = false;
        this.hasHighlightedUser = false;
        sessionStorage.setItem('currentPage', 'messagePage');
        this.init()
    }

    async init() {
        this.render();
        this.getMessage()
        this.connectWebSocket();
        this.cookie = this.getCookie("session_id")
        console.log(this.cookie);
    }
    
    getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        
        if (parts.length === 2) {
            return parts.pop().split(';').shift();
        }
        
        return null; 
    }

    render() {
        const forumContainer = document.getElementById('formContainer');
       forumContainer.innerHTML = ''
        forumContainer.innerHTML = `
            <div class="user">
               <h1>Real-Time-Forum</h1>
               <span id="logged-in-label">${sessionStorage.getItem('username')}<span>
               <button id="logoutButton">❌</button>
           </div>
            <div id="chatContainer">
                <div id="messageFormContainer">
                    <div class="user-info">
                        <div class="avatar">
                            <img src="../styles/user.png" alt="User Avatar">
                        </div>
                        <h4>${this.username}</h4>
                    </div>
                    <div id="messagesContainer"></div>
                    <form id="messageForm">
                        <textarea id="messageContent" placeholder="Type your message here..." required></textarea>
                        <div class="reaction-buttons">
                            <button type="button" class="dislike-button" id="back-button">Back</button>
                            <button type="submit" class="like-button" id="send">Send</button>
                        </div>
                    </form>
                </div>
            </div>

             <div id="userListContainer">
            <h2>Users</h2>
            <div class="user-list" id="userList">
            </div>
            </div>
        `;

  

        document.getElementById('messageForm').addEventListener('submit', (event) => {
            event.preventDefault(); // Empêche le rechargement de la page
            const commentInput = document.getElementById('messageContent');
            const messageText = commentInput.value.trim(); // Supprime les espaces
            
            if (messageText) {
                this.addComment(messageText);
                commentInput.value = '';
            } else {
                alert("Message cannot be empty."); // Alerte si le message est vide
            }
        });

        document.getElementById('back-button').addEventListener('click', () => {
            new ForumPage();
        });

        document.getElementById('logoutButton').addEventListener('click',this.handleLogout.bind(this));

        const messageBoxContent = document.getElementById('messagesContainer');
        messageBoxContent.addEventListener('scroll', this.debounce(() => {
            if (isFetching) return;
            if (messageBoxContent.scrollTop === 0) {
                console.log("Loading more messages...");
                this.older = true
                this.getMessage(); // Demande des messages plus anciens
            }
        }, 200));
    }





    async getMessage() {
        if (isFetching) return;
        isFetching = true;
        let url = `/api/message?sender=${sessionStorage.getItem('username')}&receiver=${this.username}`;

        if (this.older && lastLoadedTimestamp) {
            url += `&before=${lastLoadedTimestamp}`;
        }

        try {
            const response = await fetch(url, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
            });

            if (response.status === 401) {
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm();
                return;
            }

            if (!response.ok) {
                const errorMessage = await response.json();
                throw new Error(errorMessage || 'Message failed to fetch');
            }

    

            const messages = await response.json();
            console.log(messages);

            const messagesContainer = document.getElementById('messagesContainer')

            if (!this.older) {
                messagesContainer.innerHTML = ''
            }

            if (!Array.isArray(messages)) {
                isFetching = false;
                console.error("Expected an array of messages, but got:", messages);
                return; // Sortir de la fonction si ce n'est pas un tableau
            }

            messages.forEach((messageData, index) => {
                if (messageData.sender === sessionStorage.getItem('username')) {
                    this.displaySentMessage(messageData);
                } else {
                    this.displayReceivedMessage(messageData.content, messageData.created_at);
                    
                }
                if (index === messages.length - 1) {
                    lastLoadedTimestamp = messageData.created_at;
                    console.log(lastLoadedTimestamp);
                }
                messagesContainer.scrollTop = messagesContainer.scrollHeight;
            });

            isFetching = false;
        } catch (error) {
            alert('Error: ' + error.message);
            isFetching = false;
        }
    }

    debounce(func, delay) {
        let timer;
        return function () {
            clearTimeout(timer);
            timer = setTimeout(func, delay);
        };
    }

    connectWebSocket() {
        this.ws = new WebSocket("ws://localhost:8084/ws");
        this.ws.onopen = () => {
            console.log('WebSocket connection established');
        };

        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log("receiver", data)
            if (!data.created_at) {
                this.displayUsers(data.users,data.sender,data.receiver)
            } else {
                this.b = true
                this.displayReceivedMessage(data.content, data.created_at);
                console.log(`Received message: ${data.sender}: ${data.content}`);
                this.hasHighlightedUser = true; 
            }
        };

        this.ws.onerror = (error) => {
            console.error(`WebSocket error: ${error}`);
        };

        this.ws.onclose = () => {
            console.log("WebSocket connection closed.");
        };
    }


    displayUsers(users,sender,receiver) {
        const userList = document.getElementById('userList');
        userList.innerHTML = '';
        users.forEach(user => {
            if (user.nickname !== sessionStorage.getItem('username')) {
                const userItem = document.createElement('div');
                userItem.className = 'user-item';
                userItem.classList.add(user.online ? 'online' : 'offline');
                userItem.innerText = user.nickname;
                console.log(this.hasHighlightedUser,"khfjksdhfhjsk",sender);
                if ((user.nickname  === sender) && this.hasHighlightedUser && receiver != '') {
                         userItem.style.backgroundColor = '#4e34b6'
                         this.hasHighlightedUser = false; 
                }

                if (user.online) {
                    sessionStorage.setItem('user', user.nickname);
                    userItem.addEventListener('click', function () {
                        new Message(user.nickname);

                    });
                }
                userList.appendChild(userItem);
            }
        });
    }



   async addComment(message) {
        if (this.ws.readyState === WebSocket.OPEN) {
            let y = this.getCookie("session_id")
            
            if (this.cookie !== y) {                  
                this.ws.send(JSON.stringify({ cookie: this.cookie })); 
                const forumContainer = document.getElementById('formContainer');
                forumContainer.innerHTML = '';
                new LoginForm();
                return;
            }

            const messageData = {
                receiver: this.username,
                content: message,
                created_at: new Date().toISOString(),
            };

            this.ws.send(JSON.stringify(messageData));
            this.b = true;
            this.displaySentMessage(messageData);
        }
    }

    displaySentMessage(messageData) {
        const messagesContainer = document.getElementById('messagesContainer');
        const messageItem = document.createElement('div');
        messageItem.classList.add('message', 'sender');
        const createdAt = new Date(messageData.created_at);
        const options = {
            hour: '2-digit',
            minute: '2-digit',
            hour12: false 
        };

        const formattedTime = createdAt.toLocaleTimeString([], options);
        messageItem.textContent = `${messageData.content} ${formattedTime}`;
        if (this.b) {
            messagesContainer.appendChild(messageItem);
            this.b = false;
        } else {
            messagesContainer.prepend(messageItem);
        }

    }

    displayReceivedMessage(message, time) {
        const messagesContainer = document.getElementById('messagesContainer');
        const messageItem = document.createElement('div');
        messageItem.classList.add('message', 'receiver');
        const createdAt = new Date(time);
        const options = {
            hour: '2-digit',
            minute: '2-digit',
            hour12: false // Utiliser le format 24 heures
        };

        // Formater le temps
        const formattedTime = createdAt.toLocaleTimeString([], options);
        messageItem.textContent = `${message} ${formattedTime}`;
        if (this.b) {
            messagesContainer.appendChild(messageItem);
            this.b = false;
        } else {
            messagesContainer.prepend(messageItem);
        }
    }

    async handleLogout() {
        try {
        this.ws.send(JSON.stringify({ content: "logout" }));
        await new Promise(resolve => setTimeout(resolve, 100));
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
            sessionStorage.clear();
            const forumContainer = document.getElementById('formContainer');
            forumContainer.innerHTML = '';
            this.ws.close();
            new LoginForm();
        } catch (error) {
            alert('Error: ' + error.message);
        }
    }
}


function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    
    if (parts.length === 2) {
        return parts.pop().split(';').shift();
    }
    
    return null; // Cookie not found
}

// Example usage
const myCookie = getCookie('session_id');




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

