package route

import (
	"first-project-go/api"
	"first-project-go/authorizer"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Notif")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func OpenRoute(server *api.Server) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/user/login", server.Login)
	mux.HandleFunc("/user/register", server.CreateUser)

	mux.Handle("/users", authorizer.Middleware(http.HandlerFunc(server.GetUsers)))
	mux.Handle("/user/{userId}", authorizer.Middleware(http.HandlerFunc(server.GetUserDetail)))

	corsMux := enableCORS(mux)

	err := http.ListenAndServe(":8080", corsMux)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
