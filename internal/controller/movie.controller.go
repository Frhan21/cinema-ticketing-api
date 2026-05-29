package controller

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"cinema-ticketing-api/pkg/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type movieController struct {
	movieService service.MovieService
}

func toMovieResponse(m *model.Movie) *response.MovieResponse {
	return &response.MovieResponse{
		ID:          m.ID.String(),
		Title:       m.Title,
		Genre:       m.Genre,
		Duration:    m.Duration,
		PosterUrl:   m.PosterUrl,
		Description: m.Description,
	}
}

// Create implements [MovieController].
func (m *movieController) Create(c *gin.Context) {
	var req request.MovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	movie := model.Movie{
		Title:       req.Title,
		Genre:       req.Genre,
		Duration:    req.Duration,
		PosterUrl:   req.PosterUrl,
		Description: req.Description,
	}

	err := m.movieService.Create(&movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success create movie", toMovieResponse(&movie)))
}

// Delete implements [MovieController].
func (m *movieController) Delete(c *gin.Context) {
	id := c.Param("id")
	_, err := m.movieService.FindByID(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	err = m.movieService.Delete(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success deleting movie", nil))
}

// GetAll implements [MovieController].
func (m *movieController) GetAll(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	movies, count, err := m.movieService.FindAll(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	var res []response.MovieResponse
	for _, movie := range movies {
		res = append(res, *toMovieResponse(&movie))
	}

	meta := &response.Meta{
		Page:      req.Page,
		PerPage:   req.PerPage,
		TotalData: int(count),
		TotalPage: pagination.GetTotalPage(count, req.PerPage),
	}

	c.JSON(http.StatusOK, response.SuccessResponseWithMeta("Success get all movie", res, meta))
}

// GetByID implements [MovieController].
func (m *movieController) GetByID(c *gin.Context) {
	id := c.Param("id")
	movie, err := m.movieService.FindByID(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Success get movie by id", toMovieResponse(movie)))
}

// Update implements [MovieController].
func (m *movieController) Update(c *gin.Context) {
	var req request.UpdateMovieRequest
	id := c.Param("id")

	movie, err := m.movieService.FindByID(id)

	if err != nil {
		c.Error(err)
		c.Abort()
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

	res := toMovieResponse(movie)

	c.JSON(http.StatusOK, response.SuccessResponse("Success update movie", res))

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
