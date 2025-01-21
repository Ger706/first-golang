package route

import (
	"first-project-go/api"
	"first-project-go/authorizer"
	"net/http"
)

func OpenRoute(server *api.Server) {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", server.GetUsers)
	mux.HandleFunc("/user/{userId}", server.GetUserDetail)
	mux.HandleFunc("/user/login", server.Login)

	err := http.ListenAndServe(":8080", authorizer.Middleware(mux))
	if err != nil {
		return
	}
}
