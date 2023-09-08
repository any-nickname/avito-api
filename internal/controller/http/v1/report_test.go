package v1

import (
	"avito-rest-api/internal/entity"
	"avito-rest-api/internal/service"
	mock_service "avito-rest-api/internal/service/mocks"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReportRoutes_makeReport(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type MockBehaviour func(m *mock_service.MockReport, args args)

	testCases := []struct {
		name                 string
		args                 args
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			args: args{ctx: context.Background()},
			mockBehaviour: func(m *mock_service.MockReport, args args) {
				m.EXPECT().MakeReport(args.ctx).Return(entity.ReportCSV{
					ReportDate: "19:52:04 02.09.2023",
					Report:     "user_id,segment_name,start_date,end_date\n1,AVITO_VOICE_MESSAGES,12:35:50 01.01.2023,\n",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"report_date":"19:52:04 02.09.2023","report":"user_id,segment_name,start_date,end_date\n1,AVITO_VOICE_MESSAGES,12:35:50 01.01.2023,\n"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Инициализация зависимостей
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Инициализация мока сервиса
			report := mock_service.NewMockReport(ctrl)
			tc.mockBehaviour(report, tc.args)
			services := &service.Services{Report: report}

			// Создание тестового сервера
			e := echo.New()
			g := e.Group("/reports")
			newReportRoutes(g, services.Report)

			// Создание запроса
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/reports", nil)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
