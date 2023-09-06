package pgdb

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/package/postgres"
	"context"
	"fmt"
	"strings"
	"time"
)

type SegmentRepository struct {
	*postgres.PostgreDB
}

// NewSegmentRepository инициализирует репозиторий `segment`, инкапсулирующий
// логику хранения сущности сегмента.
func NewSegmentRepository(pg *postgres.PostgreDB) *SegmentRepository {
	return &SegmentRepository{pg}
}

// CreateSegment добавляет в базу данных новый сегмент с указанным в `segment` именем
// и возвращает имя добавленного сегмента.
func (r *SegmentRepository) CreateSegment(ctx context.Context, segment entity.Segment) (string, error) {
	sql, args, _ := r.Builder.
		Insert("segments").
		Columns("name").
		Values(segment.Name).
		Suffix("RETURNING name").
		ToSql()

	var name string

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&name); err != nil {
		return "", customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError: err,
			Comment:     fmt.Sprintf("Failed to create or recover segment %s", segment.Name),
			Location:    "SegmentRepository.CreateSegment - r.Pool.QueryRow",
		}}
	}
	return name, nil
}

// GetAllSegments используется для получения всех сегментов указанного типа `sType` ("alive" - только
// не удалённые, "deleted" - только удалённые, "both" - все).
func (r *SegmentRepository) GetAllSegments(ctx context.Context, sType int) ([]entity.Segment, error) {
	var condition []interface{}
	switch sType {
	case 0: // только не удалённые
		condition = append(condition, false)
	case 1: // только удалённые
		condition = append(condition, true)
	case 2: // любые
		condition = append(condition, false, true)
	default:
		return nil, customError.ErrSegmentValidationError{
			ErrBase: customError.ErrBase{
				OriginError:     nil,
				OriginErrorText: "",
				Comment:         fmt.Sprintf("Failed to get all segments due to incorrect `sType` param which can only accept value from range [0, 1, 2]"),
				Location:        "SegmentRepository.GetAllSegments - sType validation",
			},
		}
	}

	var sqlPlaceholders []string
	for _ = range condition {
		sqlPlaceholders = append(sqlPlaceholders, "?")
	}

	sql, args, err := r.Builder.
		Select("*").
		From("segments").
		Where(fmt.Sprintf("is_deleted in (%s)", strings.Join(sqlPlaceholders, ", ")), condition...).
		ToSql()
	if err != nil {
		return nil, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to build sql query, inspect origin error text",
			Location:        "SegmentRepository.GetAllSegments - r.Builder",
		}}
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to query all segments from database, inspect origin error text",
			Location:        "SegmentRepository.GetAllSegments - r.Pool.Query",
		}}
	}

	var segments []entity.Segment

	for rows.Next() {
		var segment entity.Segment
		err := rows.Scan(
			&segment.ID,
			&segment.Name,
			&segment.IsDeleted,
		)
		if err != nil {
			return nil, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
				OriginError:     err,
				OriginErrorText: err.Error(),
				Comment:         "Failed to scan queried segments into segments list, inspect origin error text",
				Location:        "SegmentRepository.GetAllSegments - rows.Scan",
			}}
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

// GetSegmentByName используется для нахождения сегмента в базе данных по имени,
// не проверяя, существует ли сегмент с указанным именем (проверка реализуется
// на уровне сервиса).
func (r *SegmentRepository) GetSegmentByName(ctx context.Context, name string) (entity.Segment, error) {
	sql, args, _ := r.Builder.
		Select("segment_id, name, is_deleted").
		From("segments").
		Where("name = ?", name).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Segment{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to perform sql query",
			Location:        "SegmentRepository.GetSegmentByName - r.Pool.Query",
		}}
	}

	for rows.Next() {
		var segment entity.Segment
		err = rows.Scan(
			&segment.ID,
			&segment.Name,
			&segment.IsDeleted,
		)
		if err != nil {
			return entity.Segment{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
				OriginError:     err,
				OriginErrorText: err.Error(),
				Comment:         "Failed to scan segment to structure",
				Location:        "SegmentRepository.GetSegmentByName - rows.Scan",
			}}
		}
		return segment, nil
	}

	return entity.Segment{}, customError.ErrUserNotFound{ErrBase: customError.ErrBase{
		Comment:  fmt.Sprintf("Segment with name %s does not exist", name),
		Location: "SegmentRepository.GetSegmentByName",
	}}
}

// DeleteSegment используется для удаления сегмента в базе данных с указанным именем.
// DeleteSegment осуществляет логическое удаление, изменяя значение `is_deleted` на `true`.
// DeleteSegment не проверяет существование сегмента перед выполнением операции (проверка
// реализуется на уровне сервиса).
func (r *SegmentRepository) DeleteSegment(ctx context.Context, name string) error {
	// Получим id сегмента
	sql, args, err := r.Builder.
		Select("segment_id").
		From("segments").
		Where("name = ?", name).
		ToSql()
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to build sql query to get segment id by its name %s", name),
			Location:        "SegmentRepository.DeleteSegment - r.Builder",
		}}
	}

	var id int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to query segment's id by its name %s", name),
			Location:        "SegmentRepository.DeleteSegment - r.Pool.QueryRow",
		}}
	}

	// Пометим сегмент как удалённый
	sql, args, err = r.Builder.
		Update("segments").
		Set("is_deleted", true).
		Where("segment_id = ?", id).
		ToSql()
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to build query for segment deletion by its id %d", id),
			Location:        "SegmentRepository.DeleteSegment - r.Builder",
		}}
	}

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to perform delete segment query, inspect origin error text",
			Location:        "SegmentRepository.DeleteSegment - r.Pool.QueryRow",
		}}
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
			OriginError: err,
			Comment:     fmt.Sprintf("cannot delete segment with name \"%s\" because it does not exist", name),
			Location:    "SegmentRepository.DeleteSegment - r.Pool.QueryRow",
		}}
	}

	// Установим для всех пользователей, входивших в удаляемый сегмент,
	// время выхода из сегмента, равное моменту удаления сегмента (текущему времени)
	sql, args, err = r.Builder.
		Update("users_segments").
		Set("end_date", time.Now()).
		Where("segment_id = ?", id).
		ToSql()
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to build query for updating users with deleted segment",
			Location:        "SegmentRepository.DeleteSegment - r.Builder",
		}}
	}

	res, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to update users with deleted segment",
			Location:        "SegmentRepository.DeleteSegment - r.Pool.Exec",
		}}
	}

	return nil
}

// RecoverSegment используется при вызове операции создания сегмента в том случае,
// если сегмент с указанным именем существовал ранее и был удалён.
// RecoverSegment меняет флаг `is_deleted` у сегмента на `false`.
func (r *SegmentRepository) RecoverSegment(ctx context.Context, name string) (string, error) {
	sql, args, _ := r.Builder.
		Update("segments").
		Set("is_deleted", false).
		Where("name = ?", name).
		ToSql()

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return "", customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to recover deleted segment \"%s\"", name),
			Location:        "SegmentRepository.RecoverSegment - r.Exec",
		}}
	}

	if res.RowsAffected() == 0 {
		return "", customError.ErrSegmentNotFound{ErrBase: customError.ErrBase{
			Comment:  fmt.Sprintf("Failed to recover deleted segment \"%s\" because it does not exist at all", name),
			Location: "SegmentRepository.RecoverSegment - r.Exec",
		}}
	}

	return name, nil
}

// AddUsersToSegmentByRandomPercent используется для добавления случайных `percent`% пользователей в
// указанный сегмент.
func (r *SegmentRepository) AddUsersToSegmentByRandomPercent(ctx context.Context, name string, percent int) error {
	// Получим percent% случайных пользователей из БД
	sql, args, err := r.Builder.
		Select("user_id").
		From("users").
		Suffix(fmt.Sprintf("tablesample bernoulli (%d)", percent)).
		ToSql()
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to build sql query for random fetching % of users",
			Location:        "SegmentRepository.AddUsersToSegmentByRandomPercent - r.Builder",
		}}
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to perform sql query for random fetching % of users",
			Location:        "SegmentRepository.AddUsersToSegmentByRandomPercent - r.Pool.Query",
		}}
	}

	var userIDs []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
				OriginError:     err,
				OriginErrorText: err.Error(),
				Comment:         "Failed to scan user to structure",
				Location:        "SegmentRepository.AddUsersToSegmentByRandomPercent - rows.Scan",
			}}
		}
		userIDs = append(userIDs, id)
	}

	// Получим ID сегмента
	segment, err := r.GetSegmentByName(ctx, name)
	if err != nil {
		return err
	}

	// Сформируем запрос на добавление выбранных пользователей в сегмент
	insertQuery := r.Builder.
		Insert("users_segments").
		Columns("user_id", "segment_id")

	for _, userID := range userIDs {
		insertQuery = insertQuery.Values(userID, segment.ID)
	}

	sql, args, err = insertQuery.ToSql()
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to build sql query for adding users to segment",
			Location:        "SegmentRepository.AddUsersToSegmentByRandomPercent - insertQuery.ToSql()",
		}}
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to perform sql query for adding users to segment",
			Location:        "SegmentRepository.AddUsersToSegmentByRandomPercent - r.Pool.Exec",
		}}
	}

	return nil
}
