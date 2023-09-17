package v1

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/internal/service"
	mock_service "avito-rest-api/internal/service/mocks"
	"bytes"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestUserRoutes_getAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type MockBehaviour func(m *mock_service.MockUser, args args)

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
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().GetAllUsers(args.ctx).Return([]entity.User{
					{
						ID:       1,
						Name:     "Дмитрий",
						Lastname: "Поплавский",
						Sex:      0,
						SexText:  "мужской",
						Age:      45,
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"users":[{"user_id":1,"name":"Дмитрий","lastname":"Поплавский","sex":0,"sex_text":"мужской","age":45,"is_deleted":false}]}` + "\n",
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
			req := httptest.NewRequest(http.MethodGet, "/users", nil)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUserRoutes_getAllWithSegments(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type MockBehaviour func(m *mock_service.MockUser, args args)

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
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().GetAllUsersWithSegments(args.ctx).Return([]entity.UserWithSegments{
					{
						User: entity.User{
							ID:       1,
							Name:     "Дмитрий",
							Lastname: "Поплавский",
							Sex:      0,
							SexText:  "мужской",
							Age:      45,
						},
						Segments: []entity.UserSegmentInformation{
							{
								InfoID:    1,
								UserID:    1,
								SegmentID: 1,
								Name:      "AVITO_BAKERY",
								StartDate: "10:05:24 07.09.2023",
							},
						},
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"users":[{"user":{"user_id":1,"name":"Дмитрий","lastname":"Поплавский","sex":0,"sex_text":"мужской","age":45,"is_deleted":false},"segments":[{"information_id":1,"user_id":1,"segment_id":1,"name":"AVITO_BAKERY","start_date":"10:05:24 07.09.2023","end_date":""}]}]}` + "\n",
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
			req := httptest.NewRequest(http.MethodGet, "/users/withSegments", nil)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUserRoutes_getByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	type MockBehaviour func(m *mock_service.MockUser, args args)

	testCases := []struct {
		name                 string
		args                 args
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				idInt, _ := strconv.Atoi(args.id)
				m.EXPECT().GetUserByID(args.ctx, idInt).Return(entity.User{
					ID:        1,
					Name:      "Дмитрий",
					Lastname:  "Поплавский",
					Sex:       0,
					SexText:   "мужской",
					Age:       45,
					IsDeleted: false,
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user":{"user_id":1,"name":"Дмитрий","lastname":"Поплавский","sex":0,"sex_text":"мужской","age":45,"is_deleted":false}}` + "\n",
		},
		{
			name: "Invalid id path param",
			args: args{
				ctx: context.Background(),
				id:  "sobaka-barabaka",
			},
			mockBehaviour:      func(m *mock_service.MockUser, args args) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"strconv.Atoi: parsing \"%s\": invalid syntax","title":"ErrUserValidationError","comment":"Invalid \"id\" path-param value. \"id\" should be integer","location":"UserRoutes.getByID - strconv.Atoi"}`,
				"sobaka-barabaka",
			) + "\n",
		},
		{
			name: "User with given id does not exist",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				idInt, _ := strconv.Atoi(args.id)
				m.EXPECT().GetUserByID(args.ctx, idInt).Return(entity.User{}, customError.ErrUserNotFound{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("User with id %d not found", idInt),
					Location: "UserRepository.GetUserByID",
				}})
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserNotFound","comment":"User with id %d not found","location":"UserRepository.GetUserByID"}`,
				1,
			) + "\n",
		},
		{
			name: "User with given id is deleted",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				idInt, _ := strconv.Atoi(args.id)
				m.EXPECT().GetUserByID(args.ctx, idInt).Return(entity.User{}, customError.ErrUserDeleted{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("User with id %d is deleted", idInt),
					Location: "UserService.GetUserByID - us.userRepository.GetUserByID",
				}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserDeleted","comment":"User with id %d is deleted","location":"UserService.GetUserByID - us.userRepository.GetUserByID"}`,
				1,
			) + "\n",
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
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", tc.args.id), nil)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUserRoutes_getByIWithSegments(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	type MockBehaviour func(m *mock_service.MockUser, args args)

	testCases := []struct {
		name                 string
		args                 args
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				idInt, _ := strconv.Atoi(args.id)
				m.EXPECT().GetUserWithSegmentsByUserID(args.ctx, idInt).Return(entity.UserWithSegments{
					User: entity.User{
						ID:        1,
						Name:      "Дмитрий",
						Lastname:  "Поплавский",
						Sex:       0,
						SexText:   "мужской",
						Age:       45,
						IsDeleted: false,
					},
					Segments: []entity.UserSegmentInformation{
						{
							InfoID:    1,
							UserID:    1,
							SegmentID: 1,
							Name:      "AVITO_BAKERY",
							StartDate: "10:05:24 07.09.2023",
							EndDate:   "",
						},
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user":{"user":{"user_id":1,"name":"Дмитрий","lastname":"Поплавский","sex":0,"sex_text":"мужской","age":45,"is_deleted":false},"segments":[{"information_id":1,"user_id":1,"segment_id":1,"name":"AVITO_BAKERY","start_date":"10:05:24 07.09.2023","end_date":""}]}}` + "\n",
		},
		{
			name: "Invalid id path param",
			args: args{
				ctx: context.Background(),
				id:  "sobaka-barabaka",
			},
			mockBehaviour:      func(m *mock_service.MockUser, args args) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"strconv.Atoi: parsing \"%s\": invalid syntax","title":"ErrUserValidationError","comment":"Invalid \"id\" path-param value. \"id\" should be integer","location":"UserRoutes.getByID - strconv.Atoi"}`,
				"sobaka-barabaka",
			) + "\n",
		},
		{
			name: "User with given id does not exist",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				idInt, _ := strconv.Atoi(args.id)
				m.EXPECT().GetUserWithSegmentsByUserID(args.ctx, idInt).Return(entity.UserWithSegments{}, customError.ErrUserNotFound{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("User with id %d not found", idInt),
					Location: "UserRepository.GetUserByID",
				}})
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserNotFound","comment":"User with id %d not found","location":"UserRepository.GetUserByID"}`,
				1,
			) + "\n",
		},
		{
			name: "User with given id is deleted",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				idInt, _ := strconv.Atoi(args.id)
				m.EXPECT().GetUserWithSegmentsByUserID(args.ctx, idInt).Return(entity.UserWithSegments{}, customError.ErrUserDeleted{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("User with id %d is deleted", idInt),
					Location: "UserService.GetUserWithSegmentsByUserID - us.userRepository.GetUserByID",
				}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserDeleted","comment":"User with id %d is deleted","location":"UserService.GetUserWithSegmentsByUserID - us.userRepository.GetUserByID"}`,
				1,
			) + "\n",
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
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s/withSegments", tc.args.id), nil)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUserRoutes_addUserToSegments(t *testing.T) {
	type argsInput struct {
		UserID   int
		Segments []entity.UserSegmentInformation
	}

	type args struct {
		ctx   context.Context
		input argsInput
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
				input: argsInput{
					UserID: 1,
					Segments: []entity.UserSegmentInformation{
						{
							Name:    "AVITO_BAKERY",
							EndDate: "10:00:00 25.09.2023",
						},
					},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY","end_date":"10:00:00 25.09.2023"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().AddUserToSegments(args.ctx, args.input.UserID, args.input.Segments).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"message":"user %d was successfully added to the segments"}`, 1) + "\n",
		},
		{
			name: "User with given id does not exist",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID: 1,
					Segments: []entity.UserSegmentInformation{
						{
							Name:    "AVITO_BAKERY",
							EndDate: "10:00:00 25.09.2023",
						},
					},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY","end_date":"10:00:00 25.09.2023"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().AddUserToSegments(args.ctx, args.input.UserID, args.input.Segments).Return(customError.ErrUserNotFound{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("Operation was canceled because user with id = %d does not exist", args.input.UserID),
					Location: "UserService.AddUserToSegments - us.GetUserByID",
				}})
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserNotFound","comment":"Operation was canceled because user with id = %d does not exist","location":"UserService.AddUserToSegments - us.GetUserByID"}`,
				1,
			) + "\n",
		},
		{
			name: "User with given id is deleted",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID: 1,
					Segments: []entity.UserSegmentInformation{
						{
							Name:    "AVITO_BAKERY",
							EndDate: "10:00:00 25.09.2023",
						},
					},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY","end_date":"10:00:00 25.09.2023"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().AddUserToSegments(args.ctx, args.input.UserID, args.input.Segments).Return(customError.ErrUserDeleted{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("Operation was canceled because user with id = %d is deleted", args.input.UserID),
					Location: "UserService.AddUserToSegments - us.GetUserByID",
				}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserDeleted","comment":"Operation was canceled because user with id = %d is deleted","location":"UserService.AddUserToSegments - us.GetUserByID"}`,
				1,
			) + "\n",
		},
		{
			name: "One of segments does not exist",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID: 1,
					Segments: []entity.UserSegmentInformation{
						{
							Name:    "AVITO_BAKERY",
							EndDate: "10:00:00 25.09.2023",
						},
					},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY","end_date":"10:00:00 25.09.2023"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().AddUserToSegments(args.ctx, args.input.UserID, args.input.Segments).Return(customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
					Comment: fmt.Sprintf(
						"Operation was canceled. Failed to add user (id=%d) to the segment \"%s\" because segment does not exist",
						args.input.UserID,
						args.input.Segments[0].Name,
					),
					Location: "UserService.AddUserToSegments - us.segmentRepository.GetSegmentByName",
				}})
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrSegmentNotFound","comment":"Operation was canceled. Failed to add user (id=%d) to the segment \"%s\" because segment does not exist","location":"UserService.AddUserToSegments - us.segmentRepository.GetSegmentByName"}`,
				1,
				"AVITO_BAKERY",
			) + "\n",
		},
		{
			name: "One of segments is deleted",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID: 1,
					Segments: []entity.UserSegmentInformation{
						{
							Name:    "AVITO_BAKERY",
							EndDate: "10:00:00 25.09.2023",
						},
					},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY","end_date":"10:00:00 25.09.2023"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().AddUserToSegments(args.ctx, args.input.UserID, args.input.Segments).Return(
					customError.ErrSegmentDeleted{ErrBase: customError.ErrBase{
						Comment: fmt.Sprintf("Operation was canceled. Failed to add user (id=%d) to the "+
							"segment \"%s\" because segment does not exist "+
							"(segment was deleted earlier and was not created again)",
							args.input.UserID,
							args.input.Segments[0].Name,
						),
						Location: "UserService.AddUserToSegments - isDeleted",
					}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrSegmentDeleted","comment":"Operation was canceled. Failed to add user (id=%d) to the segment \"%s\" because segment does not exist (segment was deleted earlier and was not created again)","location":"UserService.AddUserToSegments - isDeleted"}`,
				1,
				"AVITO_BAKERY",
			) + "\n",
		},
		{
			name: "Appearance of two identical segments",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID: 1,
					Segments: []entity.UserSegmentInformation{
						{
							Name: "AVITO_MUSIC",
						},
						{
							Name: "AVITO_MUSIC",
						},
					},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_MUSIC"},{"name":"AVITO_MUSIC"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().AddUserToSegments(args.ctx, args.input.UserID, args.input.Segments).Return(
					customError.ErrUserValidationError{ErrBase: customError.ErrBase{
						Comment: fmt.Sprintf(
							"Operation was canceled. The segment \"%s\" occurs more than 1 time in the list",
							args.input.Segments[1].Name,
						),
						Location: "UserService.AddUserToSegments - validation",
					}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserValidationError","comment":"Operation was canceled. The segment \"%s\" occurs more than 1 time in the list","location":"UserService.AddUserToSegments - validation"}`,
				"AVITO_MUSIC",
			) + "\n",
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
			req := httptest.NewRequest(http.MethodPost, "/users/addUserToSegments", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUserRoutes_deleteUserFromSegments(t *testing.T) {
	type argsInput struct {
		UserID   int
		Segments []string
	}

	type args struct {
		ctx   context.Context
		input argsInput
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
				input: argsInput{
					UserID:   1,
					Segments: []string{"AVITO_BAKERY"},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().DeleteUserFromSegments(args.ctx, args.input.UserID, args.input.Segments).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"message":"user %d was successfully removed from segments"}`, 1) + "\n",
		},
		{
			name: "User with given id does not exist",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID:   1,
					Segments: []string{"AVITO_BAKERY"},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().DeleteUserFromSegments(args.ctx, args.input.UserID, args.input.Segments).Return(customError.ErrUserNotFound{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("Operation was canceled. User with id = %d does not exist", args.input.UserID),
					Location: "UserService.DeleteUserFromSegments - us.GetUserByID",
				}})
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserNotFound","comment":"Operation was canceled. User with id = %d does not exist","location":"UserService.DeleteUserFromSegments - us.GetUserByID"}`,
				1,
			) + "\n",
		},
		{
			name: "User with given id is deleted",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID:   1,
					Segments: []string{"AVITO_BAKERY"},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().DeleteUserFromSegments(args.ctx, args.input.UserID, args.input.Segments).Return(customError.ErrUserDeleted{ErrBase: customError.ErrBase{
					Comment:  fmt.Sprintf("Operation was canceled because user with id = %d is deleted", args.input.UserID),
					Location: "UserService.DeleteUserFromSegments - us.GetUserByID",
				}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserDeleted","comment":"Operation was canceled because user with id = %d is deleted","location":"UserService.DeleteUserFromSegments - us.GetUserByID"}`,
				1,
			) + "\n",
		},
		{
			name: "One of segments does not exist",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID:   1,
					Segments: []string{"AVITO_BAKERY"},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().DeleteUserFromSegments(args.ctx, args.input.UserID, args.input.Segments).Return(customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
					Comment: fmt.Sprintf(
						"Operation was canceled. Failed to add user (id=%d) to the segment \"%s\" because segment does not exist",
						args.input.UserID,
						args.input.Segments[0],
					),
					Location: "UserService.DeleteUserFromSegments - us.segmentRepository.GetSegmentByName",
				}})
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrSegmentNotFound","comment":"Operation was canceled. Failed to add user (id=%d) to the segment \"%s\" because segment does not exist","location":"UserService.DeleteUserFromSegments - us.segmentRepository.GetSegmentByName"}`,
				1,
				"AVITO_BAKERY",
			) + "\n",
		},
		{
			name: "One of segments is deleted",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID:   1,
					Segments: []string{"AVITO_BAKERY"},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_BAKERY"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().DeleteUserFromSegments(args.ctx, args.input.UserID, args.input.Segments).Return(
					customError.ErrSegmentDeleted{ErrBase: customError.ErrBase{
						Comment: fmt.Sprintf("Operation was canceled. Failed to add user (id=%d) to the "+
							"segment \"%s\" because segment does not exist "+
							"(segment was deleted earlier and was not created again)",
							args.input.UserID,
							args.input.Segments[0],
						),
						Location: "UserService.DeleteUserFromSegments - isDeleted",
					}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrSegmentDeleted","comment":"Operation was canceled. Failed to add user (id=%d) to the segment \"%s\" because segment does not exist (segment was deleted earlier and was not created again)","location":"UserService.DeleteUserFromSegments - isDeleted"}`,
				1,
				"AVITO_BAKERY",
			) + "\n",
		},
		{
			name: "Appearance of two identical segments",
			args: args{
				ctx: context.Background(),
				input: argsInput{
					UserID:   1,
					Segments: []string{"AVITO_MUSIC", "AVITO_MUSIC"},
				},
			},
			inputBody: `{"id":1,"segments":[{"name":"AVITO_MUSIC"},{"name":"AVITO_MUSIC"}]}`,
			mockBehaviour: func(m *mock_service.MockUser, args args) {
				m.EXPECT().DeleteUserFromSegments(args.ctx, args.input.UserID, args.input.Segments).Return(
					customError.ErrUserValidationError{ErrBase: customError.ErrBase{
						Comment: fmt.Sprintf(
							"Operation was canceled. The segment \"%s\" occurs more than 1 time in the list",
							args.input.Segments[1],
						),
						Location: "UserService.DeleteUserFromSegments - validation",
					}})
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"origin_error_text":"","title":"ErrUserValidationError","comment":"Operation was canceled. The segment \"%s\" occurs more than 1 time in the list","location":"UserService.DeleteUserFromSegments - validation"}`,
				"AVITO_MUSIC",
			) + "\n",
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
			req := httptest.NewRequest(http.MethodPost, "/users/deleteUserFromSegments", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// Выполнение запроса
			e.ServeHTTP(w, req)

			// Проверка ответа
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
