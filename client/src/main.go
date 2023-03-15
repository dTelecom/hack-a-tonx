package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

func handleRequest(db *gorm.DB) {
	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/api/room/create", createRoom(db))
	e.POST("/api/room/join", joinRoom(db))
	e.POST("/api/room/callback", callbackRoom(db))
	e.POST("/api/room/info", infoRoom(db))

	e.Logger.Fatal(e.Start(":3000"))
}

func initialMigration(db *gorm.DB) {

	db.AutoMigrate(&Participant{}, &Room{}, &Call{})
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()

	db, err := gorm.Open("sqlite3", "sqlite3gorm.db")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	defer db.Close()

	initialMigration(db)
	handleRequest(db)
}
