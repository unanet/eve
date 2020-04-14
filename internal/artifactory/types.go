package artifactory

type VersionResponse struct {
	Version string `json:"version"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type MoveRequest struct {
	RepoKey       string `json:"repoKey"`
	Path          string `json:"path"`
	TargetRepoKey string `json:"targetRepoKey"`
	TargetPath    string `json:"targetPath"`
}
