package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/niteshsiingh/doc-server/src/config"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/middleware"
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

func connectSocket(r *gin.Engine, db *databases.Queries) {
	server := socketio.NewServer(&engineio.Options{
		PingTimeout:  180 * time.Second,
		PingInterval: 60 * time.Second,
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Connected:", s.ID())
		u, err := url.Parse(s.URL().RawQuery)
		if err != nil {
			return err
		}
		a := strings.Split(u.String(), "&")
		if len(a) > 0 {

			documentQuery := a[1]
			documentId := strings.Split(documentQuery, "=")[1]
			if documentId != "" && documentId != "undefined" {
				room := "room_" + documentId
				s.Join(room)
			} else {
				fmt.Println("No documentId provided")
			}
		}

		return nil
	})

	server.OnEvent("/", "receive-changes", func(s socketio.Conn, idArray []interface{}) {
		u, err := url.Parse(s.URL().RawQuery)
		if err != nil {
			return
		}
		a := strings.Split(u.String(), "&")
		if len(a) > 0 {

			documentQuery := a[1]
			documentId := strings.Split(documentQuery, "=")[1]
			if documentId != "" && documentId != "undefined" {
				for _, id := range idArray {
					iid := databases.IdentifierID{}
					initial := id.(map[string]interface{})["id"].(map[string]interface{})
					operation, ok := id.(map[string]interface{})["operation"]
					if !ok {
						fmt.Println("Operation not found")
						continue
					}
					input := initial
					if operation.(float64) != 1 {
						input = initial["id"].(map[string]interface{})
					}
					base, ok := input["_base"].(map[string]interface{})
					if ok {
						b, ok := base["_b"]
						if ok {
							iid.Base = []byte(fmt.Sprintf("%v", b))
						}
					}
					c, ok := input["_c"]
					if ok {
						iid.C = []byte(fmt.Sprintf("%v", c))
					}
					d, ok := input["_d"]
					if ok {
						dBytes, _ := json.Marshal(d)
						iid.D = dBytes
					}
					if s, ok := input["_s"]; ok {
						sBytes, _ := json.Marshal(s)
						iid.S = sBytes
					}
					idBytes, err := json.Marshal(iid)
					if err != nil {
						fmt.Println("Error marshalling identifier: ", err)
					}
					docIdInt, _ := strconv.Atoi(documentId)
					if operation.(float64) != 1 {
						err = db.CreateIdentifier(context.Background(), databases.CreateIdentifierParams{
							Elem:  initial["elem"].(string),
							ID:    idBytes,
							DocID: pgtype.Int4{Int32: int32(docIdInt), Valid: true},
						},
						)
						if err != nil {
							fmt.Println("Error creating identifier: ", err)
						}
					} else {
						err = db.DeleteIdentifier(context.Background(), databases.DeleteIdentifierParams{
							DocID:   pgtype.Int4{Int32: int32(docIdInt), Valid: true},
							Column2: idBytes,
						})
						if err != nil {
							fmt.Println("Error deleting identifier: ", err)
						}
					}
				}
				room := "room_" + documentId
				// server.BroadcastToRoom("/", room, "receive-changes", idArray)
				cnt := 0
				server.ForEach("/", room, func(c socketio.Conn) {
					cnt++
					if c.ID() != s.ID() {
						c.Emit("receive-changes", idArray)
					}
				})
				fmt.Println("Broadcasted to ", cnt, " clients")
			} else {
				fmt.Println("No documentId provided")
			}
		}
	})

	server.OnEvent("/", "send-changes", func(s socketio.Conn, idArray []interface{}) {
		s.Emit("receive-changes", idArray)
		u, err := url.Parse(s.URL().RawQuery)
		if err != nil {
			return
		}
		a := strings.Split(u.String(), "&")
		if len(a) > 0 && a[1][0:10] == "documentId" {
			documentQuery := a[1]
			documentId := strings.Split(documentQuery, "=")[1]
			server.BroadcastToRoom("/", "room_"+documentId, "recieve-changes", idArray)
		}
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Disconnected:", s.ID(), reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("SocketIO server error: %v", err)
		}
	}()
	// defer server.Close()

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

}
