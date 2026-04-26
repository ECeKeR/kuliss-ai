package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type SendMessageRequest struct {
	Phone	string	`json:"phone"`
	Message	string	`json:"message"`
}

func StartServer(waClient *whatsmeow.Client) {

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		isConnected := waClient != nil && waClient.IsConnected()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":		"ok",
			"whatsapp_connected":	isConnected,
		})
	})

	http.HandleFunc("/api/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Sadece POST metodu desteklenir", http.StatusMethodNotAllowed)
			return
		}

		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
			return
		}

		if waClient == nil || !waClient.IsConnected() {
			http.Error(w, "WhatsApp şu an bağlı değil", http.StatusServiceUnavailable)
			return
		}

		jid := types.JID{
			User:	req.Phone,
			Server:	types.DefaultUserServer,
		}

		msg := &waProto.Message{
			Conversation: proto.String(req.Message),
		}

		_, err := waClient.SendMessage(context.Background(), jid, msg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Mesaj gönderilemedi: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":	true,
			"message":	"Mesaj başarıyla iletildi",
		})
	})

	log.Println("REST API Sunucusu :8080 portunda başlatıldı...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("API Sunucusu hatası: %v", err)
	}
}
