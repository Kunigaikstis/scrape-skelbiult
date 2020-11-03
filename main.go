package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gocolly/colly/v2"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := initDB()
	err := migrateDB(db)
	panicOnError(err)

	listingRepo := NewListingRepository(db)
	chatRepo := NewChatRepository(db)

	listingChan := make(chan Listing)

	searchUrl := os.Getenv("SKELBIU_LT_SEARCH_RESULTS_URL")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	panicOnError(err)

	updatesChan, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: 60,
	})
	panicOnError(err)

	startScraping(searchUrl, listingRepo, listingChan)

	for {
		select {
		case <-time.After(time.Minute * 10):
			startScraping(searchUrl, listingRepo, listingChan)
		case update := <-updatesChan:
			saveChat(update, chatRepo)
		case listing := <-listingChan:
			l, err := listingRepo.GetById(listing.Id)

			if err != nil && l.Id == 0 {
				err = listingRepo.Save(listing)

				if err != nil {
					fmt.Printf("error saving new l %v", listing)
					return
				}

				// Print link
				fmt.Printf("New link found: %s -> %s\n", listing.Id, listing.Url)

				err = sendListing(listing, chatRepo, bot)

				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Sends a listing to all chats in the db
func sendListing(listing Listing, repo *ChatRepository, bot *tgbotapi.BotAPI) error {
	chats, err := repo.GetAll()

	if err != nil {
		return err
	}

	for _, chat := range chats {
		_, err := bot.Send(tgbotapi.NewMessage(chat.Id, listing.Url))
		if err != nil {
			log.Printf("error when sending message: %v", err)
		} else {
			log.Printf("listing %v sent to chat %v.", listing.Id, chat.Id)
		}
	}

	return nil
}

// Saves the Telegram chat to the database.
func saveChat(update tgbotapi.Update, chatRepo *ChatRepository) {
	chatId := update.Message.Chat.ID

	chat, err := chatRepo.GetById(chatId)

	if err != nil && chat.Id == 0 {
		chat = Chat{
			Id:       chatId,
			Username: update.Message.From.UserName,
		}

		err = chatRepo.Save(chat)

		if err != nil {
			log.Printf("error saving chat\nerr:%v\nchat:%v", err, chat)
		} else {
			log.Printf("saved new chat (%v) to the database", chat.Id)
		}
	}
}

// Initializes the db instances
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

// Finds new listings and pushes them into the listing channel
func startScraping(searchUrl string, listingRepo *ListingRepository, listingChan chan<- Listing) {
	c := colly.NewCollector(
		colly.AllowedDomains("skelbiu.lt", "aruodas.lt", "www.skelbiu.lt", "www.aruodas.lt"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("li.simpleAds > a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		strId := e.Attr("data-item-id")
		absoluteUrl := e.Request.AbsoluteURL(link)
		id, err := strconv.ParseInt(strId, 10, 0)

		if err != nil {
			fmt.Printf("could not convert %s to int", strId)
			return
		}

		listing := Listing{
			Id:  int(id),
			Url: absoluteUrl,
		}

		listingChan <- listing
	})

	// Start scraping
	c.Visit(searchUrl)
}
