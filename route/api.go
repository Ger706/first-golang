package route

import (
	"first-project-go/api"
	"first-project-go/authorizer"
	"first-project-go/notification"
	"log"
	"net/http"
)

func OpenRoute(server *api.Server) {
	mux := http.NewServeMux()
	notificationService, err := notification.NewNotificationService()
	if err != nil {
		log.Fatalf("Failed to initialize notification service: %v", err)
	}

	mux.HandleFunc("/send-notification", notificationService.SendNotificationHandler)

	mux.HandleFunc("/users", server.GetUsers)
	mux.HandleFunc("/user/{userId}", server.GetUserDetail)
	mux.HandleFunc("/user/login", server.Login)

	err = http.ListenAndServe(":8080", authorizer.Middleware(mux))
	if err != nil {
		return
	}
}
