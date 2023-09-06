package v1

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/internal/service"
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

	// Валидация
	// 1. Указано наименование сегмента
	if input.Name == "" {
		return errorHandler(c, customError.ErrSegmentValidationError{ErrBase: customError.ErrBase{
			Comment:  "Field \"name\" was not provided",
			Location: "SegmentRoutes.create - validation",
		}})
	}
	// 2. Длина имени сегмента до 1000 символов
	if len(input.Name) > 1000 {
		return errorHandler(c, customError.ErrSegmentValidationError{ErrBase: customError.ErrBase{
			Comment:  "Field \"name\" cannot have length over 1000 symbols",
			Location: "SegmentRoutes.create - validation",
		}})
	}
	// 3. percentage, если есть, является положительным целым числом до 100
	if input.PercentageOfUsersAdded < 0 || input.PercentageOfUsersAdded > 100 {
		return errorHandler(c, customError.ErrSegmentValidationError{ErrBase: customError.ErrBase{
			Comment:  "Field \"percentage\", if provided, must be a positive integer number from range [1, 100]",
			Location: "SegmentRoutes.create - validation",
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
// @Param segment_type query string false "Параметр, определяющий, сегменты какого типа (живые и(или) удалённые) необходимо вернуть. Значение `both` предполагает, что будут возвращены сегменты обоих типов (то есть абсолютно все сегменты, когда-либо созданные в системе). Значение `alive` предполагает, что будут возвращены только живые (то есть не помеченные как удалённые) сегменты. Значение `deleted` предполагает, что будут возвращены только сегменты, помеченные как удалённые. Отсутствие параметра равносильно параметру со значением `both`."
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
			Comment:  "Invalid \"segment_type\" param was given, valid values: [\"alive\", \"deleted\", \"both\"]",
			Location: "SegmentRoutes.getAll - c.QueryParam",
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

	// Валидация
	if name == "" {
		return errorHandler(c, customError.ErrSegmentValidationError{ErrBase: customError.ErrBase{
			Comment:  "Invalid query param \"name\" was provided. \"name\" cannot be empty",
			Location: "SegmentRoutes.getByName - validation",
		}})
	}

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
