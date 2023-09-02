package v1

import (
	"avito-rest-api/internal/entity"
	"avito-rest-api/internal/service"
	customError "avito-rest-api/package/error"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type segmentRoutes struct {
	segmentService service.Segment
}

func newSegmentRoutes(g *echo.Group, segmentService service.Segment) {
	r := &segmentRoutes{
		segmentService: segmentService,
	}

	g.POST("", r.create)
	g.GET("", r.getAll)
	g.GET("/:name", r.getByName)
	g.DELETE("/:name", r.deleteByName)
}

type CreateResponse struct {
	Name string `json:"name" example:"AVITO_MUSIC_SERVICE" validate:"required"`
}

// @Summary Создать сегмент
// @Description Создаёт сегмент на основе информации в теле запроса
// @Tags segments
// @Accept json
// @Produce json
// @Param data body service.SegmentCreateInput true "Структура с информацией о создаваемом сегменте"
// @Success 201 {object} CreateResponse "Наименование созданного сегмента"
// @Failure 400 {object} customError.ErrSegmentValidationError "Ошибка валидации данных запроса"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/segments [post]
func (r *segmentRoutes) create(c echo.Context) error {
	var input service.SegmentCreateInput

	if err := c.Bind(&input); err != nil {
		return errorHandler(c, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError: err,
			Comment:     err.Error(),
			Location:    "Router - segmentRoutes.create - c.Bind",
		}})
	}

	name, err := r.segmentService.CreateSegment(c.Request().Context(), input)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusCreated, CreateResponse{
		Name: name,
	})
}

type GetAllSegmentsResponse struct {
	Segments []entity.Segment `json:"segments"`
}

// @Summary Получить список всех сегментов
// @Description Возвращает список всех сегментов
// @Tags segments
// @Produce json
// @Success 200 {object} GetAllSegmentsResponse "Список всех сегментов"
// @Failure 400 {object} customError.ErrSegmentValidationError "Ошибка валидации данных запроса"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/segments [get]
func (r *segmentRoutes) getAll(c echo.Context) error {
	sType := c.QueryParam("segment_type")
	switch sType {
	case "":
		sType = "both"
	case "both":
	case "alive":
	case "deleted":
	default:
		return errorHandler(c, customError.ErrSegmentValidationError{ErrBase: customError.ErrBase{
			OriginError:     nil,
			OriginErrorText: "",
			Comment:         "Invalid \"segment_type\" param was given, valid values: [\"alive\", \"deleted\", \"both\"]",
			Location:        "SegmentRoutes.getAll - c.QueryParam",
		}})
	}

	var sTypeInt int
	switch sType {
	case "both":
		sTypeInt = 2
	case "alive":
		sTypeInt = 0
	case "deleted":
		sTypeInt = 1
	}

	segments, err := r.segmentService.GetAllSegments(c.Request().Context(), sTypeInt)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, GetAllSegmentsResponse{Segments: segments})
}

type GetSegmentByNameResponse struct {
	Segment entity.Segment `json:"segment"`
}

// @Summary Получить сегмент с указанным именем
// @Description Возвращает информацию о сегменте с указанным именем
// @Tags segments
// @Produce json
// @Param name path string true "Наименование сегмента"
// @Success 200 {object} GetSegmentByNameResponse "Сегмент с указанным именем"
// @Failure 400 {object} customError.ErrSegmentValidationError "Ошибка валидации данных запроса"
// @Failure 404 {object} customError.ErrSegmentNotFound "Сегмент с указанным именем не был найден"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/segments/{name} [get]
func (r *segmentRoutes) getByName(c echo.Context) error {
	name := c.Param("name")

	segment, err := r.segmentService.GetSegmentByName(c.Request().Context(), name)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, GetSegmentByNameResponse{segment})
}

type DeleteSegmentByNameResponse struct {
	Message string `json:"message" example:"successfully deleted segment \"AVITO_MUSIC_SERVICE\""`
}

// @Summary Удалить сегмент с указанным именем
// @Description Удаляет сегмент с указанным именем из системы.
// @Description В случае, если на момент совершения запроса в этот
// @Description сегмент входят какие-либо пользователи, они автоматически
// @Description выйдут из данного сегмента.
// @Tags segments
// @Accept json
// @Produce json
// @Param name path string true "Наименование сегмента"
// @Success 200 {object} DeleteSegmentByNameResponse "Сообщение об успехе"
// @Failure 400 {object} customError.ErrSegmentValidationError "Ошибка валидации данных запроса"
// @Failure 404 {object} customError.ErrSegmentNotFound "Сегмент с указанным именем не был найден"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/segments/{name} [delete]
func (r *segmentRoutes) deleteByName(c echo.Context) error {
	name := c.Param("name")

	err := r.segmentService.DeleteSegment(c.Request().Context(), name)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, DeleteSegmentByNameResponse{fmt.Sprintf("successfully deleted segment \"%s\"", name)})
}
