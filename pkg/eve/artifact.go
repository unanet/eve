package eve

type Artifact struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FeedType      string `json:"feed_type"`
	ProviderGroup string `json:"provider_group"`
	ImageTag      string `json:"image_tag"`
	ServicePort   int    `json:"service_port"`
	MetricsPort   int    `json:"metrics_port"`
}
