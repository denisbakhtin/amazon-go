package models

import (
	"fmt"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/utility"
)

//BrowseNode stores info about product browse node
type BrowseNode struct {
	Model
	Title             string
	Products          []Product
	SimilarProducts   []Product `gorm:"-"`
	ParentID          *uint64   `gorm:"index:browse_node_parent_idx;"`
	Parent            *BrowseNode
	Children          []BrowseNode `gorm:"foreignkey:parent_id"`
	ThreeChildren     []BrowseNode `gorm:"-"` //shown at home page
	ProductCount      int
	OwnProductCount   int
	AllParentsLoaded  bool `gorm:"-"`
	AllChildrenLoaded bool `gorm:"-"`
	Description       string
	Image             string
}

//GetURL returns the proper browse node url
func (b *BrowseNode) GetURL() string {
	slug := utility.Parameterize(b.Title[:utility.Min(30, len(b.Title))])
	if len(slug) == 0 {
		slug = "empty"
	}
	return fmt.Sprintf("/tags/%d/%s", b.ID, slug)
}

//MainImage returns node's main image
func (b *BrowseNode) MainImage() string {
	product := Product{}
	DB.Where("browse_node_id = ? and available = true and image != ?", b.ID, "").Order("id asc").First(&product)
	if product.ID > 0 {
		return product.Image
	}
	if !b.AllChildrenLoaded {
		b.LoadAllChildren()
	}
	ids := b.AppendIDs(SELFANDCHILDREN)
	DB.Where("browse_node_id IN(?) and available = true and image != ?", ids, "").Order("id asc").First(&product)
	if product.ID > 0 {
		return product.Image
	}
	return "/images/no-image.jpg"
}

//LoadAllChildren loads full children hierarchy for the current node
func (b *BrowseNode) LoadAllChildren() {
	if !b.AllChildrenLoaded {
		DB.Where("parent_id = ?", b.ID).
			Preload("Children").
			Preload("Children.Children").
			Preload("Children.Children.Children").
			Preload("Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Find(&b.Children)
		b.AllChildrenLoaded = true
	}
}

//LoadAllParents loads full parent hierarchy for the current node
func (b *BrowseNode) LoadAllParents() {
	if !b.AllParentsLoaded {
		if b.ParentID != nil {
			parent := BrowseNode{}
			DB.
				Preload("Parent").
				Preload("Parent.Parent").
				Preload("Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				Preload("Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
				First(&parent, *b.ParentID)
			b.Parent = &parent
		}
		b.AllParentsLoaded = true
	}
}

//AppendIDs appends this and children or parent ids depending on parameter hierarchy
//Children or parents have to be loaded before calling this method
func (b *BrowseNode) AppendIDs(hierarchy string) []uint64 {
	ids := make([]uint64, 0, 20)
	switch hierarchy {
	case PARENTS:
		if b.Parent != nil {
			ids = append(ids, b.Parent.AppendIDs(SELFANDPARENTS)...)
		}
	case SELFANDPARENTS:
		ids = append(ids, b.ID)
		if b.Parent != nil {
			ids = append(ids, b.Parent.AppendIDs(SELFANDPARENTS)...)
		}
	case CHILDREN:
		for i := range b.Children {
			ids = append(ids, b.Children[i].AppendIDs(SELFANDCHILDREN)...)
		}
	case SELFANDCHILDREN:
		ids = append(ids, b.ID)
		for i := range b.Children {
			ids = append(ids, b.Children[i].AppendIDs(SELFANDCHILDREN)...)
		}
	default:
	}
	return ids
}

//Breadcrumbs returns tag's breadcrumbs
func (b *BrowseNode) Breadcrumbs() []Breadcrumb {
	b.LoadAllParents()
	return b.buildBreadcrumbs()
}

func (b *BrowseNode) buildBreadcrumbs() []Breadcrumb {
	crumbs := make([]Breadcrumb, 0, 20)
	if b.Parent != nil {
		crumbs = append(b.Parent.buildBreadcrumbs(), Breadcrumb{URL: b.Parent.GetURL(), Title: b.Parent.Title})
	}
	return crumbs
}

//TopParent returns the topmost browse node parent
func (b *BrowseNode) TopParent() *BrowseNode {
	if b.Parent != nil {
		return b.Parent.TopParent()
	}
	return b
}

//MenuNodes returns a list of top-level nodes suitable for main menu
func MenuNodes() []BrowseNode {
	var nodes []BrowseNode
	DB.Where("parent_id is null and title != ? and product_count > 0", "").
		Order("title asc").
		Find(&nodes)

	return nodes
}

//TopNodes returns a slice of all top level browse nodes
func TopNodes() []BrowseNode {
	var nodes []BrowseNode
	DB.Where("parent_id is null").Find(&nodes)
	return nodes
}

//IDStr returns string representation of id
func (b *BrowseNode) IDStr() string {
	return fmt.Sprintf("%d", b.ID)
}

//GetMetaKeywords returns meta keywords
func (b *BrowseNode) GetMetaKeywords() string {
	keywords := b.Title
	if topParent := b.TopParent(); topParent != b {
		keywords = keywords + ", " + topParent.Title
	}
	return fmt.Sprintf("%s, %s", keywords, config.SiteName)
}

//GetMetaDescription returns meta description
func (b *BrowseNode) GetMetaDescription() (result string) {
	title := b.Title
	if topParent := b.TopParent(); topParent != b {
		title = title + ", " + topParent.Title
	}
	return fmt.Sprintf("Buy %s at the %s gateway shopping centre. Discount and Free Shipping on eligible items.", title, config.SiteName)
}
