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

	collector.OnHTML("#corePriceDisplay_desktop_feature_div", func(e *colly.HTMLElement) {
		inteiro := e.ChildText(".a-price-whole")
		centavo := e.ChildText(".a-price-fraction")

		if inteiro == "" && centavo == "" {
			fmt.Printf("Não foi possível encontrar o preço para URL: %s\n", e.Request.URL.String())
			return
		}

		if centavo == "" {
			centavo = "00"
		}

		if inteiro == "" {
			inteiro = "0"
		}

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

		precoAny := e.Request.Ctx.GetAny("preco")
		if precoAny == nil {
			fmt.Printf("Preço não encontrado para o produto: %s\n", title)
			return
		}

		preco, ok := precoAny.(float64)
		if !ok {
			fmt.Printf("Erro ao converter preço para o produto: %s\n", title)
			return
		}

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
		fmt.Printf("Error ao executar Collector Amazon: %v\n", err)
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
