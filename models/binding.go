package models

import (
	"fmt"
	"regexp"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/utility"
)

//Binding stores info about product binding ~ category
type Binding struct {
	Model
	Title string
	Count int `gorm:"-"` //product count
}

//GetURL returns the proper binding url
func (b *Binding) GetURL() string {
	slug := utility.Parameterize(b.Title[:utility.Min(30, len(b.Title))])
	if len(slug) == 0 {
		slug = "empty"
	}
	return fmt.Sprintf("/b/%d/%s", b.ID, slug)
}

//Breadcrumbs returns binding's breadcrumbs
func (b *Binding) Breadcrumbs() []Breadcrumb {
	return nil
}

//GetMetaKeywords returns meta keywords
func (b *Binding) GetMetaKeywords() string {
	reg := regexp.MustCompile("[\\[\\]\\(\\)\\{\\}\\.\\,\\!\\?\\:]")
	keywords := reg.ReplaceAllString(b.Title, "")
	return fmt.Sprintf("%s, %s", keywords, config.SiteName)
}

//GetMetaDescription returns meta description
func (b *Binding) GetMetaDescription() (result string) {
	return fmt.Sprintf("Buy %s at the %s gateway shopping centre. Discount and Free Shipping on eligible items.", b.Title, config.SiteName)
}
