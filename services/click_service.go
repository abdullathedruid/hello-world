package services

var clickCount = 0

type ClickService struct{}

func NewClickService() *ClickService {
	return &ClickService{}
}

func (s *ClickService) IncrementClick() int {
	clickCount++
	return clickCount
}

func (s *ClickService) GetCount() int {
	return clickCount
}

func (s *ClickService) Reset() {
	clickCount = 0
}
