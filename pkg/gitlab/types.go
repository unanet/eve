package gitlab

import (
	"time"
)

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

type Release struct {
	TagName         string    `json:"tag_name"`
	Description     string    `json:"description"`
	Name            string    `json:"name"`
	DescriptionHTML string    `json:"description_html"`
	CreatedAt       time.Time `json:"created_at"`
	ReleasedAt      time.Time `json:"released_at"`
	Author          struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"author"`
	Commit struct {
		ID             string        `json:"id"`
		ShortID        string        `json:"short_id"`
		Title          string        `json:"title"`
		CreatedAt      time.Time     `json:"created_at"`
		ParentIds      []interface{} `json:"parent_ids"`
		Message        string        `json:"message"`
		AuthorName     string        `json:"author_name"`
		AuthorEmail    string        `json:"author_email"`
		AuthoredDate   time.Time     `json:"authored_date"`
		CommitterName  string        `json:"committer_name"`
		CommitterEmail string        `json:"committer_email"`
		CommittedDate  time.Time     `json:"committed_date"`
	} `json:"commit"`
	Milestones []struct {
		ID          int       `json:"id"`
		Iid         int       `json:"iid"`
		ProjectID   int       `json:"project_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		State       string    `json:"state"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DueDate     time.Time `json:"due_date"`
		StartDate   time.Time `json:"start_date"`
		WebURL      string    `json:"web_url"`
		IssueStats  struct {
			Total  int `json:"total"`
			Closed int `json:"closed"`
		} `json:"issue_stats"`
	} `json:"milestones"`
	CommitPath string `json:"commit_path"`
	TagPath    string `json:"tag_path"`
	Assets     struct {
		Count   int `json:"count"`
		Sources []struct {
			Format string `json:"format"`
			URL    string `json:"url"`
		} `json:"sources"`
		Links []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			URL      string `json:"url"`
			External bool   `json:"external"`
			LinkType string `json:"link_type"`
		} `json:"links"`
	} `json:"assets"`
	Evidences []struct {
		Sha         string    `json:"sha"`
		Filepath    string    `json:"filepath"`
		CollectedAt time.Time `json:"collected_at"`
	} `json:"evidences"`
}
