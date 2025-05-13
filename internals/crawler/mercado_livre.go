package crawler

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"vigil/internals/models"

	"github.com/gocolly/colly/v2"
)

func ScrapeMercadoLivre(urls []string) ([]models.ProductInfo, error) {
	var results []models.ProductInfo

	var mu sync.Mutex
	var wg sync.WaitGroup

	collector := colly.NewCollector(
		colly.AllowedDomains("mercadolivre.com.br", "www.mercadolivre.com.br"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; MercadoLivreCrawler/1.0)"),
	)

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*mercadolivre.com.br*",
		Parallelism: 2,
		RandomDelay: 5 * 1000 * 1000,
	})

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

	collector.OnHTML("div.ui-pdp-price__second-line span[data-testid='price-part']", func(e *colly.HTMLElement) {
		full := e.ChildText("span.andes-money-amount__fraction")
		cents := e.ChildText("span.andes-money-amount__cents")

		priceStr := fmt.Sprintf("%s.%s", full, cents)

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Printf("erro ao converter float ML: %v\n", err)
			return
		}

		e.Request.Ctx.Put("price", price)
	})

	collector.OnHTML("h1.ui-pdp-title", func(e *colly.HTMLElement) {
		title := e.Text

		priceVal := e.Request.Ctx.GetAny("price")
		var price float64

		if p, ok := priceVal.(float64); ok {
			price = p
		}

		mu.Lock()
		results = append(results, models.ProductInfo{
			Url:   e.Request.URL.String(),
			Title: title,
			Price: price,
		})
		mu.Unlock()
	})

	collector.OnScraped(func(r *colly.Response) {
		wg.Done()
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error ao executar Collector Mercado Livre: ", err)
		wg.Done()
	})

	for _, url := range urls {
		wg.Add(1)
		collector.Visit(url)
	}

	collector.Wait()
	wg.Wait()

	return results, nil
}

func FetchNameMercadoLivre(url string) (string, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("mercadolivre.com.br", "www.mercadolivre.com.br"),
		colly.UserAgent("Mozilla/5.0 (compatible; MercadoLivreCrawler/1.0)"),
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

	collector.OnHTML("h1.ui-pdp-title", func(e *colly.HTMLElement) {
		name = strings.TrimSpace(e.Text)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return name, nil
}

func init() {
	Register("www.mercadolivre.com.br", ScrapeMercadoLivre)
	RegisterNameFetcher("www.mercadolivre.com.br", FetchNameMercadoLivre)
}
