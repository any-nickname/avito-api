package service

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/internal/repository"
	"context"
	"fmt"
	"strings"
	"time"
)

type UserService struct {
	userRepository    repository.User
	segmentRepository repository.Segment
}

func NewUserService(userRepository repository.User, segmentRepository repository.Segment) *UserService {
	return &UserService{userRepository: userRepository, segmentRepository: segmentRepository}
}

// UserCreateInput - DTO для получения данных
// для создания пользователя из тела запроса.
type UserCreateInput struct {
	Name     string `json:"name" example:"Михаил" validate:"required"`
	Lastname string `json:"lastname" example:"Иванов" validate:"required"`
	Sex      int    `json:"sex" example:"0" enums:"0,1" validate:"required"` // Пол, 0 - мужской, 1 - женский
	Age      int    `json:"age" example:"27" validate:"required"`            // Целое положительное число
}

func (us *UserService) CreateUser(ctx context.Context, input UserCreateInput) (int, error) {
	// Валидация
	if input.Name == "" {
		return 0, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     nil,
			OriginErrorText: "",
			Comment:         "\"name\" field cannot be empty",
			Location:        "UserService.CreateUser",
		}}
	}
	if input.Lastname == "" {
		return 0, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     nil,
			OriginErrorText: "",
			Comment:         "\"lastname\" field cannot be empty",
			Location:        "UserService.CreateUser",
		}}
	}
	if input.Sex < 0 || input.Sex > 1 {
		return 0, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     nil,
			OriginErrorText: "",
			Comment:         "\"sex\" field can only take value from range [0, 1]",
			Location:        "UserService.CreateUser",
		}}
	}
	if input.Age < 0 {
		return 0, customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     nil,
			OriginErrorText: "",
			Comment:         "\"age\" field cannot be negative",
			Location:        "UserService.CreateUser",
		}}
	}

	// Маппинг данных из DTO в сущность User
	user := entity.User{
		ID:       0,
		Name:     input.Name,
		Lastname: input.Lastname,
		Sex:      input.Sex,
		SexText:  "",
		Age:      input.Age,
	}
	return us.userRepository.CreateUser(ctx, user)
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return us.userRepository.GetAllUsers(ctx)
}

func (us *UserService) GetAllUsersWithSegments(ctx context.Context) ([]entity.UserWithSegments, error) {
	users, err := us.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var usersWithSegments []entity.UserWithSegments

	for _, user := range users {
		userSegments, err := us.userRepository.GetUserSegmentsByUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		usersWithSegments = append(usersWithSegments, entity.UserWithSegments{
			User:     user,
			Segments: userSegments,
		})
	}

	return usersWithSegments, nil
}

func (us *UserService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	user, err := us.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}
	if user.IsDeleted {
		return entity.User{}, customError.ErrUserDeleted{ErrBase: customError.ErrBase{
			Comment:  fmt.Sprintf("User with id %d is deleted", id),
			Location: "UserService.GetUserByID - us.userRepository.GetUserByID",
		}}
	}
	return user, nil
}

func (us *UserService) GetUserSegmentsByUserID(ctx context.Context, id int) ([]entity.UserSegmentInformation, error) {
	return us.userRepository.GetUserSegmentsByUserID(ctx, id)
}

func (us *UserService) GetUserWithSegmentsByUserID(ctx context.Context, id int) (entity.UserWithSegments, error) {
	userWithSegments := entity.UserWithSegments{}

	user, err := us.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return entity.UserWithSegments{}, err
	}
	if user.IsDeleted {
		return entity.UserWithSegments{}, customError.ErrUserDeleted{ErrBase: customError.ErrBase{
			Comment:  fmt.Sprintf("User with id %d is deleted", id),
			Location: "UserService.GetUserWithSegmentsByUserID - us.userRepository.GetUserByID",
		}}
	}
	segments, err := us.userRepository.GetUserSegmentsByUserID(ctx, id)
	if err != nil {
		return entity.UserWithSegments{}, err
	}

	userWithSegments.User = user
	userWithSegments.Segments = segments

	return userWithSegments, nil
}

func (us *UserService) AddUserToSegments(ctx context.Context, id int, segments []entity.UserSegmentInformation) error {
	// Валидация времени, переданного в сегментах
	for _, s := range segments {
		if s.EndDate != "" {
			if _, err := time.Parse("15:04:05 02.01.2006", s.EndDate); err != nil {
				return customError.ErrUserValidationError{ErrBase: customError.ErrBase{
					OriginError:     err,
					OriginErrorText: err.Error(),
					Comment:         fmt.Sprintf("Operation was canceled. Invalid \"end_date\" = %s was provided", s.EndDate),
					Location:        "UserService.AddUserToSegments - time.Parse",
				}}
			}
		}
	}

	// Проверка существования пользователя и проверка на то, что он удалён
	_, err := us.GetUserByID(ctx, id)
	if err != nil {
		if _, ok := err.(customError.ErrUserNotFound); ok {
			return customError.ErrUserNotFound{ErrBase: customError.ErrBase{
				OriginError:     nil,
				OriginErrorText: "",
				Comment:         fmt.Sprintf("Operation was canceled because user with id = %d does not exist", id),
				Location:        "UserService.AddUserToSegments - us.GetUserByID",
			}}
		}
		if _, ok := err.(customError.ErrUserDeleted); ok {
			return customError.ErrUserDeleted{ErrBase: customError.ErrBase{
				OriginError:     nil,
				OriginErrorText: "",
				Comment:         fmt.Sprintf("Operation was canceled because user with id = %d is deleted", id),
				Location:        "UserService.AddUserToSegments - us.GetUserByID",
			}}
		}
		return err
	}

	// Проверка существования сегментов
	for i, segment := range segments {
		// Проверка на то, что сегмент существует и не удалён
		tSegment, err := us.segmentRepository.GetSegmentByName(ctx, segment.Name)
		if err != nil {
			if _, ok := err.(customError.ErrSegmentNotFound); ok {
				return customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
					Comment: fmt.Sprintf("Operation was canceled. Failed to add user (id=%d) to the"+
						"segment \"%s\" because segment does not exist", id, tSegment.Name),
					Location: "UserService.AddUserToSegments - us.segmentRepository.GetSegmentByName",
				}}
			}
			return err
		}
		if tSegment.IsDeleted {
			return customError.ErrSegmentDeleted{ErrBase: customError.ErrBase{
				Comment: fmt.Sprintf("Operation was canceled. Failed to add user (id=%d) to the "+
					"segment \"%s\" because segment does not exist "+
					"(segment was deleted earlier and was not created again)", id, tSegment.Name),
				Location: "UserService.AddUserToSegments - isDeleted",
			}}
		}
		segments[i].SegmentID = tSegment.ID
	}

	// Получим сегменты, в которые пользователь входит на текущий момент,
	// и проверим отсутствие пересечений со списком `segments` за O(N)
	// с помощью мапы
	userSegmentsMap := make(map[string]bool)
	var intersection []string

	userSegments, err := us.userRepository.GetUserSegmentsByUserID(ctx, id)
	if err != nil {
		return err
	}
	for _, segment := range userSegments {
		userSegmentsMap[segment.Name] = true
	}
	for _, segment := range segments {
		if userSegmentsMap[segment.Name] {
			intersection = append(intersection, segment.Name)
		}
	}

	// Если есть пересечение, создадим ошибку
	if len(intersection) > 0 {
		return customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     nil,
			OriginErrorText: "",
			Comment: fmt.Sprintf(
				"Operation was canceled. "+
					"Failed to add segments: [%s] to user (id = %d) "+
					"since user already has this segments",
				strings.Join(intersection, ", "),
				id,
			),
			Location: "UserService.AddUserToSegments - us.userRepository.GetUserSegmentsByUserID",
		}}
	}

	// В противном случае добавляем пользователя в указанные сегменты
	return us.userRepository.AddUserToSegments(ctx, id, segments)
}

func (us *UserService) DeleteUserFromSegments(ctx context.Context, id int, segments []string) error {
	// Проверим существование юзера
	_, err := us.GetUserByID(ctx, id)
	if err != nil {
		if _, ok := err.(customError.ErrUserNotFound); ok {
			return customError.ErrUserNotFound{ErrBase: customError.ErrBase{
				OriginError:     nil,
				OriginErrorText: "",
				Comment:         fmt.Sprintf("Operation was canceled. User with id = %d does not exist", id),
				Location:        "UserService.DeleteUserFromSegments - us.GetUserByID",
			}}
		}
		return err
	}

	// Проверим существование переданных сегментов
	for _, name := range segments {
		segment, err := us.segmentRepository.GetSegmentByName(ctx, name)
		if err != nil {
			if _, ok := err.(customError.ErrSegmentNotFound); ok {
				return customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
					OriginError:     nil,
					OriginErrorText: "",
					Comment:         fmt.Sprintf("Operation was canceled. Request contains segment %s which does not exist", name),
					Location:        "UserService.DeleteUserFromSegments - us.segmentRepository.GetSegmentByName",
				}}
			}
			return err
		}
		if segment.IsDeleted {
			return customError.ErrSegmentDeleted{ErrBase: customError.ErrBase{
				OriginError:     nil,
				OriginErrorText: "",
				Comment: fmt.Sprintf("Operation was canceled. Request contains segment \"%s\" which "+
					"does not exist (this segment was deleted earlier and was not created again)", name),
				Location: "UserService.DeleteUserFromSegments - us.segmentRepository.GetSegmentByName",
			}}
		}
	}

	// Проверим, что пользователь входит в переданные сегменты
	userSegments, err := us.GetUserSegmentsByUserID(ctx, id)
	if err != nil {
		return err
	}
	userSegmentsMap := make(map[string]bool)
	var exclusion []string
	for _, segmentInfo := range userSegments {
		userSegmentsMap[segmentInfo.Name] = true
	}
	for _, name := range segments {
		if !userSegmentsMap[name] {
			exclusion = append(exclusion, name)
		}
	}

	// Если есть исключения (т.е. сегменты, в которые пользователь не входит), создадим ошибку
	if len(exclusion) > 0 {
		return customError.ErrUserValidationError{ErrBase: customError.ErrBase{
			OriginError:     nil,
			OriginErrorText: "",
			Comment: fmt.Sprintf(
				"Operation was canceled. User (id = %d) does not belong to some "+
					"segments: [%s] which were given in the request body",
				id,
				strings.Join(exclusion, ", "),
			),
			Location: "UserService.DeleteUserFromSegments - exclusion",
		}}
	}

	// Преобразование сегментов в корректные сущности
	var segmentsToDelete []entity.UserSegmentInformation
	segmentsToDeleteMap := make(map[string]bool)
	for _, name := range segments {
		segmentsToDeleteMap[name] = true
	}
	for _, segmentInfo := range userSegments {
		if segmentsToDeleteMap[segmentInfo.Name] {
			segmentsToDelete = append(segmentsToDelete, segmentInfo)
		}
	}

	return us.userRepository.DeleteUserFromSegments(ctx, id, segmentsToDelete)
}
