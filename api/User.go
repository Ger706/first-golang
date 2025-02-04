package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"first-project-go/authorizer"
	"first-project-go/model"
	notification2 "first-project-go/notification"
	"first-project-go/utility"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
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

	var userData User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Body decode error: %v", err)
		return
	}
	sha := sha256.New()
	sha.Write([]byte(userData.Password))
	hashedPassword := sha.Sum(nil)
	hashedPasswordHex := hex.EncodeToString(hashedPassword)
	var user = User{}
	user = User{}
	err := s.DB.Select("user_id, username").
		Where("password = ? AND username = ?", hashedPasswordHex, userData.Username).
		First(&user).Error
	token, err := authorizer.CreateToken(&jwt{Username: user.Username})
	user.Token = token
	if err == nil {
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("JSON encode error: %v", err)
			return
		}
	} else {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error Getting Account Information: %v", err)
	}
	notification, err := notification2.NewNotificationService()
	token = r.Header.Get("notif")
	if token == "" {
		http.Error(w, "Missing notification token in header", http.StatusBadRequest)
		return
	}

	notification.SendNonRequestNotification(token, "Welcome!", "Have a nice day "+user.Username)

	utility.SetJSONHeader(w)

}
