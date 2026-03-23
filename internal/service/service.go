package service

type repo interface{}

type Service struct {
	repo repo
}

func New(repo repo) *Service {
	return &Service{repo: repo}
}
