package info

import (
	"strings"

	"github.com/lhlyu/gitx/internal/git"
)

type Service struct {
	git *git.Client
}

func NewService() *Service {
	return &Service{
		git: git.NewClient(),
	}
}

func (s *Service) Collect() (*Info, error) {
	info := &Info{}

	if out, err := s.git.Run("config", "user.name"); err == nil {
		info.UserName = strings.TrimSpace(string(out))
	}

	if out, err := s.git.Run("config", "user.email"); err == nil {
		info.UserEmail = strings.TrimSpace(string(out))
	}

	if out, err := s.git.Run("branch", "--show-current"); err == nil {
		info.Branch = strings.TrimSpace(string(out))
	}

	if out, err := s.git.Run("remote", "get-url", "origin"); err == nil {
		info.RemoteURL = strings.TrimSpace(string(out))
	}

	if out, err := s.git.Run("status", "--porcelain"); err == nil {
		output := strings.TrimSpace(string(out))
		if output == "" {
			info.IsClean = true
		} else {
			lines := strings.Split(output, "\n")
			info.IsClean = false
			info.ChangedFiles = len(lines)
		}
	}

	return info, nil
}
