package main

import (
	"os"
	"testing"
)

func TestEnvLoad(t *testing.T) {
	loadEnvFile()
	searchUrl := os.Getenv("SKELBIU_LT_SEARCH_RESULTS_URL")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if searchUrl == "" {
		t.Error("searchUrl is empty")
	}

	if telegramBotToken == "" {
		t.Error("telegramBotToken is empty")
	}
}

func TestSetupDb(t *testing.T) {
	db := setupDB()
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")

	if err != nil {
		t.Error(err)
	}

	i := 0
	for rows.Next() {
		i++
	}

	if i == 0 {
		t.Error("expected at least a single table")
	}

	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestStartScraping(t *testing.T) {
	loadEnvFile()
	searchUrl := os.Getenv("SKELBIU_LT_SEARCH_RESULTS_URL")

	listingChan := make(chan Listing)

	startScraping(searchUrl, listingChan)

	listing := <-listingChan

	if listing.Id == 0 {
		t.Errorf("expected a valid listing, got %#v", listing)
	}
}
