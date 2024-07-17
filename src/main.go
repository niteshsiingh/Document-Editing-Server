package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
	"github.com/niteshsiingh/doc-server/src/config"
	"github.com/niteshsiingh/doc-server/src/middleware"
	ws "github.com/niteshsiingh/doc-server/src/websocket"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	db, err := config.InitializeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	_, err = config.InitEnv()
	if err != nil {
		fmt.Println("Error loading environment variables")
		log.Fatal(err)
	}

	smtp, err := config.InitSMTP()
	if err != nil {
		fmt.Println("Error loading smtp")
		log.Fatal(err)
	}

	err = initAuth()
	if err != nil {
		fmt.Println("Error initializing auth")
		log.Fatal(err)
	}

	router := createRouter(db, smtp)
	connectSocket(router, db)
	port := os.Getenv("PORT")
	_ = router.Run(":" + port)

}

func initAuth() error {
	JWT_VALIDITY, err := strconv.Atoi(os.Getenv("VALIDITY"))
	if err != nil {
		return err
	}
	JWT_REFRESH_VALIDITY, err := strconv.Atoi(os.Getenv("REFRESH_VALIDITY"))
	if err != nil {
		return err
	}
	JWT_APP_KEY := os.Getenv("JWT_APP_KEY")
	middleware.Init(JWT_APP_KEY, JWT_VALIDITY, JWT_REFRESH_VALIDITY)
	return nil
}

func connectSocket(r *gin.Engine, db *gorm.DB) {

	server := socketio.NewServer(nil)
	server.OnConnect("/", func(c socketio.Conn) error {
		u, err := url.Parse(c.URL().RawQuery)
		if err != nil {
			return err
		}
		query := u.Query()
		accessToken := query.Get("accessToken")
		documentId := query.Get("documentId")
		ws.HandleConnection(server, accessToken, documentId, db)
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Disconnected:", s.ID(), reason)
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		log.Println("Message received:", msg)
		s.Emit("reply", "Received: "+msg)
	})
	r.GET("/socket.io/*any", gin.WrapH(server), func(c *gin.Context) {

		query := c.Request.URL.Query()
		accessToken := query.Get("accessToken")
		documentId := query.Get("documentId")
		ws.HandleConnection(server, accessToken, documentId, db)
	})
}
