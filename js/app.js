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
        this.element.className = 'content';
        document.getElementById('formContainer').appendChild(this.element);
        this.showHomePage();
    }

    showHomePage() {
        const content = `
            <h1>WELCOME TO REAL-TIME-FORUM</h1>
            <p class="subtext">Join our community and explore new ideas</p>
            <p class="subtext">Connect, share, and learn with fellow enthusiasts</p>
            <div class="buttons">
                <button class="btn" id="showLogin">Sign In</button>
                <button class="btn" id="showRegister">Sign Up</button>
                <button class="btn" onclick="window.location.href = '/guest'">Guest</button>
            </div>
        `;
        this.element.innerHTML = content; 
    }
}

class LoginForm {
    constructor() {
        this.render();
    }

    render() {
        const formContainer = document.getElementById('formContainer');
        formContainer.innerHTML = `
            <h1>Connexion</h1>
            <form id="loginForm">
                <input type="text" name="identifier" placeholder="Nickname or E-mail" required>
                <input type="password" name="password" placeholder="Password" required>
                <button type="submit">Log In</button>
            </form>
            <p>Don't have an account? <a href="#" id="showRegister">Register</a></p>
        `;

        document.getElementById('loginForm').addEventListener('submit', this.handleSubmit.bind(this));
    }

    async handleSubmit(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = Object.fromEntries(formData);

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                body: JSON.stringify(data),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.error || 'Login failed');
            }

            alert('Login successful!');
            new ShowHomePage(); 
        } catch (error) {
            alert(error.message);
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
            <h1>Inscription</h1>
            <form id="registerForm">
                <input type="text" name="nickname" placeholder="Nickname" required>
                <input type="number" name="age" placeholder="Age" required>
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
                <button type="submit">Register</button>
            </form>
            <p>Already have an account? <a href="#" id="showLogin">Log in</a></p>
        `;

        document.getElementById('registerForm').addEventListener('submit', this.handleSubmit.bind(this));
    }

    async handleSubmit(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = Object.fromEntries(formData);

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                body: JSON.stringify(data),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.error || 'Registration failed');
            }

            alert('Registration successful!');
            new ShowHomePage();
        } catch (error) {
            alert(error.message);
        }
    }
}