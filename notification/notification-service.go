package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type NotificationService struct {
	client *messaging.Client
}

func NewNotificationService() (*NotificationService, error) {
	opt := option.WithCredentialsFile("C:\\Users\\Gerry anderson\\Project\\first-golang\\golang-first-7c19e-firebase-adminsdk-fbsvc-19d5ca7c3d.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase messaging client: %v", err)
	}

	return &NotificationService{client: client}, nil
}

type NotificationRequest struct {
	Token string `json:"token"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Notification struct {
	Token string `json:"token"`
}

func (ns *NotificationService) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	message := &messaging.Message{
		Token: req.Token,
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
	}

	response, err := ns.client.Send(context.Background(), message)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending FCM message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Notification sent successfully",
		"response": response,
	})
}

func (ns *NotificationService) SendNonRequestNotification(token string, title string, body string) {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	_, err := ns.client.Send(context.Background(), message)
	if err != nil {

		return
	}

}
