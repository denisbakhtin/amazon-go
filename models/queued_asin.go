package models

//QueuedAsin stores info about queued asins
type QueuedAsin struct {
	Model
	Asin     string  `form:"asin" binding:"required"`
	FeedID   *uint64 `form:"feed_id" gorm:"index:queued_asin_feed_idx"`
	Priority int64   `form:"priority"` //dont know why it is not float64, but let it be
	Feed     Feed    `form:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
}
