package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.New("").ParseGlob("templates/*.html"))
)

var upgrader = websocket.Upgrader{}

func (fe *frontendServer) homeHandler(w http.ResponseWriter, r *http.Request) {

	if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{}); err != nil {
		logrus.Error(err)
	}
}

// WebSocket connection handler
func (fe *frontendServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	fe.startSQSReceiver(conn)

}

func (fe *frontendServer) startSQSReceiver(conn *websocket.Conn) {

	msgChan := make(chan string)
	errorChan := make(chan error)
	if fe.eventsManager == nil {
		msg := "No events will be received because the events manager component has not been initialized."
		if err := conn.WriteJSON(map[string]string{"Message": msg}); err != nil {
			fmt.Println("Error broadcasting message:", err)
		}
		conn.Close()
		return
	}
	go fe.eventsManager.ReceiveMessages(msgChan, errorChan)
	//receive msgs and errors
	go func() {
		defer conn.Close()
		for {
			select {
			case err := <-errorChan:
				logrus.Error(err)
			case msg := <-msgChan:
				if err := conn.WriteJSON(map[string]string{"Message": msg}); err != nil {
					fmt.Println("Error broadcasting message:", err)
					conn.Close()
				}
			}
		}
	}()
}
