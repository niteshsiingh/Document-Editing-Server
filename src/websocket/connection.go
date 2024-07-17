package websocket

import (
	"fmt"
	"strconv"

	socketio "github.com/googollee/go-socket.io"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/services"
	"gorm.io/gorm"
)

var clients = make(map[*socketio.Server]string)

func HandleConnection(conn *socketio.Server, accessToken string, documentId string, db *gorm.DB) {
	// defer conn.Close()
	if accessToken == "" || documentId == "" {
		conn.Close()
		return
	}

	token, parsedToken, err := middleware.GetAuth().ParseToken(accessToken)
	if err != nil || !token.Valid {
		fmt.Println(err)
		fmt.Println(!token.Valid)
		// conn.Close()
		return
	}
	userId := parsedToken.UserID
	docId, _ := strconv.Atoi(documentId)

	document, err := services.FindDocumentByID(uint(docId), uint(userId), db)
	if err != nil {
		conn.Close()
		return
	}
	conn.BroadcastToRoom("", "document", string([]byte("Connected to document: "+string(document.Content))))

	// conn.WriteMessage(websocket.TextMessage, []byte("Connected to document: "+string(document.Content)))

	// clients[conn] = emailId

	// for {
	// 	_, message, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Printf("error: %v", err)
	// 		delete(clients, conn)
	// 		break
	// 	}

	// 	// for client := range clients {
	// 	// 	if client != conn {
	// 	// 		client.WriteMessage(websocket.TextMessage, message)
	// 	// 	}
	// 	// }
	// }
}
