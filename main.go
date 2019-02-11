package main

import (
	"log"

	"github.com/denisbakhtin/amazon-go/aws"
	"github.com/denisbakhtin/amazon-go/cache"
	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/controllers"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/sitemap"
	"github.com/robfig/cron"
)

func main() {
	//------------------------- initialization ---------------------------
	config.Init()
	defer config.DisposeLogger()
	cache.Init()
	controllers.InitRoutes()
	models.InitDatabase() //default settings will be used

	//------------------------- service tasks ---------------------------
	//go aws.DealsParse()
	//go aws.RssParse()
	//go aws.FeedParse()
	go func() {
		//aws.ClearProcessedAsins()
		//aws.QueueAvailableAsins()
		//aws.QueueUnavailableAsins()
		aws.ProductUpdate()
	}()

	if config.IsRelease() {
		//------------------------- scheduled tasks --------------------------
		//TODO: later move these tasks to system cron, to make it more stable and independent of app uptime
		c := cron.New()
		c.AddFunc("0 0 0 * * *", func() {
			err := sitemap.Create()
			if err != nil {
				log.Println(err)
			}
		}) //every day
		c.AddFunc("23 50 0 * * *", func() { aws.MaintainDB() })                 //every day at 23:50
		c.AddFunc("0 2 0 * * *", func() { aws.ClearProcessedAsins() })          //every day at 00:02
		c.AddFunc("0 3 0 * * *", func() { aws.ClearProcessedSpecifications() }) //every day at 00:03 //obsolete
		c.AddFunc("0 5 0 * * *", func() { aws.QueueAvailableAsins() })          //every day at 00:05
		c.AddFunc("0 7 0 * * *", func() { models.InitializeCache() })           //every day at 00:07
		c.AddFunc("0 10 0 4 * *", func() { aws.QueueUnavailableAsins() })       //every 4 days at 00:10
		c.Start()

		//launch tasks in at the start of application in release mode
		go models.InitializeCache()
		go aws.RssParse()
		go aws.FeedParse()
		go aws.ProductUpdate()
		go aws.SpecificationUpdate()
	}

	//------------------------- run web server ---------------------------
	log.Fatal(controllers.GetRouter().Run(":9000"))
}
