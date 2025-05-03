package scheduler

import (
	"fmt"

	"github.com/ViniMendes2515/price-crawler/config"
	"github.com/ViniMendes2515/price-crawler/internals/crawler"
	"github.com/ViniMendes2515/price-crawler/internals/models"
	"github.com/ViniMendes2515/price-crawler/internals/notifier"
	"github.com/ViniMendes2515/price-crawler/internals/services"
	"github.com/robfig/cron/v3"
)

func BuscaPrecos(c *cron.Cron, cfg config.Config, tg *notifier.TelegramNotifier) {

	c.AddFunc("0 9,15,21 * * *", func() {

		var produtos []models.ProductInfo

		if len(cfg.KabumURLs) > 0 {
			produtosKabum, err := crawler.ScrapeKabum(cfg.KabumURLs)
			if err != nil {
				fmt.Println("Erro ao fazer o scraping Kabum: ", err)
				return
			}
			produtos = append(produtos, produtosKabum...)
		}

		if len(cfg.AmazonURLs) > 0 {
			produtosAmazon, err := crawler.ScrapeAmazon(cfg.AmazonURLs)
			if err != nil {
				fmt.Println("Erro ao fazer o scraping Amazon: ", err)
				return
			}
			produtos = append(produtos, produtosAmazon...)
		}

		services.Monitorar(produtos, *tg)
	})

}
