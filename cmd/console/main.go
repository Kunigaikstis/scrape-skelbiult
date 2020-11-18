package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Kunigaikstis/scrape-skelbiult/pkg/listing"
	"github.com/Kunigaikstis/scrape-skelbiult/pkg/message"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gocolly/colly/v2"
	_ "github.com/mattn/go-sqlite3"
)

func loadEnvFile() {
	err := godotenv.Load()
	panicOnError(err)
}

func setupDB() *sql.DB {
	db := initDB()
	err := migrateDB(db)
	panicOnError(err)

	return db
}

func main() {
	loadEnvFile()
	searchUrl := os.Getenv("SKELBIU_LT_SEARCH_RESULTS_URL")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	db := setupDB()

	c := colly.NewCollector(
		colly.AllowedDomains("skelbiu.lt", "aruodas.lt", "www.skelbiu.lt", "www.aruodas.lt"),
	)
	s := listing.NewScraper(c)
	listingRepo := listing.NewStorage(db)
	listingService := listing.NewService(listingRepo, s)
	chatRepo := message.NewStorage(db)
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	panicOnError(err)
	chatService := message.NewService(chatRepo, bot)

	updatesChan, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: 60,
	})
	panicOnError(err)

	for {
		select {
		case <-time.After(time.Minute * 10):
			log.Print("checking for new ads....")
			adsToDispatch, err := listingService.GetNewListings(searchUrl)
			panicOnError(err)
			log.Printf("found %v new ads", len(adsToDispatch))

			for _, ad := range adsToDispatch {
				msg := fmt.Sprintf("%s %s\n%s, %s\n%s", ad.Price, ad.SqFootage, ad.Street, ad.Neighbourhood, ad.Url)
				log.Printf("sending message: %s", msg)
				err := chatService.SendMessageToAllChats(msg)

				panicOnError(err)
			}
		case update := <-updatesChan:
			err := chatService.AddChat(message.Chat{
				Id:       update.Message.Chat.ID,
				Username: update.Message.Chat.UserName,
			})

			if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
				panicOnError(err)
			}
		}
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Initializes the db instance
func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./listings.db")
	panicOnError(err)

	return db
}

// Applies migrations to the db
func migrateDB(db *sql.DB) error {
	file, err := ioutil.ReadFile("./init.sql")

	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			return err
		}
	}

	return nil
}
