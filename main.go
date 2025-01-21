package main

import (
	"first-project-go/api"
	"first-project-go/route"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	dbUser, _ := os.LookupEnv("DB_USER")
	listenAddr, _ := os.LookupEnv("DB_LISTEN_ADDR")
	dbPassword, _ := os.LookupEnv("DB_PASSWORD")
	dbName, _ := os.LookupEnv("DB_NAME")

	server := api.NewServer(listenAddr, dbUser, dbPassword, dbName)
	route.OpenRoute(server)

	log.Printf("Server is listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
