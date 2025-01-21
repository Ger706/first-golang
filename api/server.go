package api

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Server struct {
	listenAddr string
	DB         *gorm.DB
}

func NewServer(listenAddr string, dbUser, dbPassword, dbName string) *Server {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to MySQL database successfully.")

	return &Server{
		listenAddr: listenAddr,
		DB:         db,
	}
}
