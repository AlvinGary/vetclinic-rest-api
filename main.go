package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"vetclinic-rest-api/database"
	"vetclinic-rest-api/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	err error
)

func main() {
	err = godotenv.Load(".env")
    if err != nil {
        panic("Error loading .env file")
    }
	dbInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
        os.Getenv("PGHOST"),
        os.Getenv("PGPORT"),
        os.Getenv("PGUSER"),
        os.Getenv("PGPASSWORD"),
        os.Getenv("PGDATABASE"),
    )
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot reach the database:", err)
	}

	database.DBMigrate(db)
	router := gin.Default()
	routers.SetupRoutes(router, db)

	port := os.Getenv("PORT")
    if port == "" {
        port = "8000" // local
    }

    log.Println("Server is running on port:", port)
    router.Run(":" + port)
}