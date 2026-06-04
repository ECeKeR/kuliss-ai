package handler

import (
	"context"
	"log"
	"time"

	"message-go/internal/ai"
	"message-go/internal/db"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type MessageHandler struct {
	WA		*whatsmeow.Client
	AI		*ai.Client
	StartTime	time.Time
	OnNewMessage	func(phone, role, content string)
	OnProfilePic	func(phone string)
}

func (h *MessageHandler) Handle(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		h.handleMessage(v)
	}
}

func (h *MessageHandler) handleMessage(evt *events.Message) {

	if evt.Info.IsGroup {
		return
	}

	msg := evt.Message.GetConversation()
	if msg == "" {
		msg = evt.Message.GetExtendedTextMessage().GetText()
	}

	if evt.Info.IsFromMe {
		if msg != "" {
			phone := evt.Info.Chat.ToNonAD().String()
			db.SaveMessage(phone, "Me", "assistant", msg)
		}
		return
	}

	if !h.StartTime.IsZero() && evt.Info.Timestamp.Before(h.StartTime) {
		return
	}

	if time.Since(evt.Info.Timestamp) > 2*time.Minute {
		return
	}

	imgMsg := evt.Message.GetImageMessage()
	if imgMsg != nil {
		caption := imgMsg.GetCaption()
		if caption != "" {
			msg = "📷 " + caption
		} else {
			msg = "📷 [Fotoğraf]"
		}
	}

	if msg == "" {
		return
	}

	// Chat JID her zaman @s.whatsapp.net formatındadır; @lid sender JID'i ile mesaj gönderilemez.
	chatJID := evt.Info.Chat.ToNonAD()
	phone := chatJID.String()
	senderJID := evt.Info.Sender // profil fotoğrafı için orijinal sender
	log.Printf("[IN] %s: %s", phone, msg)
	log.Printf("[SYSTEM] Queued for device %s...", phone)

	if err := db.SaveMessage(phone, evt.Info.PushName, "user", msg); err != nil {
		log.Printf("mesaj kaydedilemedi: %v", err)
	}

	go func(sender types.JID, p string) {
		if h.WA != nil {
			pic, err := h.WA.GetProfilePictureInfo(context.Background(), sender, &whatsmeow.GetProfilePictureParams{})
			if err == nil && pic != nil && pic.URL != "" {
				db.SaveContactProfilePic(p, pic.URL)
				if h.OnProfilePic != nil {
					h.OnProfilePic(p)
				}
			}
		}
	}(senderJID, phone)

	if h.OnNewMessage != nil {
		h.OnNewMessage(phone, "user", msg)
	}

	history, err := db.GetHistory(phone, 20)
	if err != nil {
		log.Printf("geçmiş alınamadı: %v", err)
		history = nil
	}

	done := make(chan bool)
	go func() {
		h.WA.SendChatPresence(context.Background(), chatJID, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		ticker := time.NewTicker(8 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				h.WA.SendChatPresence(context.Background(), chatJID, types.ChatPresenceComposing, types.ChatPresenceMediaText)
			case <-done:
				return
			}
		}
	}()

	log.Printf("[SYSTEM] Waiting for response from Ollama AI module... (Model Active)")
	reply, err := h.AI.Chat(history, msg)

	done <- true
	h.WA.SendChatPresence(context.Background(), chatJID, types.ChatPresencePaused, types.ChatPresenceMediaText)

	if err != nil {
		log.Printf("ai yanıtında hata: %v", err)
		return
	}

	// AI düşünürken (yaklaşık 5-10sn) kullanıcı manuel mesaj atmış olabilir mi kontrol et.
	latestHistory, _ := db.GetHistory(phone, 1)
	if len(latestHistory) > 0 && latestHistory[0]["role"] == "assistant" {
		log.Printf("[İPTAL] AI düşünürken siz manuel cevap verdiniz, AI cevabı iptal edildi: %s", phone)
		return
	}

	if reply == "" {
		log.Printf("[UYARI] AI boş yanıt döndürdü, mesaj gönderilmiyor: %s", phone)
		return
	}

	log.Printf("[OUT] %s: %s", phone, reply)

	if err := db.SaveMessage(phone, "", "assistant", reply); err != nil {
		log.Printf("yanıt kaydedilemedi: %v", err)
	}
	if h.OnNewMessage != nil {
		h.OnNewMessage(phone, "assistant", reply)
	}

	h.sendMessage(chatJID, reply)
}

func (h *MessageHandler) sendMessage(to types.JID, text string) {
	_, err := h.WA.SendMessage(
		context.Background(),
		to,
		&waProto.Message{
			Conversation: proto.String(text),
		},
	)
	if err != nil {
		log.Printf("gönderim hatası: %v", err)
	}
}
