package v1

import (
	"avito-rest-api/internal/service"
	mock_service "avito-rest-api/internal/service/mocks"
	"bytes"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserRoutes_create(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.UserCreateInput
	}

	type MockBehaviour func(m *mock_service.MockUser, args args)

	testCases := []struct {
		name                 string
		args                 args
		inputBody            string
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			args: args{
				ctx: context.Background(),
				input: service.UserCreateInput{
					Name:     "Дмитрий",
					Lastname: "Поплавский",
					Sex:      0,
					Age:      45,
				},
			},
			inputBody: `{"name":"Дмитрий","lastname":"Поплавский","sex":0,"age":45}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}` + "\n",
		},
		{
			name:                 "Invalid name: empty name",
			args:                 args{ctx: context.Background()},
			inputBody:            `{"lastname":"Поплавский","sex":0,"age":45}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Field \"name\" cannot be empty","location":"UserRoutes.create - validation"}` + "\n",
		},
		{
			name: "Invalid name: length over 1.000 symbols",
			args: args{ctx: context.Background()},
			inputBody: `{"name":"Абвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбв",
						"lastname":"Поплавский","sex":0,"age":45}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Length of the field \"name\" cannot be over 1.000 symbols","location":"UserRoutes.create - validation"}` + "\n",
		},
		{
			name:                 "Invalid lastname: empty lastname",
			args:                 args{ctx: context.Background()},
			inputBody:            `{"name":"Дмитрий","sex":0,"age":45}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Field \"lastname\" cannot be empty","location":"UserRoutes.create - validation"}` + "\n",
		},
		{
			name:                 "Invalid lastname: length over 1.000 symbols",
			args:                 args{ctx: context.Background()},
			inputBody:            `{"name":"Дмитрий","lastname":"Абвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбвбв","sex":0,"age":45}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Length of the field \"lastname\" cannot be over 1.000 symbols","location":"UserRoutes.create - validation"}` + "\n",
		},
		{
			name:                 "Invalid sex: empty sex",
			args:                 args{ctx: context.Background()},
			inputBody:            `{"name":"Дмитрий","lastname":"Поплавский","age":45}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Required field \"sex\" was not provided","location":"UserRoutes.create - validation"}` + "\n",
		},
		{
			name:                 "Invalid sex: value out of range [0, 1]",
			args:                 args{ctx: context.Background()},
			inputBody:            `{"name":"Дмитрий","lastname":"Поплавский","sex":4,"age":45}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Field \"sex\" must equals to 0 (man) or 1 (woman)","location":"UserRoutes.create - validation"}` + "\n",
		},
		{
			name:                 "Invalid age: empty age",
			args:                 args{ctx: context.Background()},
			inputBody:            `{"name":"Дмитрий","lastname":"Поплавский","sex":0}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Required field \"age\" was not provided","location":"UserRoutes.create - validation"}` + "\n",
		},
		{
			name:                 "Invalid age: less or equal to 0",
			args:                 args{ctx: context.Background()},
			inputBody:            `{"name":"Дмитрий","lastname":"Поплавский","sex":0,"age":0}`,
			mockBehaviour:        func(m *mock_service.MockUser, args args) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"origin_error_text":"","title":"ErrUserValidationError","comment":"Invalid request body. Field \"age\" must be positive integer number","location":"UserRoutes.create - validation"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Инициализация зависимостей
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Инициализация мока сервиса
			user := mock_service.NewMockUser(ctrl)
			tc.mockBehaviour(user, tc.args)
			services := &service.Services{User: user}

			// Создание тестового сервера
			e := echo.New()
			g := e.Group("/users")
			newUserRoutes(g, services.User)

			// Создание запроса
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
