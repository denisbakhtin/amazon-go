package models

import (
	"sync"
	"time"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" //do not delete, required for gorm
)

// Hierarchy relation constants
const (
	PARENTS         string = "Parents"
	SELFANDPARENTS  string = "SelfAndParents"
	CHILDREN        string = "Children"
	SELFANDCHILDREN string = "SelfAndChildren"
)

type cacheStruct struct {
	RWMutex sync.RWMutex

	CategoryDescriptions     []string
	CategoryMetaDescriptions []string
	CategoryMetaKeywords     []string

	TagDescriptions     []string
	TagMetaDescriptions []string
	TagDiscountScale    []string
	TagPriceScale       []string
}

var (
	//DB stores gorm db handler
	DB *gorm.DB
	//Cache stores global app cache
	Cache cacheStruct
)

//Model is a modification of Model struct with uint64 index
type Model struct {
	ID        uint64     `form:"id" gorm:"primary_key"`
	CreatedAt time.Time  `binding:"-" form:"-"`
	UpdatedAt time.Time  `binding:"-" form:"-"`
	DeletedAt *time.Time `binding:"-" form:"-"`
}

//InitDatabase initializes db handler
func InitDatabase() {
	var err error

	DB, err = gorm.Open("postgres", config.DBConnectionString)

	if err != nil {
		panic(err.Error())
	}

	if config.IsDebug() {
		//DB.LogMode(true)
	}
	if err := DB.AutoMigrate(&Account{}, &Brand{}, &BrowseNode{}, &Category{}, &Company{}, &Dimension{}, &Feed{},
		&Language{}, &Menu{}, &Operation{}, &Page{}, &ProcessedAsin{}, &ProcessedSpecification{},
		&ProductGroupType{}, &ProductGroup{}, &ProductTranslation{}, &ProductType{}, &Product{}, &QueuedAsin{}, &QueuedSpecification{},
		&QueuedTranslation{}, &Specification{}, &SyncLog{}, &TranslationRequest{}, &Variation{}, &Watch{}, &Binding{}, &Department{}, &ItemAttributes{}).Error; err != nil {
		panic(err)
	}
	seedDB()
}

//InitializeCache initializes app cache
//Obsolete
func InitializeCache() {
	category := Category{}
	DB.Select("id").Where("id != ?", config.NewArrivalsID).Order("id desc").First(&category)
	categoryMaxID := category.ID

	Cache.RWMutex.Lock()

	//Initialize slices with empty values
	Cache.CategoryDescriptions = make([]string, categoryMaxID+1)
	Cache.CategoryMetaDescriptions = make([]string, categoryMaxID+1)
	Cache.CategoryMetaKeywords = make([]string, categoryMaxID+1)

	Cache.RWMutex.Unlock()

	//Populate category slices with actual data
	var categories []Category
	DB.Where("id != ?", config.NewArrivalsID).Find(&categories)
	for i := range categories {
		desc := categories[i].compileDescription()
		metaDesc := categories[i].compileMetaDescription()
		metaKeyw := categories[i].compileMetaKeywords()
		//moving long computations out of lock (as well as granular per category locking) significantly reduces lock time
		Cache.RWMutex.Lock()
		Cache.CategoryDescriptions[categories[i].ID] = desc
		Cache.CategoryMetaDescriptions[categories[i].ID] = metaDesc
		Cache.CategoryMetaKeywords[categories[i].ID] = metaKeyw
		Cache.RWMutex.Unlock()
	}

}

//seedDB seeds DB with initial vital data
func seedDB() {
	//seed languages
	count := 0
	DB.Model(&Language{}).Count(&count)
	if count == 0 {
		enLang := Language{
			Model:       Model{ID: config.EnLangID},
			Code:        "en",
			Title:       "English",
			NativeTitle: "English",
		}
		if err := DB.Create(&enLang).Error; err != nil {
			panic(err)
		}
		ruLang := Language{
			Model:       Model{ID: config.RuLangID},
			Code:        "ru",
			Title:       "Russian",
			NativeTitle: "Русский",
		}
		if err := DB.Create(&ruLang).Error; err != nil {
			panic(err)
		}
	}

	//seed accounts
	count = 0
	DB.Model(&Account{}).Count(&count)
	if count == 0 {
		account := Account{
			Model:     Model{ID: config.AdminID},
			Email:     config.AdminEmail,
			FirstName: "Denis",
			LastName:  "Bakhtin",
			Role:      config.AdminRole,
			Password:  config.AdminPassword,
		}
		if err := DB.Create(&account).Error; err != nil {
			panic(err)
		}
	}

	//seed companies
	count = 0
	DB.Model(&Company{}).Count(&count)
	if count == 0 {
		comp := Company{
			Model:       Model{ID: config.MyCompanyID},
			AccountID:   config.AdminID,
			Description: "Default company",
			Show:        false,
		}
		if err := DB.Create(&comp).Error; err != nil {
			panic(err)
		}
	}

	//seed categories
	count = 0
	DB.Model(&Category{}).Count(&count)
	if count == 0 {
		newArrivals := Category{
			Model:       Model{ID: config.NewArrivalsID},
			Title:       "New arrivals",
			Description: "New, uncategorized products",
		}
		if err := DB.Create(&newArrivals).Error; err != nil {
			panic(err)
		}
	}
}
