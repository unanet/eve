package gitlab

type TagOptions struct {
	ProjectID int    `url:"-"`
	TagName   string `url:"tag_name,omitempty"`
	GitHash   string `url:"ref,omitempty"`
}

type Tag struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Commit  struct {
		ID             string   `json:"id"`
		ShortID        string   `json:"short_id"`
		CreatedAt      string   `json:"created_at"`
		ParentIds      []string `json:"parent_ids"`
		Title          string   `json:"title"`
		Message        string   `json:"message"`
		AuthorName     string   `json:"author_name"`
		AuthorEmail    string   `json:"author_email"`
		AuthoredDate   string   `json:"authored_date"`
		CommitterName  string   `json:"committer_name"`
		CommitterEmail string   `json:"committer_email"`
		CommittedDate  string   `json:"committed_date"`
		WebURL         string   `json:"web_url"`
	} `json:"commit"`
	Release   interface{} `json:"release"`
	Protected bool        `json:"protected"`
}
