package models

import (
	"fmt"
	"strings"
)

//Variation stores info about product variation
type Variation struct {
	Model
	Asin             string
	Title            string
	Available        bool
	Feature          string
	RegularPrice     float64
	SpecialPrice     float64
	Currency         string
	Discount         float64
	DiscountPercent  float64
	OfferListingID   string
	ProductID        uint64 `gorm:"index:variation_product_idx"`
	AvailabilityNote string
	FreeShipping     bool
	Images           string
	Dim1Id           *uint64 //this and below is ugly, but...
	Dim1Value        string
	Dim2Id           *uint64
	Dim2Value        string
	Dim3Id           *uint64
	Dim3Value        string
	Dim4Id           *uint64
	Dim4Value        string
	Dim5Id           *uint64
	Dim5Value        string
	CategoryID       uint64 `gorm:"index:variation_category_idx"`
	Category         Category
}

//MainImage returns variation's main image
func (v Variation) MainImage() string {
	//Empty string, i guess is {}
	if len(v.Images) > 0 {
		images := strings.Split(v.Images, ",")
		return images[0]
	}
	return "/images/no-image.jpg"
}

//MainImageTitle returns variation's main image title
func (v Variation) MainImageTitle() string {
	//Empty string, i guess is {}
	if len(v.Images) > 0 {
		return fmt.Sprintf("%s photo", v.Title)
	}
	return "No image currently available"
}

//ImageSlice returns variation imates
func (v Variation) ImageSlice() []string {
	//Empty string, i guess is {}
	if len(v.Images) > 0 {
		images := strings.Split(v.Images, ",")
		return images
	}
	return nil
}

//FeatureSlice returns the feature slice for product view
func (v Variation) FeatureSlice() (features []string) {
	if len(v.Feature) > 0 {
		features = strings.Split(v.Feature, "<br/>")
	}
	return
}

//Attributes returns a string with all variation dimensions combined in readable format
func (v Variation) Attributes() string {
	attrs := make([]string, 0, 5)
	dim1, dim2, dim3, dim4, dim5 := getDim(v.Dim1Id), getDim(v.Dim2Id), getDim(v.Dim3Id), getDim(v.Dim4Id), getDim(v.Dim5Id)
	if dim1.ID > 0 {
		attrs = append(attrs, fmt.Sprintf("%s: %s", dim1.Title, v.Dim1Value))
	}
	if dim2.ID > 0 {
		attrs = append(attrs, fmt.Sprintf("%s: %s", dim2.Title, v.Dim2Value))
	}
	if dim3.ID > 0 {
		attrs = append(attrs, fmt.Sprintf("%s: %s", dim3.Title, v.Dim3Value))
	}
	if dim4.ID > 0 {
		attrs = append(attrs, fmt.Sprintf("%s: %s", dim4.Title, v.Dim4Value))
	}
	if dim5.ID > 0 {
		attrs = append(attrs, fmt.Sprintf("%s: %s", dim5.Title, v.Dim5Value))
	}
	return strings.Join(attrs, ", ")
}

func getDim(id *uint64) Dimension {
	dim := Dimension{}
	if id != nil {
		DB.First(&dim, *id)
	}
	return dim
}
