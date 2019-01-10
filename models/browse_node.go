package models

import (
	"fmt"

	"github.com/denisbakhtin/amazon-go/utility"
	"github.com/denisbakhtin/amazon-go/viewmodels"
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
	ProductCount      int
	OwnProductCount   int
	AllParentsLoaded  bool `gorm:"-"`
	AllChildrenLoaded bool `gorm:"-"`
	Description       string
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

//LoadAllParents loads full parent hierarchy for the current node
func (b *BrowseNode) LoadAllParents() {
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

//GetThreeChildren returns a slice of 3 children and/or successors
func (b *BrowseNode) GetThreeChildren() []BrowseNode {
	if !b.AllChildrenLoaded {
		b.LoadAllChildren()
	}
	result := make([]BrowseNode, 0, 3)
	result = appendChildren(3, b, result)
	return result
}

func appendChildren(number int, b *BrowseNode, result []BrowseNode) []BrowseNode {
	//append direct children first
	for i := 0; len(result) < number && i < len(b.Children); i++ {
		if b.Children[i].OwnProductCount > 0 {
			result = append(result, b.Children[i])
		}
	}
	if len(result) < number {
		for i := range b.Children {
			result = appendChildren(number, &b.Children[i], result)
		}
	}
	return result
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

//Breadcrumbs returns browse node parent breadcrumbs
func (b *BrowseNode) Breadcrumbs() []viewmodels.Breadcrumb {
	if !b.AllParentsLoaded {
		b.LoadAllParents()
	}
	return b.buildBreadcrumbs()
}

func (b *BrowseNode) buildBreadcrumbs() []viewmodels.Breadcrumb {
	crumbs := make([]viewmodels.Breadcrumb, 0, 20)
	if b.Parent != nil {
		crumbs = append(b.Parent.buildBreadcrumbs(), viewmodels.Breadcrumb{URL: b.Parent.GetURL(), Title: b.Parent.Title})
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
