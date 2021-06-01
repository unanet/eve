package eve

type Feed struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PromotionOrder int    `json:"promotion_order"`
	FeedType       string `json:"feed_type"`
	Alias          string `json:"alias"`
}
