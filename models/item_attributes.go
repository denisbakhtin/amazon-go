package models

//ItemAttributes stores most of item attributes
type ItemAttributes struct {
	Model
	VariationID            uint64 `gorm:"index:itemattributes_variation_idx"`
	Feature                string
	Studio                 string
	Publisher              string
	Manufacturer           string
	Author                 string
	Label                  string
	PartNumber             string
	ItemModel              string //xml: Model
	MPN                    string
	UPC                    string
	EAN                    string
	SKU                    string
	EISBN                  string
	ISBN                   string
	Warranty               string
	OperationSystem        string
	NumberOfDiscs          string
	NumberOfIssues         string
	IssuesPerYear          string
	NumberOfItems          string
	NumberOfPages          string
	NumberOfTracks         string
	MediaType              string
	LegalDisclaimer        string
	Creator                string
	Actor                  string
	Artist                 string
	AspectRatio            string
	AudienceRating         string
	AudioFormat            string
	ClothingSize           string
	Color                  string
	Director               string
	Edition                string
	EpisodeSequence        string
	CEROAgeRating          string
	ESRBAgeRating          string
	Format                 string
	Genre                  string
	HardwarePlatform       string
	HazardousMaterialType  string
	IsAdultProduct         string
	Height                 string
	Length                 string
	Width                  string
	Weight                 string
	PackageHeight          string
	PackageLength          string
	PackageWidth           string
	PackageWeight          string
	PackageQuantity        string
	Language               string
	PublicationDate        string
	ReleaseDate            string
	MetalType              string
	MaterialType           string
	RunningTime            string
	Size                   string
	ManufacturerMinimumAge string
	ManufacturerMaximumAge string
	PictureFormat          string
}
