package models

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/denisbakhtin/amazon-go/utility"
	"github.com/denisbakhtin/amazon-go/viewmodels"
)

//Category stores info about product category
type Category struct {
	Model
	Title        string
	ParentID     *uint64 `gorm:"index:category_parent_idx"`
	Description  string
	Image        string
	Children     []Category
	Products     []Product
	ProductCount int64 `sql:"-"` //product count cache
}

//BeforeSave gorm hook
func (c *Category) BeforeSave() error {
	if len(c.Title) == 0 {
		return errors.New("Title is empty")
	}
	if c.ParentID != nil {
		//check parent exists and its top level
		parent := Category{}
		DB.Where(*c.ParentID).First(&parent)
		if parent.ID == 0 || parent.ParentID != nil {
			return errors.New("Parent category does not exist or can't be selected")
		}
	}
	if len(c.Description) > 5000 {
		return errors.New("Description length exceeds 5000 character limit")
	}
	return nil
}

//GetURL returns the proper category url
func (c *Category) GetURL() string {
	return fmt.Sprintf("/categories/%d/%s", c.ID, utility.Parameterize(c.Title))
}

//GetCategoryIDs returs self.ID + all Children's ids
func (c *Category) GetCategoryIDs() []uint64 {
	result := make([]uint64, len(c.Children)+1)
	result[0] = c.ID
	for i := range c.Children {
		result[i+1] = c.Children[i].ID
	}
	return result
}

//MainImage returns category image url
func (c *Category) MainImage() string {
	switch {
	case len(c.Image) > 0:
		return "/uploads/" + c.Image
	default:
		return "/images/no-image.jpg"
	}
}

//GetDescription returns category description
func (c *Category) GetDescription() template.HTML {
	meta := ""
	Cache.RWMutex.RLock()
	if len(Cache.CategoryDescriptions) > int(c.ID) {
		meta = Cache.CategoryDescriptions[c.ID]
		Cache.RWMutex.RUnlock()
	} else {
		Cache.RWMutex.RUnlock()   //unlock to prevent deadlock
		meta = c.SetDescription() //write lock inside
	}
	return template.HTML(meta)
}

//SetDescription creates and caches category description
func (c *Category) SetDescription() string {
	Cache.RWMutex.Lock()
	Cache.CategoryDescriptions = utility.AppendToCache(Cache.CategoryDescriptions, int(c.ID), c.compileDescription())
	meta := Cache.CategoryDescriptions[c.ID]
	Cache.RWMutex.Unlock()
	return meta
}

//compileDescription creates category description via pattern
func (c *Category) compileDescription() string {
	result := ""

	c.LoadChildren()
	categoryIDs := c.GetCategoryIDs()

	var totalProductCount int
	DB.Model(Product{}).Where("category_id IN(?)", categoryIDs).Count(&totalProductCount)

	var products []Product
	DB.Select("brand_id, count(id)").Where("category_id IN(?)", categoryIDs).Group("brand_id").Order("count(id) desc").Limit(10).Find(&products)
	brandIDs := make([]uint64, 0, len(products))
	for i := range products {
		if products[i].BrandID != nil {
			brandIDs = append(brandIDs, *products[i].BrandID)
		}
	}

	var brands []Brand
	DB.Where("id IN(?)", brandIDs).Select("title").Find(&brands)
	brandTitles := make([]string, len(brands))
	for i := range brands {
		brandTitles[i] = brands[i].Title
	}

	/* var tags []Tag
	DB.Where("category_id IN(?)", categoryIDs).Select("title").Limit(5).Find(&tags)
	tagTitles := make([]string, len(tags))
	for i := range tags {
		tagTitles[i] = tags[i].Title
	}

	if len(tags) > 0 {
		result = fmt.Sprintf("<p>Welcome to SaveLikea.Pro! Browse a list of the most popular %s, where %d items from %s to %s are waiting for you. Feel free to use search, if you are looking for specific %s.</p>", c.Title, totalProductCount, tags[0].Title, tags[len(tags)-1].Title, c.Title)
	} else {
		result = fmt.Sprintf("<p>Welcome to SaveLikea.Pro! Browse a list of the most popular %s, where %d items are waiting for you. Feel free to use search, if you are looking for specific %s.</p>", c.Title, totalProductCount, c.Title)
	} */
	result = result + fmt.Sprintf("<p>Many of our customers have become very brand savvy, and for them we carry %s, just to mention a few. <span class='text-success'>Get label quality without having to pay a label price.</span></p>", strings.Join(brandTitles, ", "))

	freeDeliveryCount := 0
	DB.Model(Product{}).Where("category_id IN(?) and free_shipping = ?", categoryIDs, true).Count(&freeDeliveryCount)
	if freeDeliveryCount > 0 {
		result = result + fmt.Sprintf("<p>Items that comply to Free Shipping (currently %d), for your convenience, are marked with <span class='glyphicon glyphicon-gift'></span><span class='glyphicon glyphicon-send'></span> icons.</p>", freeDeliveryCount)
	} else {
		result = result + fmt.Sprintf("<p>Items that comply to Free Shipping, for your convenience, are marked with <span class='glyphicon glyphicon-gift'></span><span class='glyphicon glyphicon-send'></span> icons.</p>")
	}
	result = result + fmt.Sprintf("<h3>Why SaveLikea.Pro?</h3> <strong>Because we offer best prices on %s on the Internet</strong>. These prices are real, believe it. ", c.Title)
	/* if len(tags) > 0 {
		result = result + fmt.Sprintf("Buy wide range of <strong>%s</strong>, choosing from %s etc, and enjoy the ease and speed of online shopping. ", c.Title, strings.Join(tagTitles, ", "))
	} else {
		result = result + fmt.Sprintf("Buy wide range of <strong>%s</strong> and enjoy the ease and speed of online shopping. ", strings.ToLower(c.Title))
	} */
	result = result + "Save your money and time like a Pro."

	return result
}

//Breadcrumbs returns category breadcrumbs
func (c *Category) Breadcrumbs() []viewmodels.Breadcrumb {
	return nil
}

//LoadChildren loads category children
func (c *Category) LoadChildren() {
	if len(c.Children) == 0 {
		DB.Where("parent_id = ?", c.ID).Find(&c.Children)
	}
}

//GetMetaDescription returns category meta description
func (c *Category) GetMetaDescription() string {
	meta := ""
	Cache.RWMutex.RLock()
	if len(Cache.CategoryMetaDescriptions) > int(c.ID) {
		meta = Cache.CategoryMetaDescriptions[c.ID]
		Cache.RWMutex.RUnlock()
	} else {
		Cache.RWMutex.RUnlock()       //unlock to prevent deadlock
		meta = c.SetMetaDescription() //write lock inside
	}
	return meta
}

//SetMetaDescription creates and caches category meta description
func (c *Category) SetMetaDescription() string {
	Cache.RWMutex.Lock()
	Cache.CategoryMetaDescriptions = utility.AppendToCache(Cache.CategoryMetaDescriptions, int(c.ID), c.compileMetaDescription())
	meta := Cache.CategoryMetaDescriptions[c.ID]
	Cache.RWMutex.Unlock()
	return meta
}

//compileMetaDescription creates meta description following the pattern
func (c *Category) compileMetaDescription() string {
	result := ""
	c.LoadChildren()
	categoryIDs := make([]uint64, 1+len(c.Children))
	categoryIDs[0] = c.ID
	for i := range c.Children {
		categoryIDs[i+1] = c.Children[i].ID
	}

	var totalProductCount int
	DB.Model(Product{}).Where("category_id IN(?)", categoryIDs).Count(&totalProductCount)
	if totalProductCount > 0 {
		variation := Variation{}
		DB.Where("category_id IN(?)", categoryIDs).
			Order("discount_percent desc").
			Select("discount_percent,discount").
			First(&variation)
		result = fmt.Sprintf("Save up to %.1f%% ($%.1f) on %s purchase. Free shipping, daily updates...", variation.DiscountPercent, variation.Discount, strings.ToLower(c.Title))
	} else {
		result = fmt.Sprintf("Explore %s department and find out the best way to save substancial amount of money on every purchase. Limited time offer, daily updates...", strings.ToLower(c.Title))
	}
	return result
}

//GetMetaKeywords returns category meta keywords
func (c *Category) GetMetaKeywords() string {
	meta := ""
	Cache.RWMutex.RLock()
	if len(Cache.CategoryMetaKeywords) > int(c.ID) {
		meta = Cache.CategoryMetaKeywords[c.ID]
		Cache.RWMutex.RUnlock()
	} else {
		Cache.RWMutex.RUnlock()    //unlock to prevent deadlock
		meta = c.SetMetaKeywords() //write lock inside
	}
	return meta
}

//SetMetaKeywords creates and caches category meta keywords
func (c *Category) SetMetaKeywords() string {
	Cache.RWMutex.Lock()
	Cache.CategoryMetaKeywords = utility.AppendToCache(Cache.CategoryMetaKeywords, int(c.ID), c.compileMetaKeywords())
	meta := Cache.CategoryMetaKeywords[c.ID]
	Cache.RWMutex.Unlock()
	return meta
}

//compileMetaKeywords creates meta keywords
func (c *Category) compileMetaKeywords() string {
	keywords := []string{c.Title, fmt.Sprintf("%s deals", c.Title), fmt.Sprintf("%s discounts", c.Title), fmt.Sprintf("%s mark down", c.Title), fmt.Sprintf("cheap %s", c.Title)}
	if c.ParentID != nil {
		parent := Category{}
		DB.First(&parent, *c.ParentID)
		keywords = append(keywords, fmt.Sprintf("%s %s", c.Title, parent.Title))
	}

	return strings.Join(keywords, ", ")
}

/*
//GetBrandScale returns tag brands
func (t *Tag) GetBrandScale(current string) (result []Scale) {
	var brands []Brand
	DB.Model(Brand{}).Where("id IN(select brand_id from products where tag_id =? group by brand_id limit 100)", t.ID).Find(&brands)
	result = make([]Scale, len(brands)+1)
	tagURL := t.GetURL()
	result[0] = Scale{Title: "All", URL: tagURL}
	currentID, _ := strconv.ParseInt(current, 10, 64)
	for i := 0; i < len(brands); i++ {
		scale := Scale{Title: brands[i].Title, URL: fmt.Sprintf("%s?brand=%d", tagURL, brands[i].ID)}
		if brands[i].ID == uint64(currentID) {
			scale.Class = "active"
		}
		result[i+1] = scale
	}
	return
}

//GetPriceScale returns tag price scales, current param is obsolete
func (t *Tag) GetPriceScale(current string) []Scale {
	strScale := ""
	Cache.RWMutex.RLock()
	if len(Cache.TagPriceScale) > int(t.ID) {
		strScale = Cache.TagPriceScale[t.ID]
		Cache.RWMutex.RUnlock()
	} else {
		Cache.RWMutex.RUnlock()             //unlock to prevent deadlock
		strScale = t.SetPriceScale(current) //write lock inside
	}

	scales := make([]Scale, 0)
	json.Unmarshal([]byte(strScale), &scales) //ignore error, if cant unmarshal (e.g empty string)
	return scales
}

//SetPriceScale calculates and caches tag price scale
func (t *Tag) SetPriceScale(current string) string {
	Cache.RWMutex.Lock()
	scales := t.compilePriceScale(current)
	strScale, _ := json.Marshal(scales) //ignore error, not crucial
	Cache.TagPriceScale = utility.AppendToCache(Cache.TagPriceScale, int(t.ID), string(strScale))
	Cache.RWMutex.Unlock()

	return string(strScale)
}

//compilePriceScale calculates tag price scales
func (t *Tag) compilePriceScale(current string) (result []Scale) {
	var priceRanges []PriceRange
	DB.Model(PriceRange{}).Where("EXISTS (SELECT NULL FROM variations v WHERE v.special_price >= price_ranges.from AND v.special_price < price_ranges.to AND v.tag_id = ? LIMIT 1)", t.ID).Order("\"from\"").Find(&priceRanges)
	result = make([]Scale, len(priceRanges))
	tagURL := t.GetURL()
	for i := 0; i < len(priceRanges); i++ {
		var scale Scale
		if priceRanges[i].Code == "all" {
			scale = Scale{Title: priceRanges[i].Title, URL: tagURL}
		} else {
			scale = Scale{Title: priceRanges[i].Title, URL: tagURL + "?price=" + priceRanges[i].Code}
		}
		if priceRanges[i].Code == current {
			scale.Class = "active"
		}
		result[i] = scale
	}
	return
}

//GetDiscountScale returns tag discount scales, current param is obsolete
func (t *Tag) GetDiscountScale(current string) []Scale {
	strScale := ""
	Cache.RWMutex.RLock()
	if len(Cache.TagDiscountScale) > int(t.ID) {
		strScale = Cache.TagDiscountScale[t.ID]
		Cache.RWMutex.RUnlock()
	} else {
		Cache.RWMutex.RUnlock()                //unlock to prevent deadlock
		strScale = t.SetDiscountScale(current) //write lock inside
	}

	scales := make([]Scale, 0)
	json.Unmarshal([]byte(strScale), &scales) //ignore error, if cant unmarshal (e.g empty string)
	return scales
}

//SetDiscountScale calculates and caches discount scales
func (t *Tag) SetDiscountScale(current string) string {
	Cache.RWMutex.Lock()
	scales := t.compileDiscountScale(current)
	strScale, _ := json.Marshal(scales) //ignore error, not crucial
	Cache.TagDiscountScale = utility.AppendToCache(Cache.TagDiscountScale, int(t.ID), string(strScale))
	Cache.RWMutex.Unlock()

	return string(strScale)
}

//compileDiscountScale calculates discount scales
func (t *Tag) compileDiscountScale(current string) (result []Scale) {
	var discountRanges []DiscountRange
	DB.Model(DiscountRange{}).Where("EXISTS (SELECT NULL FROM variations v WHERE v.discount_percent >= discount_ranges.from AND v.discount_percent < discount_ranges.to AND v.tag_id = ? LIMIT 1)", t.ID).Order("\"from\"").Find(&discountRanges)
	result = make([]Scale, len(discountRanges))
	tagURL := t.GetURL()
	for i := 0; i < len(discountRanges); i++ {
		var scale Scale
		if discountRanges[i].Code == "all" {
			scale = Scale{Title: discountRanges[i].Title, URL: tagURL}
		} else {
			scale = Scale{Title: discountRanges[i].Title, URL: tagURL + "?discount=" + discountRanges[i].Code}
		}
		if discountRanges[i].Code == current {
			scale.Class = "active"
		}
		result[i] = scale
	}
	return
}

*/
