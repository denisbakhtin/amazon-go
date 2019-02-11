package controllers

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/utility"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	csrf "github.com/utrack/gin-csrf"
)

var tmpl *template.Template

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) (H gin.H) {
	defer func() {
		if e := recover(); e != nil {
			H["Authenticated"] = false
			H["Account"] = nil
			return
		}
	}()
	H = gin.H{
		"Title":              "", //page title
		"Context":            c,
		"Controller":         "",
		"SiteName":           config.SiteName,
		"SiteTitle":          config.SiteTitle,
		"GoogleVerification": config.GoogleVerification,
	}
	account, _ := c.Get("account")
	H["Csrf"] = csrf.GetToken(c)
	H["Authenticated"] = account != nil
	H["Account"] = account
	return
}

//loadTemplates loads html templates
func loadTemplates() *template.Template {
	tmpl = template.New("").Funcs(template.FuncMap{
		"SubmitAndBackTitle":         submitAndBackTitle,
		"SubmitAndViewTitle":         submitAndViewTitle,
		"SiteTitle":                  siteTitle,
		"SiteName":                   siteName,
		"SignUpEnabled":              signUpEnabled,
		"SidebarCategories":          sidebarCategories,
		"SidebarTags":                sidebarTags,
		"SearchDepartments":          searchDepartments,
		"Truncate":                   truncate,
		"GetPriceData":               getPriceData,
		"EqIDParentID":               eqIDParentID,
		"FormatDateTime":             formatDateTime,
		"Now":                        now,
		"SuperSaverURL":              superSaverURL,
		"PrivacyPolicyURL":           privacyPolicyURL,
		"AgreementURL":               agreementURL,
		"DisclaimerURL":              disclaimerURL,
		"WeAreSorryURL":              weAreSorryURL,
		"PaymentMethodsURL":          paymentMethodsURL,
		"ContactsURL":                contactsURL,
		"AboutURL":                   aboutURL,
		"TotalProducts":              totalProducts,
		"NoEscape":                   noEscape,
		"GenerateColor":              generateColor,
		"AccountRoles":               accountRoles,
		"PriceRanges":                PriceRanges,
		"DiscountRanges":             DiscountRanges,
		"NoItemsMessage":             noItemsMessage,
		"StringSliceContains":        utility.StringSliceContains,
		"NotEmpty":                   notEmpty,
		"NotEmptyDim":                notEmptyDim,
		"ItemDimensions":             itemDimensions,
		"PackageDimensions":          packageDimensions,
		"ItemWeight":                 itemWeight,
		"PackageWeight":              packageWeight,
		"ManufacturerRecommendedAge": manufacturerRecommendedAge,
		"ItemDate":                   itemDate,
		"UTCTime":                    utcTimeStr,
	})

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".gohtml") {
			var err error
			tmpl, err = tmpl.ParseFiles(path)
			if err != nil {
				return fmt.Errorf("%s: %s", path, err.Error())
			}
		}
		return nil
	}

	err := filepath.Walk(path.Join(config.AppPath, "views"), fn)
	if err != nil {
		log.Panic(err)
	}
	return tmpl
}

//submitAndBackTitle returns title for "submit and go back to list" button
func submitAndBackTitle() string {
	return "Submit and Back to list"
}

//submitAndViewTitle returns title for "submit and view" button
func submitAndViewTitle() string {
	return "Submit and View"
}

//siteTitle return site meta title
func siteTitle() string {
	return config.SiteTitle
}

//siteName return site name
func siteName() string {
	return config.SiteName
}

//signUpEnabled returns SignUpEnabled flag
func signUpEnabled() bool {
	return config.SignUpEnabled
}

//sidebarCategories returns sidebar categories
//obsolete
func sidebarCategories() []models.Category {
	var categories []models.Category
	models.DB.Select("id, title").Where("parent_id is NULL and id != ?", config.NewArrivalsID).Find(&categories)
	return categories
}

//sidebarTags returns sidebar tags
func sidebarTags(category models.Category) []models.BrowseNode {
	var tags []models.BrowseNode
	models.DB.Select("id, title").Where("parent_id is NULL and product_count > 0").Order("title").Find(&tags)
	return tags
}

//searchDepartments returns a slice of search scopes
func searchDepartments() []models.SearchDepartment {
	result := make([]models.SearchDepartment, 0, 20) //length = 0, capacity = 20, should be fine
	//basic search list
	result = append(result, models.SearchDepartment{ID: 0, Title: "All Categories", Class: "global"})
	nodes := models.MenuNodes()
	for i := range nodes {
		result = append(result, models.SearchDepartment{ID: nodes[i].ID, Title: nodes[i].Title})
	}
	return result
}

//truncate string to num chars
func truncate(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func getPriceData(str string, data map[string]string) string {
	return data[str]
}

func eqIDParentID(id uint64, parentID *uint64) bool {
	if parentID == nil {
		return false
	}
	return id == *parentID
}

func now() time.Time {
	return time.Now()
}

//formatDateTime prints timestamp in human format
func formatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func superSaverURL() string {
	superSaverPage := models.Page{}
	models.DB.First(&superSaverPage, config.SuperSaverPageID)
	return superSaverPage.GetURL()
}

func privacyPolicyURL() string {
	privacyPolicyPage := models.Page{}
	models.DB.First(&privacyPolicyPage, config.PrivacyPolicyID)
	return privacyPolicyPage.GetURL()
}

func agreementURL() string {
	agreementPage := models.Page{}
	models.DB.First(&agreementPage, config.SiteAgreementID)
	return agreementPage.GetURL()
}

func disclaimerURL() string {
	disclaimerPage := models.Page{}
	models.DB.First(&disclaimerPage, config.AffiliateDisclaimerID)
	return disclaimerPage.GetURL()
}

func weAreSorryURL() string {
	sorryPage := models.Page{}
	models.DB.First(&sorryPage, config.WeAreSorryID)
	return sorryPage.GetURL()
}

func paymentMethodsURL() string {
	paymentPage := models.Page{}
	models.DB.First(&paymentPage, config.PaymentMethodsID)
	return paymentPage.GetURL()
}

func contactsURL() string {
	contactsPage := models.Page{}
	models.DB.First(&contactsPage, config.ContactsID)
	return contactsPage.GetURL()
}

func aboutURL() string {
	aboutPage := models.Page{}
	models.DB.First(&aboutPage, config.AboutID)
	return aboutPage.GetURL()
}

func totalProducts() uint64 {
	var count uint64
	models.DB.Model(&models.Product{}).Where("available = true").Count(&count)
	return count
}

func noEscape(str string) template.HTML {
	return template.HTML(str)
}

func generateColor(step, colorCount int, alpha float64) string {
	hueStep := int(math.Floor(330 / float64(colorCount)))
	hue := step * hueStep
	if hue > 100 {
		hue += 30
	}
	saturation := 65
	if (step & 1) == 1 {
		saturation = 90
	}
	luminosity := 55
	if (step & 1) == 1 {
		luminosity = 80
	}
	if alpha == 1 {
		hslToHex(float64(hue)/360, float64(saturation)/100, float64(luminosity)/100)
	}
	return hslToRgba(float64(hue)/360, float64(saturation)/100, float64(luminosity)/100, alpha)
}

//hslToHex converts hsl color to the html hex string, where h, s, l are [0,1] floats
func hslToHex(h, s, l float64) string {
	var r, g, b float64

	if s == 0 {
		// achromatic
		r = l
		g = l
		b = l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hue2rgb(p, q, h+1/float64(3))
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1/float64(3))
	}

	return fmt.Sprintf("#%x%x%x", int(math.Min(r*255, float64(255))), int(math.Min(g*255, float64(255))), int(math.Min(b*255, float64(255))))
}

//hslToRgba converts hsl color to rgba string, where h, s, l are [0,1] floats
func hslToRgba(h, s, l, alpha float64) string {
	var r, g, b float64

	if s == 0 {
		// achromatic
		r = l
		g = l
		b = l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hue2rgb(p, q, h+1/float64(3))
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1/float64(3))
	}

	return fmt.Sprintf("rgba(%d,%d,%d,%0.2f)", int(math.Min(r*255, float64(255))), int(math.Min(g*255, float64(255))), int(math.Min(b*255, float64(255))), alpha)
}

//helper function
func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if (t * 6) < 1 {
		return p + (q-p)*6*t
	}
	if (t * 2) < 1 {
		return q
	}
	if (t * 3) < 2 {
		return p + (q-p)*(2/float64(3)-t)*6
	}
	return p
}

func accountRoles() []models.AccountRole {
	return []models.AccountRole{
		models.AccountRole{Code: config.UserRole, Title: "User"},
		models.AccountRole{Code: config.AdminRole, Title: "Admin"},
	}
}

//PriceRanges returns a slice of price ranges
func PriceRanges() []models.Range {
	return []models.Range{
		models.Range{From: 0.0, To: 100000000.0, Title: "Any price", Code: ""},
		models.Range{From: 0.0, To: 25.0, Title: "Under $25", Code: "under_25"},
		models.Range{From: 25.0, To: 50.0, Title: "$25 to $50", Code: "25_to_50"},
		models.Range{From: 50.0, To: 100.0, Title: "$50 to $100", Code: "50_to_100"},
		models.Range{From: 100.0, To: 100000000.0, Title: "Over $100", Code: "over_100"},
	}
}

//DiscountRanges returns a slice of discount ranges
func DiscountRanges() []models.Range {
	return []models.Range{
		models.Range{From: 0.0, To: 100000000.0, Title: "Any value", Code: ""},
		models.Range{From: 0.0, To: 25.0, Title: "Under 25%", Code: "under_25"},
		models.Range{From: 25.0, To: 50.0, Title: "25% to 50%", Code: "25_to_50"},
		models.Range{From: 50.0, To: 75.0, Title: "50% to 75%", Code: "50_to_75"},
		models.Range{From: 75.0, To: 100000000.0, Title: "Over 75%", Code: "over_75"},
	}
}

//applyProductFilters applies product filters to gorm query
func applyProductFilters(c *gin.Context, query *gorm.DB) *gorm.DB {
	if c.Query("on_sale") == "1" {
		query = query.Where("discount_percent > 0.01")
	}
	if c.Query("free_shipping") == "1" {
		query = query.Where("free_shipping = true")
	}
	if len(c.Query("price")) > 0 {
		pranges := PriceRanges()
		for _, r := range pranges {
			if r.Code == c.Query("price") {
				query = query.
					Where("EXISTS ( SELECT null FROM variations WHERE variations.special_price >= ? AND variations.special_price < ? AND variations.product_id = products.id )", r.From, r.To).
					Preload("Variations", func(db *gorm.DB) *gorm.DB {
						return db.Where("special_price >= ? AND special_price < ?", r.From, r.To)
					})
				break
			}
		}
	}
	if len(c.Query("discount")) > 0 {
		dranges := DiscountRanges()
		for _, r := range dranges {
			if r.Code == c.Query("discount") {
				query = query.
					Where("EXISTS ( SELECT null FROM variations WHERE variations.discount_percent >= ? AND variations.discount_percent < ? AND variations.product_id = products.id )", r.From, r.To).
					Preload("Variations", func(db *gorm.DB) *gorm.DB {
						return db.Where("discount_percent >= ? AND discount_percent < ?", r.From, r.To)
					})
				break
			}
		}
	}
	if len(c.QueryArray("brand[]")) > 0 {
		query = query.Where("brand_id IN(?)", c.QueryArray("brand[]"))
	}
	if len(c.QueryArray("tag[]")) > 0 {
		query = query.Where("browse_node_id IN(?)", c.QueryArray("tag[]"))
	}
	return query
}

func noItemsMessage() string {
	return "Sorry, your search returned no results."
}

func notEmpty(s string) bool {
	return len(s) > 0
}

func notEmptyDim(s string) bool {
	return strings.Index(s, "0.0") != 0
}

func dimSlice(dim string) []string {
	d := strings.Split(dim, " ")
	if len(d) > 1 {
		return []string{d[0], normalizeDimensionUnit(strings.Join(d[1:len(d)], " "))}
	}
	return []string{d[0], ""}
}

func dimsSlice(length, width, height string) [][]string {
	return [][]string{dimSlice(length), dimSlice(width), dimSlice(height)}
}

//converts hundredths to base units atm
func convertDim(dim []string) []string {
	if strings.Index(dim[1], "hundredths ") > -1 {
		value, _ := strconv.Atoi(dim[0])
		dim[0] = fmt.Sprintf("%.1f", float64(value)/100)
		dim[1] = strings.TrimPrefix(dim[1], "hundredths ")
	}
	return dim
}

//converts hundredths to base units atm
func convertDims(dims [][]string) [][]string {
	for i := range dims {
		dims[i] = convertDim(dims[i])
	}
	return dims
}

func normalizeDimensionUnit(unit string) string {
	return strings.ToLower(strings.Replace(unit, "-", " ", -1))
}

func itemDimensions(attrs models.ItemAttributes) string {
	dims := dimsSlice(attrs.Length, attrs.Width, attrs.Height)
	dims = convertDims(dims)

	return fmt.Sprintf("%s x %s x %s %s", dims[0][0], dims[1][0], dims[2][0], dims[0][1])
}

func itemWeight(attrs models.ItemAttributes) string {
	dim := dimSlice(attrs.Weight)
	dim = convertDim(dim)
	return fmt.Sprintf("%s %s", dim[0], dim[1])
}

func packageWeight(attrs models.ItemAttributes) string {
	dim := dimSlice(attrs.PackageWeight)
	dim = convertDim(dim)
	return fmt.Sprintf("%s %s", dim[0], dim[1])
}

func packageDimensions(attrs models.ItemAttributes) string {
	dims := dimsSlice(attrs.PackageLength, attrs.PackageWidth, attrs.PackageHeight)
	dims = convertDims(dims)

	return fmt.Sprintf("%s x %s x %s %s", dims[0][0], dims[1][0], dims[2][0], dims[0][1])
}

func manufacturerRecommendedAge(attrs models.ItemAttributes) string {
	minage, _ := strconv.Atoi(attrs.ManufacturerMinimumAge)
	maxage, _ := strconv.Atoi(attrs.ManufacturerMaximumAge)
	return fmt.Sprintf("%.0f - %.0f years", float64(minage)/12, float64(maxage)/12)
}

func itemDate(date string) string {
	t, _ := time.Parse("2006-01-02", date)
	return t.Format("January 02, 2006")
}

func utcTimeStr(t time.Time) string {
	return utility.UTCTime(t).Format(time.RFC1123)
}
