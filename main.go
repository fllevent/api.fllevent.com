package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const versionNumber = "1.0.0.0"

var (
	help       bool
	version    bool
	portNumber int

	DB *sql.DB
)

func init() {
	flag.BoolVar(&help, "help", false, "Prints out available comands")
	flag.BoolVar(&version, "version", false, "Prints the version number")
	flag.IntVar(&portNumber, "port", 8000, "Set custome port number")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := dbconnect(dbhost(), dbname(), dbusername(), dbpassword())
	DB = db
	if err != nil {
		panic(err)
	}
}

func main() {
	if help {
		printHelp()
	}
	if version {
		printVersion()
	}

	webServerPort := fmt.Sprintf(":%d", portNumber)

	router := gin.Default()
	router.GET("/api/healthcheck", healthcheck())
	router.POST("/api/events", getallevents(DB))
	router.GET("/api/events", getallevents(DB))
	router.POST("/api/event/:name", getsingleevent(DB))
	router.GET("/api/event/:name", getsingleevent(DB))
	router.POST("/api/addevent", addevent(DB))
	router.POST("/api/removeevent", removeevent(DB))

	fmt.Println("Starting Server on port" + webServerPort)
	router.Run(webServerPort) // listen and serve on port
}
