package models

//Range is a view model that represents info about price or discount range
type Range struct {
	From  float64
	To    float64
	Title string
	Code  string
}
