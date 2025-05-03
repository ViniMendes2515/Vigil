package crawler

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ViniMendes2515/price-crawler/internals/models"
	"github.com/gocolly/colly/v2"
)

// ScrapeKabum faz o scraping dos produtos da Kabum e retorna uma lista de ProductInfo
func ScrapeKabum(url []string) ([]models.ProductInfo, error) {
	var results []models.ProductInfo

	var mu sync.Mutex
	var wg sync.WaitGroup
	collector := colly.NewCollector(
		colly.AllowedDomains("kabum.com.br", "www.kabum.com.br"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; KabumCrawler/1.0)"),
	)

	// Limit requests
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*kabum.*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	collector.OnHTML("h4.finalPrice", func(e *colly.HTMLElement) {
		raw := strings.TrimSpace(e.Text)
		raw = strings.ReplaceAll(raw, "R$", "")
		raw = strings.ReplaceAll(raw, ".", "")
		raw = strings.ReplaceAll(raw, ",", ".")

		var price float64
		fmt.Sscanf(raw, "%f", &price)

		e.Request.Ctx.Put("price", fmt.Sprintf("%f", price))
	})

	collector.OnHTML("h1.sc-58b2114e-6", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		priceStr := e.Request.Ctx.Get("price")

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Printf("Error parsing price: %v\n", err)
			wg.Done()
			return
		}

		mu.Lock()

		results = append(results, models.ProductInfo{
			Url:   e.Request.URL.String(),
			Title: title,
			Price: price,
		})

		mu.Unlock()
	})

	// Encerra a goroutine
	collector.OnScraped(func(r *colly.Response) {
		wg.Done()
	})

	for _, url := range url {
		wg.Add(1)
		collector.Visit(url)
	}

	collector.Wait()
	wg.Wait()

	return results, nil
}
