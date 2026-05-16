package controller

import (
	"cinema-ticketing-api/models"
	"cinema-ticketing-api/request"
	"cinema-ticketing-api/response"
	"cinema-ticketing-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type movieController struct {
	movieService service.MovieService
}

// CreateMovie implements [MovieController].
func (m *movieController) CreateMovie(c *gin.Context) {
	var req request.MovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	movie := models.Movie{
		Title:     req.Title,
		Genre:     req.Genre,
		Duration:  req.Duration,
		PosterUrl: req.PosterUrl,
	}

	err := m.movieService.Create(&movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success create movie", movie))
}

// DeleteMovie implements [MovieController].
func (m *movieController) DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	movie, err := m.movieService.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Movie not found"))
		return
	}

	err = m.movieService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success delete movie", movie))
}

// GetAllMovie implements [MovieController].
func (m *movieController) GetAllMovie(c *gin.Context) {
	movies, err := m.movieService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success get all movie", movies))
}

// GetByIdMovie implements [MovieController].
func (m *movieController) GetByIdMovie(c *gin.Context) {
	id := c.Param("id")
	movie, err := m.movieService.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Movie not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success get movie by id", movie))
}

// UpdateMovie implements [MovieController].
func (m *movieController) UpdateMovie(c *gin.Context) {
	var req request.UpdateMovieRequest
	id := c.Param("id")

	movie, err := m.movieService.GetById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Movie not found"))
		return
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	if req.Title != nil {
		movie.Title = *req.Title
	}
	if req.Genre != nil {
		movie.Genre = *req.Genre
	}
	if req.Duration != nil {
		movie.Duration = *req.Duration
	}
	if req.PosterUrl != nil {
		movie.PosterUrl = *req.PosterUrl
	}
	updatedMovie, err := m.movieService.Update(id, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success update movie", updatedMovie))

}

type MovieController interface {
	CreateMovie(c *gin.Context)
	UpdateMovie(c *gin.Context)
	GetAllMovie(c *gin.Context)
	GetByIdMovie(c *gin.Context)
	DeleteMovie(c *gin.Context)
}

func NewMovieController(movieService service.MovieService) MovieController {
	return &movieController{movieService: movieService}
}
