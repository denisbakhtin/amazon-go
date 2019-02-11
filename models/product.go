package models

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/utility"
)

//Product stores info about product
type Product struct {
	Model
	Asin               string
	Available          bool `gorm:"index:product_available_idx"`
	RegularPrice       float64
	SpecialPrice       float64
	Discount           float64
	DiscountPercent    float64
	Title              string
	URL                string
	CustomerReviews    string
	Content            string
	CategoryID         uint64  `gorm:"index:product_category_idx"`
	FeedID             *uint64 `gorm:"index:product_feed_idx"`
	BrowseNodeID       uint64  `gorm:"index:product_browse_node_idx"`
	FreeShipping       bool
	BrandID            *uint64 `gorm:"index:product_brand_idx"`
	ProductGroupTypeID *uint64 `gorm:"index:product_group_type_idx"`
	Rate               float64
	Image              string
	RedirectToAsin     string
	AccountID          uint64  `gorm:"index:product_account_idx"`
	CompanyID          *uint64 `gorm:"index:product_company_idx"`
	CreatedByUser      bool
	UserImage          string
	LanguageID         uint64  `gorm:"index:product_language_idx"`
	BindingID          *uint64 `gorm:"index:product_binding_idx"`
	DepartmentID       *uint64 `gorm:"index:product_department_idx"`
	Variations         []Variation
	Category           Category
	BrowseNode         BrowseNode
	Brand              Brand      `gorm:"association_autoupdate:false;association_autocreate:false"`
	Binding            Binding    `gorm:"association_autoupdate:false;association_autocreate:false"`
	Department         Department `gorm:"association_autoupdate:false;association_autocreate:false"`
}

//MainImage returns product main image url
func (p *Product) MainImage() string {
	switch {
	case len(p.Image) > 0:
		return p.Image
	default:
		return "/images/no-image.jpg"
	}
}

//TitleWithoutBrand returns product title without brand name
func (p *Product) TitleWithoutBrand() string {
	if len(p.Brand.Title) > 0 && strings.HasPrefix(p.Title, p.Brand.Title) {
		return strings.Trim(strings.TrimPrefix(p.Title, p.Brand.Title), " ")
	}
	return p.Title
}

//ImageTitle returns product image title
func (p *Product) ImageTitle() string {
	return strings.ToLower(p.Title) + " image"
}

//DoShowReviews checks if reviews should be visible
func (p *Product) DoShowReviews() bool {
	return p.Available && len(p.CustomerReviews) > 0
}

//ActualReviewsURL returns a valid up to date reviews url
func (p *Product) ActualReviewsURL() string {
	//check if original link has not expired
	exp := regexp.MustCompile("exp=([0-9]{4}-[0-9]{2}-[0-9]{2})")
	rdate := exp.FindStringSubmatch(p.CustomerReviews)
	if len(rdate) > 1 {
		if date, err := time.Parse("2006-01-02", rdate[1]); err == nil {
			serverBod := utility.BeginningOfDay(utility.AmazonTime(time.Now()))
			pageBod := utility.BeginningOfDay(utility.AmazonTime(date))
			//leave untouched
			if pageBod.Equal(serverBod) || pageBod.After(serverBod) {
				return p.CustomerReviews
			}
		}
	}
	tomorrow := time.Now().Add(24 * time.Hour)
	url := exp.ReplaceAllString(p.CustomerReviews, fmt.Sprintf("exp=%04d-%02d-%02d", tomorrow.Year(), tomorrow.Month(), tomorrow.Day()))
	return url
}

//PromoVariation returns variation with maximum discount %
func (p *Product) PromoVariation() Variation {
	if len(p.Variations) == 0 {
		return Variation{}
	}
	//variations should be sorted by discout % desc
	return p.Variations[0]
}

//RatingSlice builds star rating slice for product view
func (p *Product) RatingSlice() (starRating [5]bool) {
	for i := 0; i < 5; i++ {
		starRating[i] = (math.Ceil(p.Rate) >= float64(i+1))
	}
	return
}

//Similar builds the slice of similar products
func (p *Product) Similar() (similar []Product) {
	similar = make([]Product, 0, config.SimilarLimit)
	DB.Where("id != ? and available = true and browse_node_id = ?", p.ID, p.BrowseNodeID).
		Preload("BrowseNode").
		Preload("Brand").
		Limit(config.SimilarLimit).
		Find(&similar)
	if len(similar) < config.SimilarLimit {
		var similar2 []Product
		DB.Where("id != ? and brand_id = ? and available = true and browse_node_id != ?", p.ID, p.BrandID, p.BrowseNodeID).
			Preload("BrowseNode").
			Preload("Brand").
			Limit(config.SimilarLimit - len(similar)).
			Find(&similar2)
		similar = append(similar, similar2...)
	}
	if len(similar) < config.SimilarLimit {
		var similar2 []Product
		DB.Where("id != ? and brand_id != ? and category_id = ? and available = true and browse_node_id != ?", p.ID, p.BrandID, p.CategoryID, p.BrowseNodeID).
			Preload("BrowseNode").
			Preload("Brand").
			Limit(config.SimilarLimit - len(similar)).
			Find(&similar2)
		similar = append(similar, similar2...)
	}
	if len(similar) < config.SimilarLimit {
		var similar2 []Product
		DB.Where("id != ? and brand_id != ? and category_id != ? and available = true and browse_node_id != ?", p.ID, p.BrandID, p.CategoryID, p.BrowseNodeID).
			Preload("BrowseNode").
			Preload("Brand").
			Limit(config.SimilarLimit - len(similar)).
			Find(&similar2)
		similar = append(similar, similar2...)
	}
	return
}

//GetURL returns the proper product url
func (p *Product) GetURL() string {
	return fmt.Sprintf("/products/%d/%s", p.ID, utility.Parameterize(p.Title[:utility.Min(30, len(p.Title))]))
}

//PriceColumns returns the array of product prices
func (p *Product) PriceColumns() (columns []Column) {
	if len(p.Variations) == 0 {
		return
	}
	var dims []Dimension
	DB.Where(p.getDimIDs()).Select("id, name, title").Find(&dims)
	columns = make([]Column, len(dims)+4) //+ static columns, see below

	columns[0] = Column{ID: "asin", Title: "Asin"}
	i := 1
	for _, d := range dims {
		columns[i] = Column{ID: fmt.Sprintf("%d", d.ID), Title: d.GetTitle()}
		i++
	}
	columns[i] = Column{ID: "special_price", Title: "Special price"}
	i++
	columns[i] = Column{ID: "regular_price", Title: "Regular price"}
	i++
	columns[i] = Column{ID: "discount_percent", Title: "Discount"}
	return
}

func (p *Product) getDimIDs() []uint64 {
	result := make([]uint64, 0, 5)
	if len(p.Variations) == 0 {
		return result
	}
	v := &p.Variations[0]
	if v.Dim1Id != nil {
		result = append(result, *v.Dim1Id)
	}
	if v.Dim2Id != nil {
		result = append(result, *v.Dim2Id)
	}
	if v.Dim3Id != nil {
		result = append(result, *v.Dim3Id)
	}
	if v.Dim4Id != nil {
		result = append(result, *v.Dim4Id)
	}
	if v.Dim5Id != nil {
		result = append(result, *v.Dim5Id)
	}
	return result
}

//SelectableDims return a slice of product dimensions suitable for selects (with more than 1 unique value)
func (p *Product) SelectableDims() []Dimension {
	var dims []Dimension
	DB.Where(p.getDimIDs()).Select("id, name, title").Order("id asc").Find(&dims)
	result := make([]Dimension, 0, len(dims))
	for i := range dims {
		uv := p.UniquePriceColumnData(dims[i].IDStr(), p.PriceData())
		if len(uv) > 1 {
			result = append(result, dims[i])
		}
	}
	return result
}

//PriceData returns product prices by dimensions
func (p *Product) PriceData() (dimData []map[string]string) {
	if len(p.Variations) == 0 {
		return
	}
	dimData = make([]map[string]string, len(p.Variations))
	for i, v := range p.Variations {
		dimMap := make(map[string]string)
		if v.Dim1Id != nil {
			dimMap[fmt.Sprintf("%d", *v.Dim1Id)] = v.Dim1Value
		}
		if v.Dim2Id != nil {
			dimMap[fmt.Sprintf("%d", *v.Dim2Id)] = v.Dim2Value
		}
		if v.Dim3Id != nil {
			dimMap[fmt.Sprintf("%d", *v.Dim3Id)] = v.Dim3Value
		}
		if v.Dim4Id != nil {
			dimMap[fmt.Sprintf("%d", *v.Dim4Id)] = v.Dim4Value
		}
		if v.Dim5Id != nil {
			dimMap[fmt.Sprintf("%d", *v.Dim5Id)] = v.Dim5Value
		}
		dimMap["asin"] = v.Asin
		dimMap["special_price"] = fmt.Sprintf("$%.2f", v.SpecialPrice)
		dimMap["regular_price"] = fmt.Sprintf("$%.2f", v.RegularPrice)
		dimMap["discount_percent"] = fmt.Sprintf("%.1f%%", v.DiscountPercent)
		dimData[i] = dimMap
	}
	return
}

//PriceColumnData returns price column data by id
func (p *Product) PriceColumnData(id string, dimData []map[string]string) []string {
	result := make([]string, len(dimData))
	for i := range dimData {
		result[i] = dimData[i][id]
	}
	return result
}

//UniquePriceColumnData returns price column data by id
func (p *Product) UniquePriceColumnData(id string, dimData []map[string]string) []string {
	data := p.PriceColumnData(id, dimData)
	return utility.UniqueStrings(data)
}

//Breadcrumbs returns product breadcrumbs
func (p *Product) Breadcrumbs() []Breadcrumb {
	if !p.BrowseNode.AllParentsLoaded {
		p.BrowseNode.LoadAllParents()
	}
	return p.buildBreadcrumbs()
}

func (p *Product) buildBreadcrumbs() []Breadcrumb {
	crumbs := make([]Breadcrumb, 0, 20)
	if p.BrowseNode.ID != 0 {
		crumbs = append(p.BrowseNode.buildBreadcrumbs(), Breadcrumb{URL: p.BrowseNode.GetURL(), Title: p.BrowseNode.Title})
	}
	return crumbs
}

//GetMetaKeywords returns meta keywords
func (p *Product) GetMetaKeywords() string {
	reg := regexp.MustCompile("[\\[\\]\\(\\)\\{\\}\\.\\,\\!\\?\\:]")
	keywords := strings.Split(reg.ReplaceAllString(p.Title, ""), " ")
	keywords = utility.SubtractArray(keywords, []string{"for", "and", "by", "a", "the", "this", "to", "from", "on", "under", ""})
	return fmt.Sprintf("%s, %s", strings.ToLower(strings.Join(keywords, " ")), config.SiteName)
}

//GetMetaDescription returns meta description
func (p *Product) GetMetaDescription() (result string) {
	tag := ""
	if p.BrowseNode.ID > 0 {
		parent := p.BrowseNode.TopParent()
		tag = parent.Title
	} else if len(p.Binding.Title) > 0 {
		tag = p.Binding.Title
	} else {
		tag = config.SiteName
	}
	brand := ""
	if len(p.Brand.Title) > 0 {
		brand = p.Brand.Title
	} else {
		brand = p.Title
	}
	return fmt.Sprintf("Buy %s at the %s gateway shopping centre. Discount and Free Shipping on eligible items.", brand, tag)
}

//DimValuesJS returns a string with two dimensional js array
func (p *Product) DimValuesJS(dims []Dimension) string {
	prices := p.PriceData()
	result := "["
	for i := range prices {
		result += "["
		for j := range dims {
			result += fmt.Sprintf("%q,", prices[i][dims[j].IDStr()])
		}
		result += fmt.Sprintf("%q", prices[i]["asin"])
		result += "]"
		if i < len(prices)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}

//WithSameBrand returns a slice of 'limit' products with the same brand
func (p *Product) WithSameBrand(limit int) []Product {
	var result []Product
	if p.BrandID == nil {
		return result
	}
	DB.Preload("Brand").Preload("BrowseNode").Where("brand_id = ? and id != ? and available = true", *p.BrandID, p.ID).Order("id asc").Limit(limit).Find(&result)
	return result
}

//CountWithSameBrand returns a number of products that have the same brand as p
func (p *Product) CountWithSameBrand() int {
	count := 0
	if p.BrandID == nil {
		return count
	}
	DB.Model(&Product{}).Where("brand_id = ? and available = true", *p.BrandID).Count(&count)
	return count
}

//GetSpecialPrice returns product's meta special price or best variation's special price if they have been preloaded
func (p *Product) GetSpecialPrice() float64 {
	if len(p.Variations) == 0 {
		return p.SpecialPrice
	}
	price := p.Variations[0].SpecialPrice
	discount := p.Variations[0].DiscountPercent
	//find variation with maximum discount
	for _, v := range p.Variations {
		if v.DiscountPercent > discount {
			price = v.SpecialPrice
			discount = v.DiscountPercent
		}
	}
	return price
}

//GetRegularPrice returns product's meta regular price or best variation's regular price if they have been preloaded
func (p *Product) GetRegularPrice() float64 {
	if len(p.Variations) == 0 {
		return p.RegularPrice
	}
	price := p.Variations[0].RegularPrice
	discount := p.Variations[0].DiscountPercent
	//find variation with maximum discount
	for _, v := range p.Variations {
		if v.DiscountPercent > discount {
			price = v.RegularPrice
			discount = v.DiscountPercent
		}
	}
	return price
}

//GetDiscountPercent returns product's meta discount percent or best variation's discount percent if they have been preloaded
func (p *Product) GetDiscountPercent() float64 {
	if len(p.Variations) == 0 {
		return p.DiscountPercent
	}
	discount := p.Variations[0].DiscountPercent
	//find variation with maximum discount
	for _, v := range p.Variations {
		if v.DiscountPercent > discount {
			discount = v.DiscountPercent
		}
	}
	return discount
}
