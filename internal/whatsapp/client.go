package whatsapp

import (
	"context"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/mdp/qrterminal/v3"
)

func Connect(EventHandler func(interface{})) *whatsmeow.Client {
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", "file:kuliss-session.db?_foreign_keys=on", waLog.Noop)
	if err != nil {
		log.Fatalf("whatsmeow store hatası: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatalf("device store hatası: %v", err)
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	client.AddEventHandler(EventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx)
		if err = client.Connect(); err != nil {
			log.Fatalf("bağlantı hatası: %v", err)
		}

		log.Println("whatsapp ı bağlamak için QR kodu tara:")
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				log.Printf("QR durumu: %s", evt.Event)
			}
		}
	} else {
		if err = client.Connect(); err != nil {
			log.Fatalf("reconnect hatası: %v", err)
		}
		log.Println("Whatssap bağlı (mevcut sessions)")
	}

	return client
}
