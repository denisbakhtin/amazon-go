package models

//Menu stores info about menu
type Menu struct {
	Model
	Title        string
	URL          string
	Priority     int
	Show         bool
	ProductCount int
}

//BuildMenu builds main menu
/* func BuildMenu() {
	DB.Model(&Menu{}).Update("show", false)
	nodes := MenuNodes()
	for i := range nodes {
		menu := Menu{}
		if err := DB.FirstOrCreate(&menu, Menu{Title: nodes[i].Title}).Error; err != nil {
			log.Println(err.Error())
			return
		}
		menu.URL = nodes[i].GetURL()
		menu.Priority = i
		menu.Show = true
		menu.ProductCount = nodes[i].ProductCount
		if err := DB.Save(&menu).Error; err != nil {
			log.Println(err.Error())
			return
		}
	}
} */
