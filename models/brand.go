package models

import (
	"fmt"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/utility"
)

//Brand stores info about product brand
type Brand struct {
	Model
	Title string
	Body  string
	Count int `gorm:"-"` //product count
}

//GetURL returns the proper brand url
func (b *Brand) GetURL() string {
	slug := utility.Parameterize(b.Title[:utility.Min(30, len(b.Title))])
	if len(slug) == 0 {
		slug = "empty"
	}
	return fmt.Sprintf("/brands/%d/%s", b.ID, slug)
}

//Breadcrumbs returns brand's breadcrumbs
func (b *Brand) Breadcrumbs() []Breadcrumb {
	return nil
}

//IDStr returns string representation of id
func (b *Brand) IDStr() string {
	return fmt.Sprintf("%d", b.ID)
}

//GetMetaKeywords returns meta keywords
func (b *Brand) GetMetaKeywords() string {
	return fmt.Sprintf("%s, %s", b.Title, config.SiteName)
}

//GetMetaDescription returns meta description
func (b *Brand) GetMetaDescription() (result string) {
	return fmt.Sprintf("Buy %s products at the %s gateway shopping centre. Discount and Free Shipping on eligible items.", b.Title, config.SiteName)
}
