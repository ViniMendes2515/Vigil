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

func ScrapeAmazon(url []string) ([]models.ProductInfo, error) {
	var results []models.ProductInfo

	var mu sync.Mutex
	var wg sync.WaitGroup

	collector := colly.NewCollector(
		colly.AllowedDomains("amazon.com.br", "www.amazon.com.br"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; AmazonCrawler/1.0)"),
	)

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*amazon.*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	collector.OnHTML("#corePriceDisplay_desktop_feature_div", func(e *colly.HTMLElement) {
		inteiro := e.ChildText(".a-price-whole")
		centavo := e.ChildText(".a-price-fraction")

		inteiro = strings.ReplaceAll(inteiro, ".", "")
		inteiro = strings.ReplaceAll(inteiro, ",", "")
		inteiro = strings.TrimSpace(inteiro)
		centavo = strings.TrimSpace(centavo)

		precoStr := fmt.Sprintf("%s.%s", inteiro, centavo)

		preco, err := strconv.ParseFloat(precoStr, 64)
		if err != nil {
			fmt.Printf("Erro ao converter preço: %v\n", err)
			return
		}

		e.Request.Ctx.Put("preco", preco)
	})

	collector.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		preco := e.Request.Ctx.GetAny("preco").(float64)

		mu.Lock()

		results = append(results, models.ProductInfo{
			Url:   e.Request.URL.String(),
			Title: title,
			Price: preco,
		})

		mu.Unlock()
	})

	collector.OnScraped(func(r *colly.Response) {
		wg.Done()
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error ao executar Collector: %v\n", err)
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

func FetchNameAmazon(url string) (string, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("amazon.com.br", "www.amazon.com.br"),
		colly.UserAgent("Mozilla/5.0 (compatible; AmazonCrawler/1.0)"),
	)

	var name string

	collector.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
		name = strings.TrimSpace(e.Text)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return name, nil
}

func init() {
	Register("www.amazon.com.br", ScrapeAmazon)
	RegisterNameFetcher("www.amazon.com.br", FetchNameAmazon)
}
