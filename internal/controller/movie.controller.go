package controller

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type movieController struct {
	movieService service.MovieService
}

// Create implements [MovieController].
func (m *movieController) Create(c *gin.Context) {
	var req request.MovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	movie := model.Movie{
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

// Delete implements [MovieController].
func (m *movieController) Delete(c *gin.Context) {
	id := c.Param("id")
	movie, err := m.movieService.FindByID(id)
	if err != nil {
		c.JSON(movieStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	err = m.movieService.Delete(id)
	if err != nil {
		c.JSON(movieStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success delete movie", movie))
}

// GetAll implements [MovieController].
func (m *movieController) GetAll(c *gin.Context) {
	movies, err := m.movieService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success get all movie", movies))
}

// GetByID implements [MovieController].
func (m *movieController) GetByID(c *gin.Context) {
	id := c.Param("id")
	movie, err := m.movieService.FindByID(id)
	if err != nil {
		c.JSON(movieStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success get movie by id", movie))
}

// Update implements [MovieController].
func (m *movieController) Update(c *gin.Context) {
	var req request.UpdateMovieRequest
	id := c.Param("id")

	movie, err := m.movieService.FindByID(id)

	if err != nil {
		c.JSON(movieStatusCode(err), response.ErrorResponse(err.Error()))
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
	err = m.movieService.Update(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success update movie", movie))

}

type MovieController interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Delete(c *gin.Context)
}

func NewMovieController(movieService service.MovieService) MovieController {
	return &movieController{movieService: movieService}
}

func movieStatusCode(err error) int {
	switch {
	case errors.Is(err, service.ErrInvalidMovieID):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrMovieNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
