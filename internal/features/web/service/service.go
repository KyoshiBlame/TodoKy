package web_service

type WebService struct {
	webRepository WebRepository
}

type WebRepository interface {
	GetFile(filePath string) ([]byte, error)
}

func NewWebService(
	webRepo WebRepository,
) *WebService {
	return &WebService{
		webRepository: webRepo,
	}
}
