package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	models "handler/DataBase/Models"
	database "handler/DataBase/Sqlite"
)

func ValidateRegisterFornData(user models.RegisterRequest) map[string]string {
	errors := make(map[string]string)

	if user.Nickname == "" {
		errors["nickname"] = "nickname cannot be empty"
	} else {
		if len(user.Nickname) < 2 || len(user.Nickname) > 15 {
			errors["nickname"] = "nickname must be between 2 and 15 characters in length"
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(user.Nickname) {
			errors["nickname"] = "nickname can only contain letters, numbers, and underscores"
		}
	}

	if user.FirstName == "" {
		errors["first_name"] = "first name cannot be empty"
	}

	if user.LastName == "" {
		errors["last_name"] = "last name cannot be empty"
	}

	if user.Email == "" {
		errors["email"] = "email cannot be empty"
	} else if !regexp.MustCompile(`\S+@\S+\.\S+`).MatchString(user.Email) {
		errors["email"] = "invalid email format"
	}

	age, err := strconv.Atoi(user.Age)
	if err != nil || age < 0 {
		errors["age"] = "invalid age"
	}

	if user.Gender != "male" && user.Gender != "female" {
		errors["gender"] = "please select your gender"
	}

	if user.Password == "" {
		errors["password"] = "password cannot be empty"
	} else if len(user.Password) < 6 || len(user.Password) > 30 {
		errors["password"] = "password must be between 6 and 30 characters in length"
	}

	return errors
}

func ValidateLoginFormInput(NicknameOREmail, password string) map[string]string {
	errors := make(map[string]string)
	var req models.LoginRequest
	if NicknameOREmail == "" {
		errors["usernameOrEmail"] = "username or email is required"
	} else {
		emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
		isEmail := emailRegex.MatchString(NicknameOREmail)
		if isEmail {
			req.Email = NicknameOREmail
		} else {
			req.Nickname = NicknameOREmail
		}
	}

	if password == "" {
		errors["password"] = "password is required"
	} else if len(password) < 6 {
		errors["password"] = "password should be at least 6 characters long"
	}

	req.Password = password

	return errors
}

func ValidationPostInput(post models.Post) map[string]string {
	errors := make(map[string]string)

	if post.Title == "" {
		errors["title"] = "title cannot be empty"
	} else if len(post.Title) > 50 {
		errors["title"] = "title cannot exceed 50 characters"
	}

	if post.Content == "" {
		errors["content"] = "content cannot be empty"
	} else if len(post.Content) > 1000 {
		errors["content"] = "content cannot exceed 1000 characters"
	}

	if len(post.Categories) == 0 {
		errors["categories"] = "at least one category has to be selected"
	}

	return errors
}

func ValidateCommentInput(comment models.Comment) map[string]string {
	errors := make(map[string]string)

	if comment.Content == "" {
		errors["content"] = "comment cannot be empty"
	} else if len(comment.Content) > 250 {
		errors["content"] = "comment cannot exceed 250 characters"
	}

	return errors
}

func ValidateSession(sessionID string) error {
	var userID string
	query := `SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?`
	err := database.DB.QueryRow(query, sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return fmt.Errorf("unauthorized: session not found")
	}

	return nil
}
