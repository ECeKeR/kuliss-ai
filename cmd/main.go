package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"message-go/internal/ai"
	"message-go/internal/api"
	"message-go/internal/config"
	"message-go/internal/db"
	"message-go/internal/handler"
	"message-go/internal/whatsapp"
)

const defaultPromptTemplate = `Lütfen Kuliss AI arayüzünden prompt ayarlarınızı yapılandırın.
Kurallarınızı ve ürünlerinizi maddeler halinde yazıp 'Hızlı JSON Oluştur ve Kaydet' butonuna basabilirsiniz.
`

const defaultEnvTemplate = `# ==========================================
# ⚡ KULISS AI - AYARLAR
# ==========================================

# Veritabanı (sqlite = kurulum gerektirmez)
DB_TYPE=sqlite
DB_URL=messages.db

# Ollama AI
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=gemma4:e4b
`

func initScaffolding() {

	if _, err := os.Stat("prompt.txt"); os.IsNotExist(err) {
		log.Println("✨ İlk kurulum algılandı! 'prompt.txt' şablonu oluşturuluyor...")
		_ = os.WriteFile("prompt.txt", []byte(defaultPromptTemplate), 0644)
		log.Println("✅ prompt.txt başarıyla oluşturuldu. İşletmenize göre düzenleyebilirsiniz.")
	}

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("✨ Ayar dosyası algılanmadı! '.env' şablonu oluşturuluyor...")
		_ = os.WriteFile(".env", []byte(defaultEnvTemplate), 0644)
		log.Println("✅ .env dosyası oluşturuldu. Modeli değiştirmek isterseniz bu dosyayı kullanabilirsiniz.")
	}
}

func main() {

	initScaffolding()

	cfg := config.Load()

	db.Connect(cfg.DBType, cfg.DBURL)

	aiClient := ai.NewClient(cfg.OllamaURL, cfg.OllamaModel)

	msgHandler := &handler.MessageHandler{
		AI: aiClient,
	}

	waClient := whatsapp.Connect(msgHandler.Handle)
	msgHandler.WA = waClient

	log.Println("Kuliss AI çalışıyor mesaj bekleniyor...")

	go api.StartServer(waClient)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("kapatılıyor..")
	waClient.Disconnect()
}
