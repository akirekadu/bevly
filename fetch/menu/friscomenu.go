package menu

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/fetch/metadata/frisco"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
	"log"
	"net/url"
	"strings"
)

func init() {
	menuFetcherRegistry["frisco"] = friscoMenu
}

func friscoMenu(provider model.MenuProvider) ([]model.Beverage, error) {
	agent := frisco.Agent()
	response, err := agent.Get(provider.Url())
	if err != nil {
		log.Printf("friscoMenu: Get(%s) failed: %s\n", provider.Url(), err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Printf("friscoMenu: Failed to create document from %s: %s\n", provider.Url(), err)
		return nil, err
	}
	beverages, err := friscoDrafts(doc, response.Request.URL)
	log.Printf("Frisco: parsed %d beverages from %s\n", len(beverages), provider.Url())
	return beverages, err
}

func friscoDrafts(doc *goquery.Document, url *url.URL) ([]model.Beverage, error) {
	beers := []model.Beverage{}
	doc.Find("#drafts a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		beerName := text.Normalize(s.Text())
		if strings.Contains(href, "beer_details") && beerName != "" {
			beer := model.CreateBeverage(beerName)
			hrefUrl, err := url.Parse(href)
			if err == nil {
				beer.SetAttribute(frisco.ProfileURLProperty, hrefUrl.String())
			}
			beers = append(beers, beer)
		}
	})
	if len(beers) == 0 {
		return nil, ErrEmptyMenu
	}
	return beers, nil
}
