package v1

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/internal/service"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(g *echo.Group, userService service.User) {
	r := &userRoutes{
		userService: userService,
	}

	g.POST("", r.create)
	g.GET("", r.getAll)
	g.GET("/withSegments", r.getAllWithSegments)
	g.GET("/:id", r.getByID)
	g.GET("/:id/withSegments", r.getByIDWithSegments)
	g.POST("/addUserToSegments", r.addUserToSegments)
	g.POST("/deleteUserFromSegments", r.deleteUserFromSegments)
}

type UserCreateResponse struct {
	ID int `json:"id" example:"26"`
}

// @Summary Создать пользователя
// @Description Создаёт пользователя на основе информации в теле запроса
// @Tags users
// @Accept json
// @Produce json
// @Param data body service.UserCreateInput true "Структура с информацией о создаваемом пользователе"
// @Success 201 {object} UserCreateResponse "Идентификатор созданного пользователя"
// @Failure 400 {object} customError.ErrUserValidationError "Ошибка валидации данных запроса"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/users [post]
func (r *userRoutes) create(c echo.Context) error {
	var input service.UserCreateInput

	if err := c.Bind(&input); err != nil {
		return errorHandler(c, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "invalid request body",
			Location:        "UserRoutes.create - c.Bind",
		}})
	}

	id, err := r.userService.CreateUser(c.Request().Context(), service.UserCreateInput{
		Name:     input.Name,
		Lastname: input.Lastname,
		Sex:      input.Sex,
		Age:      input.Age,
	})

	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusCreated, UserCreateResponse{
		ID: id,
	})
}

type GetAllUsersResponse struct {
	Users []entity.User `json:"users"`
}

// @Summary Получить список всех пользователей
// @Description Возвращает список абсолютно всех пользователей
// @Tags users
// @Produce json
// @Success 200 {object} GetAllUsersResponse "Список всех пользователей"
// @Failure 400 {object} customError.ErrUserValidationError "Ошибка валидации данных запроса"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/users [get]
func (r *userRoutes) getAll(c echo.Context) error {
	users, err := r.userService.GetAllUsers(c.Request().Context())
	if err != nil {
		return errorHandler(c, err)
	}
	return c.JSON(http.StatusOK, GetAllUsersResponse{users})
}

type GetAllUsersWithSegmentsResponse struct {
	Users []entity.UserWithSegments `json:"users"`
}

// @Summary Получить список всех пользователей, включая их сегменты
// @Description Возвращает список всех пользователей, включая список активных для каждого пользователя сегментов на момент совершения запроса
// @Tags users
// @Produce json
// @Success 200 {object} GetAllUsersWithSegmentsResponse "Список пользователей с их активными сегментами"
// @Failure 400 {object} customError.ErrUserValidationError "Ошибка валидации данных запроса"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/users/withSegments [get]
func (r *userRoutes) getAllWithSegments(c echo.Context) error {
	usersWithSegments, err := r.userService.GetAllUsersWithSegments(c.Request().Context())
	if err != nil {
		return errorHandler(c, err)
	}
	return c.JSON(http.StatusOK, GetAllUsersWithSegmentsResponse{usersWithSegments})
}

type GetUserByIDResponse struct {
	User entity.User `json:"user"`
}

// @Summary Получить пользователя по ID
// @Description Возвращает пользователя с указанным ID
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} GetUserByIDResponse "Пользователь с указанным ID"
// @Failure 400 {object} customError.ErrUserValidationError "Ошибка валидации данных запроса"
// @Failure 404 {object} customError.ErrUserNotFound "Пользователь с указанным ID не был найден"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/users/{id} [get]
func (r *userRoutes) getByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorHandler(c, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Invalid \"id\" param value. \"id\" should be integer",
			Location:        "UserRoutes.getByID - strconv.Atoi - c.Param",
		}})
	}

	user, err := r.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, GetUserByIDResponse{User: user})
}

type GetUserByIDWithSegmentsResponse struct {
	User entity.UserWithSegments `json:"user"`
}

// @Summary Получить пользователя с его сегментами по ID
// @Description Возвращает пользователя с указанным ID, включая в тело ответа список сегментов, в которые пользователь входит на момент совершения запроса
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} GetUserByIDWithSegmentsResponse "Пользователь с его активными сегментами"
// @Failure 400 {object} customError.ErrUserValidationError "Ошибка валидации данных запроса"
// @Failure 404 {object} customError.ErrUserNotFound "Пользователь с указанным ID не был найден"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/users/{id}/withSegments [get]
func (r *userRoutes) getByIDWithSegments(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorHandler(c, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Invalid \"id\" param value. \"id\" should be integer",
			Location:        "UserRoutes.getByID - strconv.Atoi - c.Param",
		}})
	}

	userWithSegments, err := r.userService.GetUserWithSegmentsByUserID(c.Request().Context(), id)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, GetUserByIDWithSegmentsResponse{User: userWithSegments})
}

// AddUserToSegmentsInput - DTO для маппинга данных
// из запроса на добавление пользователя в определённые сегменты
type AddUserToSegmentsInput struct {
	ID       int `json:"id" example:"16" validate:"required"` // Идентификатор пользователя
	Segments []struct {
		Name    string `json:"name" example:"AVITO_MUSIC_SERVICE" validate:"required"` // Наименование сегмента, в который необходимо добавить пользователя
		EndDate string `json:"end_date" example:"10:00:00 25.09.2023"`                 // Необязательное поле, если отсутствует, значит дата выхода пользователя из сегмента не определена
	} `json:"segments"` // Сегменты, в которые необходимо добавить пользователя
}

type AddUserToSegmentsResponse struct {
	Message string `json:"message" example:"user 16 was successfully added to the segments"`
}

// @Summary Добавить пользователя в сегменты
// @Description Добавляет пользователя с указанным ID в указанные сегменты
// @Tags users
// @Param data body AddUserToSegmentsInput true "Структура, содержащая ID пользователя и наименование сегментов, в которые необходимо добавить пользователя. Поле `end_date` у сегмента является опциональным, и, если не  установлено, сигнализирует о том, что время выхода пользователя из сегмента не определено (пока сегмент  не будет удалён или пользователь не будет удалён из этого сегмента)"
// @Success 200 {object} AddUserToSegmentsResponse "Сообщение об успехе"
// @Failure 400 {object} customError.ErrUserValidationError "Ошибка валидации данных запроса"
// @Failure 404 {object} customError.ErrUserNotFound "Пользователь с указанным ID не был найден или некоторые из указанных сегментов не существуют"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/users/addUserToSegments [post]
func (r *userRoutes) addUserToSegments(c echo.Context) error {
	var segmentsInput AddUserToSegmentsInput
	err := c.Bind(&segmentsInput)
	if err != nil {
		return errorHandler(c, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to parse request's body",
			Location:        "UserRoutes.addUserToSegments - c.Bind",
		}})
	}

	var userID int = segmentsInput.ID
	var segments []entity.UserSegmentInformation

	for _, s := range segmentsInput.Segments {
		segments = append(segments, entity.UserSegmentInformation{
			Name:    s.Name,
			EndDate: s.EndDate,
		})
	}

	err = r.userService.AddUserToSegments(c.Request().Context(), userID, segments)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, AddUserToSegmentsResponse{Message: fmt.Sprintf(
		"user %d was successfully added to the segments",
		userID,
	)})
}

// DeleteUserFromSegmentsInput - DTO для маппинга данных из запроса на
// удаление пользователя из определённых сегментов
type DeleteUserFromSegmentsInput struct {
	ID       int `json:"id" example:"16" validate:"required"` // Идентификатор пользователя
	Segments []struct {
		Name string `json:"name" example:"AVITO_MUSIC_SERVICE" validate:"required"` // Наименование сегмента
	} `json:"segments"` // Сегменты, из которых необходимо удалить пользователя
}

type DeleteUserFromSegmentsResponse struct {
	Message string `json:"message" example:"user 179 was successfully removed from segments"`
}

// @Summary Удалить пользователя из сегментов
// @Description Удаляет пользователя с указанным ID из указанных сегментов
// @Tags users
// @Param data body DeleteUserFromSegmentsInput true "Структура, содержащая ID пользователя и наименования сегментов, из которых пользователя необходимо удалить"
// @Success 200 {object} DeleteUserFromSegmentsResponse "Сообщение об успехе"
// @Failure 400 {object} customError.ErrUserValidationError "Ошибка валидации данных запроса, может возникать, если пользователь не входит в указанные сегменты"
// @Failure 404 {object} customError.ErrUserNotFound "Пользователь с указанным ID не был найден или некоторые из указанных сегментов не существуют"
// @Failure 500 {object} customError.ErrInternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/users/deleteUserFromSegments [post]
func (r *userRoutes) deleteUserFromSegments(c echo.Context) error {
	var input DeleteUserFromSegmentsInput

	err := c.Bind(&input)
	if err != nil {
		return errorHandler(c, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Invalid request body",
			Location:        "UserRoutes.deleteUsersFromSegments",
		}})
	}

	var id int = input.ID

	var segments []string
	for _, s := range input.Segments {
		segments = append(segments, s.Name)
	}

	err = r.userService.DeleteUserFromSegments(c.Request().Context(), id, segments)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(http.StatusOK, DeleteUserFromSegmentsResponse{
		Message: fmt.Sprintf("user %d was successfully removed from segments", id),
	})
}
