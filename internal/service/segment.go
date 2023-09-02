package service

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/internal/repository"
	"fmt"
	"golang.org/x/net/context"
)

type SegmentService struct {
	segmentRepository repository.Segment
}

// NewSegmentService инициализирует сервис для сегментов
func NewSegmentService(segmentRepository repository.Segment) *SegmentService {
	return &SegmentService{segmentRepository: segmentRepository}
}

// SegmentCreateInput - DTO для маппинга данных из тела
// POST-запроса на создание сегмента.
type SegmentCreateInput struct {
	// Имя сегмента
	Name string `json:"name" example:"AVITO_MUSIC_SERVICE" validate:"required"`
	// Необязательное поле, процент пользователей, которое автоматически войдёт в сегмент при его создании
	PercentageOfUsersAdded int `json:"percentage" example:"57" minimum:"0" maximum:"100"`
}

// doesSegmentExist используется для проверки существования сегмента опираясь указанное название.
func (s *SegmentService) doesSegmentExist(ctx context.Context, name string) (bool, error) {
	_, err := s.segmentRepository.GetSegmentByName(ctx, name)
	if err != nil {
		if _, ok := err.(customError.ErrSegmentNotFound); ok {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// isSegmentDeleted проверяет, помечен ли сегмент с указанным именем как удалённый.
// isSegmentDeleted не осуществляет проверку на существование сегмента перед обращением
// к нему, для этого используется doesSegmentExist.
func (s *SegmentService) isSegmentDeleted(ctx context.Context, name string) (bool, error) {
	segment, err := s.segmentRepository.GetSegmentByName(ctx, name)
	if err != nil {
		return false, err
	}
	return segment.IsDeleted, nil
}

// CreateSegment используется для создания сегмента в базе данных и включает в себя проверки
// на дупликацию сегмента.
func (s *SegmentService) CreateSegment(ctx context.Context, input SegmentCreateInput) (string, error) {
	// Валидация
	if input.Name == "" {
		return "", customError.ErrSegmentValidationError{ErrBase: customError.ErrBase{
			OriginError: nil,
			Comment:     "Validation of segment's data failed, field \"name\" cannot be empty",
			Location:    "SegmentService.CreateSegment",
		}}
	}
	// Проверим, существует ли сегмент с таким же именем
	exist, err := s.doesSegmentExist(ctx, input.Name)
	if err != nil {
		return "", err
	}
	// Если существует и удалён - восстановим, в противном случае создадим ошибку
	if exist {
		isDeleted, err := s.isSegmentDeleted(ctx, input.Name)
		if err != nil {
			return "", err
		}
		if isDeleted {
			// Восстанавливаем сегмент
			return s.segmentRepository.RecoverSegment(ctx, input.Name)
		} else {
			return "", customError.ErrSegmentAlreadyExists{ErrBase: customError.ErrBase{
				Comment:  fmt.Sprintf("Segment with the given name \"%s\" already exists", input.Name),
				Location: "SegmentService.CreateSegment",
			}}
		}
	}
	// Если сегмент не существует, создадим с нуля
	name, err := s.segmentRepository.CreateSegment(ctx, entity.Segment{Name: input.Name})
	if err != nil {
		return "", err
	}
	// Если задан случайный процент пользователей для попадания в сегмент, добавим их
	if input.PercentageOfUsersAdded > 0 && input.PercentageOfUsersAdded <= 100 {
		err = s.segmentRepository.AddUsersToSegmentByRandomPercent(ctx, input.Name, input.PercentageOfUsersAdded)
		if err != nil {
			return "", err
		}
	}

	return name, nil
}

// GetAllSegments используется для получения всех существующих в системе сегментов, учитывая
// переданный параметр `sType` ("alive" - все сегменты, не помеченные как удалённые,
// "deleted" - все сегменты, помеченные как удалённые, "both" - абсолютно все сегменты).
func (s *SegmentService) GetAllSegments(ctx context.Context, sType int) ([]entity.Segment, error) {
	return s.segmentRepository.GetAllSegments(ctx, sType)
}

// GetSegmentByName используется для получения сегмента по имени и включает в себя проверки
// на существование сегмента.
func (s *SegmentService) GetSegmentByName(ctx context.Context, name string) (entity.Segment, error) {
	exist, err := s.doesSegmentExist(ctx, name)
	if err != nil {
		return entity.Segment{}, customError.ErrInternalServerError{}
	}
	if !exist {
		return entity.Segment{}, customError.ErrSegmentNotFound{}
	}

	segment, err := s.segmentRepository.GetSegmentByName(ctx, name)
	if err != nil {
		return entity.Segment{}, customError.ErrInternalServerError{}
	}

	return segment, nil
}

// DeleteSegment используется для удаления сегмента, проверяя, существует ли сегмент
// и не помечен ли он как удалённый.
func (s *SegmentService) DeleteSegment(ctx context.Context, name string) error {
	exist, err := s.doesSegmentExist(ctx, name)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError: err,
			Comment:     err.Error(),
			Location:    "SegmentService.DeleteSegment - s.doesSegmentExist",
		}}
	}
	if !exist {
		return customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
			OriginError: nil,
			Comment:     fmt.Sprintf("Unable to delete segment %s because it does not exist", name),
			Location:    "SegmentService.DeleteSegment - s.doesSegmentExist",
		}}
	}

	isDeleted, err := s.isSegmentDeleted(ctx, name)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError: err,
			Comment:     err.Error(),
			Location:    "SegmentService.DeleteSegment - s.doesSegmentExist",
		}}
	}
	if isDeleted {
		return customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
			OriginError: nil,
			Comment:     fmt.Sprintf("Unable to delete segment %s because it does not exist (segment was deleted earlier and was not created again)", name),
			Location:    "SegmentService.DeleteSegment - s.isSegmentDeleted",
		}}
	}

	err = s.segmentRepository.DeleteSegment(ctx, name)
	if err != nil {
		return err
	}
	return nil
}
