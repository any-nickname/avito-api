package v1

import (
	_ "avito-rest-api/internal/error"
	"avito-rest-api/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type reportRoutes struct {
	reportService service.Report
}

func newReportRoutes(g *echo.Group, reportService service.Report) {
	r := &reportRoutes{
		reportService: reportService,
	}

	g.GET("", r.makeReport)
}

// MakeReportResponse - структура ответа на запрос о создании отчёта
type MakeReportResponse struct {
	// Дата формирования отчёта
	ReportDate string `json:"report_date"`
	// Отчёт в виде csv-строки с разделителями "," и символом перехода на новую строку "\n"
	Report string `json:"report"`
}

// @Summary Получить отчёт в формате csv
// @Description Возвращает отчёт,
// @Description содержащий столбцы `user_id`, `segment_name`, `start_date`,
// @Description `end_date`, обозначающие идентификатор пользователя,
// @Description наименование сегмента, дату добавления пользователя в сегмент и
// @Description дату выхода пользователя из сегмента соответственно. Строки отчёта
// @Description отсортированы в порядке возрастания по дате добавления пользователя в сегмент.
// @Tags reports
// @Success 200 {object} MakeReportResponse "Структура, содержащая дату формирования отчёта и ссылку на файл с отчётом или отчёт в виде csv-строки"
// @Failure 400 {object} error.ErrReportValidationError "Ошибка валидации данных запроса"
// @Failure 500 {object} error.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/reports [get]
func (r *reportRoutes) makeReport(c echo.Context) error {
	result, err := r.reportService.MakeReport(c.Request().Context())
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, MakeReportResponse{
		ReportDate: result.ReportDate,
		Report:     result.Report,
	})
}
