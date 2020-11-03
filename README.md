Scrapes www.skelbiu.lt every 10 minutes for new listings and sends them to you via Telegram.

# Setup
1. Clone the repo: `git clone https://github.com/Kunigaikstis/scrape-skelbiult.git`.
2. Open the directory (`cd scrape-skelbiult`).
3. Run `make install`.

### Crate a .env file
You'll need a .env file at the root of the directory. Two variables are required, which have to be populated by yourself:

1. `SKELBIU_LT_SEARCH_RESULTS_URL`.
2. `TELEGRAM_BOT_TOKEN`.

An example can be found in `/examples/.env`.

### Skelbiu.lt
Refine your search on https://www.skelbiu.lt and save the resulting URL into the `SKELBIU_LT_SEARCH_RESULTS_URL` `.env` variable.

### Telegram
Create a Telegram BOT using the [official documentation](https://core.telegram.org/bots#6-botfather) and save the API token into the `TELEGRAM_BOT_TOKEN` `.env` variable.

## Important
You'll need to subscribe to your newly created bot after running the scraper. This can be done by finding the BOT through the nickname you assigned to it and sending it a message.
In case of *_unsubscribing_* (aka deleting the SQLite `.db` file), send your bot a new message once the scraper is running.

# Running the scraper
`make start`

# Authors notes
This repo was created as a fairly quick prototype without the intention of making it into a generic Telegram bot. I run it on my own machine but it could easily be configured to run on a _Cloud Function_ somewhere, given that you switch out SQLite for a cloud-hosted database (like Google Firestore).

# License
MIT License.
