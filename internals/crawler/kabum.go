package crawler

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vigil/internals/models"

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

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Sec-Ch-Ua", "\"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

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
			Name:  title,
		})

		mu.Unlock()
	})

	// Encerra a goroutine
	collector.OnScraped(func(r *colly.Response) {
		wg.Done()
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error ao executar Collector Kabum: %v\n", err)
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

func FecthNameKabum(url string) (string, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("kabum.com.br", "www.kabum.com.br"),
		colly.UserAgent("Mozilla/5.0 (compatible; KabumCrawler/1.0)"),
	)

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Sec-Ch-Ua", "\"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	var name string

	collector.OnHTML("h1.sc-58b2114e-6", func(e *colly.HTMLElement) {
		name = strings.TrimSpace(e.Text)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return name, nil
}

func init() {
	Register("www.kabum.com.br", ScrapeKabum)
	RegisterNameFetcher("www.kabum.com.br", FecthNameKabum)
}
