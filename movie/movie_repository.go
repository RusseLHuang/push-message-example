package movie

type MovieRepository struct{}

type MovieModel struct {
	Name string
}

func NewMovieRepository() *MovieRepository {
	return &MovieRepository{}
}

func (d *MovieRepository) FindById(id int) {

}
