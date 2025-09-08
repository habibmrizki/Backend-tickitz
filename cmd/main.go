package main

import (
	"log"
	"os"

	"github.com/habibmrizki/back-end-tickitz/internal/configs"
	"github.com/habibmrizki/back-end-tickitz/internal/routers"
	"github.com/joho/godotenv"
	// gin-swagger middleware
)

// @title 		Back-End Tickitz
// @version		1.0
// @description	RESTful API created using gin for Back-End Tickitz
// @host		localhost:3000
// @basepath	/
func main() {
	// Manual load env
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env\nCause: ", err.Error())
		return
	}
	log.Println(os.Getenv("DBUSERS"))
	log.Println(os.Getenv("JWT_SECRET"))
	log.Println(os.Getenv("JWTISSUER"))

	// Inisialisasi DB
	db, err := configs.InitDb()
	if err != nil {
		log.Println("Failed to connect to database\nCause: ", err.Error())
		return
	}
	defer db.Close()

	if err := configs.TesbDB(db); err != nil {
		log.Println("Ping to DB failed\nCuase: ", err.Error())
		return
	}
	log.Println("Connect to db")

	router := routers.InitRouter(db)
	router.Run("localhost:3000")
}
