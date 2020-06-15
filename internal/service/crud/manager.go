package crud

type Repo interface {
}

func NewManager(r Repo) *Manager {
	return &Manager{
		repo: r,
	}
}

type Manager struct {
	repo Repo
}
