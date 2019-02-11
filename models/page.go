package models

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/denisbakhtin/amazon-go/utility"
)

//Page stores info about static web page
type Page struct {
	Model
	Title           string `form:"title" binding:"required"`
	MetaKeywords    string `form:"meta_keywords"`
	MetaDescription string `form:"meta_description"`
	Body            string `form:"body"`
	Show            bool   `form:"show"`
	Submit          string `gorm:"-" form:"submit" binding:"required"`
}

//BeforeSave gorm hook
func (p *Page) BeforeSave() error {
	if len(p.Title) == 0 {
		return errors.New("Title is empty")
	}
	if len(p.Body) == 0 {
		return errors.New("Body is empty")
	}
	return nil
}

//GetURL returns the proper product url
func (p *Page) GetURL() string {
	return fmt.Sprintf("/pages/%d/%s", p.ID, utility.Parameterize(p.Title))
}

//GetBody returns html body
func (p *Page) GetBody() template.HTML {
	return template.HTML(p.Body)
}
