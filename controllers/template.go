package controllers

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/viewmodels"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

var tmpl *template.Template

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	account, _ := c.Get("account")

	return gin.H{
		"Title":         "", //page title
		"Context":       c,
		"Csrf":          csrf.GetToken(c),
		"Authenticated": account != nil,
		"Account":       account,
		"Controller":    "",
	}
}

//loadTemplates loads html templates
func loadTemplates() *template.Template {
	tmpl = template.New("").Funcs(template.FuncMap{
		"SubmitAndBackTitle": submitAndBackTitle,
		"SubmitAndViewTitle": submitAndViewTitle,
		"SiteTitle":          siteTitle,
		"SiteName":           siteName,
		"SignUpEnabled":      signUpEnabled,
		"SidebarCategories":  sidebarCategories,
		"SidebarTags":        sidebarCategoryTags,
		"SearchDepartments":  searchDepartments,
		"Truncate":           truncate,
		"GetPriceData":       getPriceData,
		"EqIDParentID":       eqIDParentID,
		"FormatDateTime":     formatDateTime,
		"Now":                now,
		"SuperSaverURL":      superSaverURL,
		"PrivacyPolicyURL":   privacyPolicyURL,
		"AgreementURL":       agreementURL,
		"DisclaimerURL":      disclaimerURL,
		"WeAreSorryURL":      weAreSorryURL,
		"PaymentMethodsURL":  paymentMethodsURL,
		"ContactsURL":        contactsURL,
		"AboutURL":           aboutURL,
		"HomeTopProducts":    homeTopProducts,
		"TotalProducts":      totalProducts,
		"NoEscape":           noEscape,
		"GenerateColor":      generateColor,
		"AccountRoles":       accountRoles,
		"PriceRanges":        PriceRanges,
		"DiscountRanges":     DiscountRanges,
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
func sidebarCategories() []models.Scale {
	var sCategories []models.Category
	models.DB.Where("parent_id is NULL and id != ?", config.NewArrivalsID).Find(&sCategories)
	sidebarCategories := make([]models.Scale, len(sCategories))
	for i := range sCategories {
		sidebarCategories[i] = models.Scale{Title: sCategories[i].Title, URL: sCategories[i].GetURL()}
	}
	return sidebarCategories
}

//sidebarCategoryTags returns sidebar category tags
func sidebarCategoryTags(category models.Category) []models.Scale {
	var sCategories []models.Category
	models.DB.Where("parent_id is NULL and id != ?", config.NewArrivalsID).Find(&sCategories)
	sidebarCategories := make([]models.Scale, len(sCategories))
	for i := range sCategories {
		sidebarCategories[i] = models.Scale{Title: sCategories[i].Title, URL: sCategories[i].GetURL()}
	}
	return sidebarCategories
}

//searchDepartments returns a slice of search scopes
func searchDepartments() []viewmodels.SearchDepartment {
	result := make([]viewmodels.SearchDepartment, 0, 20) //length = 0, capacity = 20, should be fine
	//basic search list
	result = append(result, viewmodels.SearchDepartment{ID: 0, Title: "All Categories", Class: "global"})
	nodes := models.MenuNodes()
	for i := range nodes {
		result = append(result, viewmodels.SearchDepartment{ID: nodes[i].ID, Title: nodes[i].Title})
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

func homeTopProducts() []models.Product {
	var ids []uint64
	//models.DB.Model(&models.Product{}).Where("available = true").Order("discount_percent desc, id desc").Limit(30).Pluck("image", &images)
	models.DB.Table("products").Where("discount_percent > 0 and available = true and image != ?", "").Select("min(id) as idd").Group("browse_node_id").Order("browse_node_id desc").Limit(30).Pluck("idd", &ids)
	var products []models.Product
	models.DB.Where("id IN(?)", ids).Order("id asc").Find(&products)
	return products
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

func accountRoles() []viewmodels.AccountRole {
	return []viewmodels.AccountRole{
		viewmodels.AccountRole{Code: config.UserRole, Title: "User"},
		viewmodels.AccountRole{Code: config.AdminRole, Title: "Admin"},
	}
}

//PriceRanges returns a slice of price ranges
func PriceRanges() []viewmodels.Range {
	return []viewmodels.Range{
		viewmodels.Range{From: 0.0, To: 100000000.0, Title: "Any price", Code: "any"},
		viewmodels.Range{From: 0.0, To: 25.0, Title: "Under $25", Code: "under_25"},
		viewmodels.Range{From: 25.01, To: 50.0, Title: "$25 to $50", Code: "25_to_50"},
		viewmodels.Range{From: 50.01, To: 100.0, Title: "$50 to $100", Code: "50_to_100"},
		viewmodels.Range{From: 100.01, To: 100000000.0, Title: "Over $100", Code: "over_100"},
	}
}

//DiscountRanges returns a slice of discount ranges
func DiscountRanges() []viewmodels.Range {
	return []viewmodels.Range{
		viewmodels.Range{From: 0.0, To: 100000000.0, Title: "Any value", Code: "any"},
		viewmodels.Range{From: 0.0, To: 25.0, Title: "Under 25%", Code: "under_25"},
		viewmodels.Range{From: 25.01, To: 50.0, Title: "25% to 50%", Code: "25_to_50"},
		viewmodels.Range{From: 50.01, To: 75.0, Title: "50% to 75%", Code: "50_to_75"},
		viewmodels.Range{From: 75.01, To: 100000000.0, Title: "Over 75%", Code: "over_75"},
	}
}
