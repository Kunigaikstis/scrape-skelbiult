package listing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	collector *colly.Collector
}

func NewScraper(allowedDomains ...string) *Scraper {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
	)
	return &Scraper{collector: c}
}

func (s *Scraper) GetListings(url string) []Ad {
	detailCollector := s.collector.Clone()
	// On every a element which has href attribute call callback
	s.collector.OnHTML("li.simpleAds > a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteUrl := e.Request.AbsoluteURL(link)
		detailCollector.Visit(absoluteUrl)
	})

	listings := make([]Ad, 0)

	detailCollector.OnHTML(".itemscope > .left-content", func(e *colly.HTMLElement) {
		requestUrl := e.Request.URL.RequestURI()

		id, err := getListingIdFromUrl(requestUrl)
		if err != nil {
			fmt.Println(err)
			return
		}

		price := e.ChildAttr("meta[itemprop=price]", "content")

		listingEntry := Ad{Id: id, Url: e.Request.AbsoluteURL(requestUrl), Price: price}

		sqFootage := ""
		street := ""
		neighbourhood := ""
		location := ""

		details := e.ChildTexts(".detail")
		for _, txt := range details {
			parsed := strings.Split(txt, ":")
			name := parsed[0]
			value := strings.TrimSpace(parsed[1])
			if strings.Contains(name, "Mikrorajonas") {
				neighbourhood = value
			}

			if strings.Contains(name, "Gatv") {
				street = value
			}

			if strings.Contains(name, "Plotas") {
				sqFootage = value
			}

			if strings.Contains(name, "Gyvenviet") {
				location = value
			}
		}

		listingEntry.SqFootage = sqFootage
		listingEntry.Street = street
		listingEntry.Neighbourhood = neighbourhood
		listingEntry.Location = location

		listings = append(listings, listingEntry)
	})

	// Start scraping
	s.collector.Visit(url)

	return listings
}

func getListingIdFromUrl(url string) (int, error) {
	indexOfLastHyphen := strings.LastIndex(url, "-")

	if indexOfLastHyphen == -1 {
		return 0, fmt.Errorf("could not find the id")
	} else {
		indexOfLastHyphen++
	}

	idWithRemainingDot := url[indexOfLastHyphen:]

	strId := strings.Split(idWithRemainingDot, ".")[0]

	id, err := strconv.Atoi(strId)

	if err != nil {
		return 0, err
	}

	return id, nil
}
