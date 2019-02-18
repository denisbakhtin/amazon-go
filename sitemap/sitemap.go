package sitemap

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/utility"
)

const (
	header = `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	footer   = ` </urlset>`
	template = `
	 <url>
	   <loc>%s</loc>
	   <lastmod>%s</lastmod>
	   <changefreq>%s</changefreq>
	   <priority>%.1f</priority>
	 </url> 	`

	indexHeader = `<?xml version="1.0" encoding="UTF-8"?>
      <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	indexFooter = `
</sitemapindex>
	`
	indexTemplate = `
    <sitemap>
       <loc>%s%s</loc>
       <lastmod>%s</lastmod>
    </sitemap>
	`
	linkLimit = 50000
)

type item struct {
	Loc        string
	LastMod    time.Time
	ChangeFreq string
	Priority   float32
}

func (item *item) String() string {
	//2012-08-30T01:23:57+08:00
	//Mon Jan 2 15:04:05 -0700 MST 2006
	return fmt.Sprintf(template, item.Loc, item.LastMod.Format("2006-01-02T15:04:05+08:00"), item.ChangeFreq, item.Priority)
}

//Create creates sitemap files
func Create() error {
	//firstly delete old files, then create fresh ones
	folder := config.PublicPath + "/system"
	fs, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	log.Println("Starting sitemap creation process")

	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".xml.gz") {
			log.Printf("Deleting old sitemap file: %s\n", f.Name())
			os.Remove(folder + "/" + f.Name())
		}
	}

	baseURL := config.SiteURL
	items := make([]item, 0, 100000)
	var nodes []models.BrowseNode
	models.DB.Where("product_count > 0").Find(&nodes)
	for i := range nodes {
		items = append(items, item{Loc: baseURL + nodes[i].GetURL(), LastMod: nodes[i].UpdatedAt, ChangeFreq: "daily", Priority: 0.5})
	}

	var products []models.Product
	models.DB.Select("id, title, created_at, category_id, browse_node_id").Find(&products)
	for i := range products {
		items = append(items, item{Loc: baseURL + products[i].GetURL(), LastMod: products[i].UpdatedAt, ChangeFreq: "weekly", Priority: 0.9})
	}

	var pages []models.Page
	models.DB.Where("show = true").Find(&pages)
	for i := range pages {
		items = append(items, item{Loc: baseURL + pages[i].GetURL(), LastMod: pages[i].UpdatedAt, ChangeFreq: "monthly", Priority: 0.4})
	}

	//split items by linkLimit and write to separate files
	for i := 0; i < int(math.Ceil(float64(len(items))/linkLimit)); i++ {
		fileName := fmt.Sprintf("%s/%s%d.xml.gz", folder, "sitemap", i)
		err := createFile(fileName, items[i*linkLimit:utility.Min((i+1)*linkLimit, len(items))])
		if err != nil {
			return err
		}
	}

	err = createIndex(folder, folder+"/sitemap.xml.gz", baseURL+"/system/")
	if err != nil {
		return err
	}

	log.Println("Stopping sitemap creation process")
	return nil
}

//create sitemap file (with max links of 50000)
func createFile(f string, items []item) error {
	log.Printf("Creating sitemap file: %s\n", f)
	var buffer bytes.Buffer
	buffer.WriteString(header)
	for _, item := range items {
		_, err := buffer.WriteString(item.String())
		if err != nil {
			return err
		}
	}
	fo, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fo.Close()
	buffer.WriteString(footer)

	zip := gzip.NewWriter(fo)
	defer zip.Close()
	_, err = zip.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

//create index sitemap file (with a list of detailed ones)
func createIndex(folder, indexFile, baseurl string) error {
	log.Printf("Creating sitemap index file: %s\n", indexFile)
	var buffer bytes.Buffer
	buffer.WriteString(indexHeader)
	fs, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	fo, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	defer fo.Close()

	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".xml.gz") {
			s := fmt.Sprintf(indexTemplate, baseurl, f.Name(), time.Now().Format("2006-01-02T15:04:05+08:00"))
			buffer.WriteString(s)
		}
	}
	buffer.WriteString(indexFooter)

	zip := gzip.NewWriter(fo)
	defer zip.Close()
	_, err = zip.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}
