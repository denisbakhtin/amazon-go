package viewmodels

//Page used for binding new/edit Page form values
type Page struct {
	Title           string `form:"title" binding:"required"`
	MetaKeywords    string `form:"meta_keywords"`
	MetaDescription string `form:"meta_description"`
	Body            string `form:"body"`
	Show            bool   `form:"show"`
	Submit          string `form:"submit" binding:"required"`
}
