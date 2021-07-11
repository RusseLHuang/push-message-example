package movie

type MovieService struct{}

func NewMovieService() *MovieService {
	return &MovieService{}
}

func (m *MovieService) GetUserRecommendation(userId int) string {
	return "Iron Man"
}
