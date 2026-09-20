package web_service

import (
	"fmt"
	"os"
	"path"
)

func (s *WebService) GetAuthPage() ([]byte, error) {
	htmlPath := path.Join(
		os.Getenv("PROJECT_ROOT"),
		"/public/auth.html",
	)

	html, err := s.webRepository.GetFile(htmlPath)
	if err != nil {
		return nil, fmt.Errorf("get file from repository: %w", err)
	}
	return html, nil
}
