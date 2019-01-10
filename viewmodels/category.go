package viewmodels

//Category used for binding new/edit category form values
type Category struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
	ParentID    uint64 `form:"parent_id"`
	Submit      string `form:"submit" binding:"required"`
}
