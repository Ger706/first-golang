package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"first-project-go/authorizer"
	"first-project-go/model"
	"first-project-go/utility"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

type User = model.User
type jwt = authorizer.JwtFormat

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var users []User
	if err := s.DB.Select("username, user_id").
		Where("deleted_at IS NOT NULL").
		Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}
	utility.SetJSONHeader(w)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("JSON encode error: %v", err)
		return
	}
}

func (s *Server) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	var userId = strings.TrimPrefix(r.URL.Path, "/user/")

	var userData User
	if err := s.DB.Where("user_id = ?", userId).First(&userData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			log.Printf("Database error: %s", err)
		}
		return
	}

	// Respond with JSON
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("JSON encode error: %v", err)
		return
	}
	utility.SetJSONHeader(w)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	var body map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utility.ErrorMessage(w, err, nil)
		return
	}

	data := body["data"].(map[string]interface{})

	password := data["password"].(string)
	username := data["username"].(string)

	sha := sha256.New()
	sha.Write([]byte(password))
	hashedPassword := sha.Sum(nil)
	hashedPasswordHex := hex.EncodeToString(hashedPassword)

	var user User
	err := s.DB.Select("user_id, username").
		Where("password = ? AND username = ?", hashedPasswordHex, username).
		First(&user).Error

	if err == nil {

		token, err := authorizer.CreateToken(&jwt{Username: user.Username, UserId: user.UserID})
		if err != nil {
			utility.ErrorMessage(w, err, nil)
			return
		}
		user.Token = token

		if err := json.NewEncoder(w).Encode(user); err != nil {
			utility.ErrorMessage(w, err, nil)
			return
		}

		//notification, err := notification2.NewNotificationService()
		//if err != nil {
		//	utility.ErrorMessage(w, err, nil)
		//} else {
		//	token := r.Header.Get("notif")
		//	if token == "" {
		//		utility.ErrorMessage(w, err, nil)
		//		//http.Error(w, "Missing notification token in header", http.StatusBadRequest)
		//		return
		//	}
		//	notification.SendNonRequestNotification(token, "Welcome!", "Have a nice day "+user.Username)
		//}

		utility.WriteJSONResponse(w, nil, user)
		return
	} else {
		msg := "Invalid Credentials"
		utility.ErrorMessage(w, err, &msg)
		return
	}
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utility.ErrorMessage(w, err, nil)
		return
	}

	data := body["data"].(map[string]interface{})

	password := data["password"].(string)
	username := data["username"].(string)
	email := data["email"].(string)

	sha := sha256.New()
	sha.Write([]byte(password))
	hashedPassword := sha.Sum(nil)
	hashedPasswordHex := hex.EncodeToString(hashedPassword)

	var checkUser User

	result := s.DB.Select("user_id, username").
		Where("username = ?", username).
		First(&checkUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("User not found")
		} else {
			fmt.Println("Database error:", result.Error)
			utility.ErrorMessage(w, nil, nil)
			return
		}
	} else {
		fmt.Println("User found:", checkUser)
		msg := "User Already Exists"
		utility.ErrorMessage(w, nil, &msg)
		return
	}

	user := User{
		Username:  username,
		Password:  hashedPasswordHex,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	err := s.DB.Create(&user).Error
	if err == nil {
		msg := "User Registered Successfully"
		utility.WriteJSONResponse(w, &msg, nil)
		return
	} else {
		utility.ErrorMessage(w, err, nil)
		return
	}
}
