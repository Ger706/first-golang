package route

import (
	"first-project-go/api"
	"first-project-go/authorizer"
	"first-project-go/notification"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Notif")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func OpenRoute(server *api.Server) error {
	mux := http.NewServeMux()
	notificationService, err := notification.NewNotificationService()
	if err != nil {
		log.Fatalf("Failed to initialize notification service: %v", err)
	}

	mux.HandleFunc("/send-notification", notificationService.SendNotificationHandler)

	mux.HandleFunc("/users", server.GetUsers)
	mux.HandleFunc("/user/{userId}", server.GetUserDetail)
	mux.HandleFunc("/user/login", server.Login)
	corsMux := enableCORS(mux)
	securedMux := authorizer.Middleware(corsMux)

	log.Println("Server running on port 8080...")
	return http.ListenAndServe(":8080", securedMux)
}
