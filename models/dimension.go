package models

import "fmt"

//Dimension stores info about product dimension
type Dimension struct {
	Model
	Name     string
	Priority float64
	Pname    string
	Title    string
}

//GetTitle returns human readable title or system name if title is empty
func (d *Dimension) GetTitle() string {
	if len(d.Title) > 0 {
		return d.Title
	}
	return d.Name
}

//IDStr returns string id
func (d *Dimension) IDStr() string {
	return fmt.Sprintf("%d", d.ID)
}
