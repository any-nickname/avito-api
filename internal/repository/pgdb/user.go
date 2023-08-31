package pgdb

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/package/error"
	"avito-rest-api/package/postgres"
	"context"
	sqlLibrary "database/sql"
	"fmt"
	"strings"
	"time"
)

type UserRepository struct {
	*postgres.PostgreDB
}

// NewUserRepository инициализирует репозиторий `user`, инкапсулирующий логику хранения
// сущности пользователя.
func NewUserRepository(pg *postgres.PostgreDB) *UserRepository {
	return &UserRepository{pg}
}

// CreateUser добавляет в базу данных нового пользователя с указанными в `user` данными
// и возвращает `id` добавленной записи.
func (r *UserRepository) CreateUser(ctx context.Context, user entity.User) (int, error) {
	sql, args, _ := r.Builder.
		Insert("users").
		Columns("name", "lastname", "sex", "age").
		Values(user.Name, user.Lastname, user.Sex, user.Age).
		Suffix("RETURNING user_id").
		ToSql()

	var id int

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "failed to perform sql query",
			Location:        "UserRepository.CreateUser - r.Pool.QueryRow",
		}}
	}

	return id, nil
}

// GetAllUsers возвращает всех пользователей из базы данных.
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		ToSql()
	if err != nil {
		return nil, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to build query to fetch all users, inspect origin error text",
			Location:        "UserRepository.GetAllUsers - r.Builder",
		}}
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to fetch all users from database, inspect origin error text",
			Location:        "UserRepository.GetAllUsers - r.Pool.Query",
		}}
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Lastname,
			&user.Sex,
			&user.SexText,
			&user.Age,
			&user.IsDeleted,
		)
		if err != nil {
			return nil, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
				OriginError:     err,
				OriginErrorText: err.Error(),
				Comment:         "Failed to scan database's result to structure",
				Location:        "UserRepository.GetAllUsers - rows.Scan",
			}}
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByID возвращает информацию о пользователе с указанным `id`.
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("users").
		Where("user_id = ?", id).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.User{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to perform sql query to get user by id = %d", id),
			Location:        "UserRepository.GetUserByID - r.Pool.Exec",
		}}
	}

	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Lastname,
			&user.Sex,
			&user.SexText,
			&user.Age,
			&user.IsDeleted,
		)
		if err != nil {
			return entity.User{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
				OriginError:     err,
				OriginErrorText: err.Error(),
				Comment:         "Failed to scan user to structure",
				Location:        "UserRepository.GetUserByID - rows.Scan",
			}}
		}
		return user, nil
	}

	return entity.User{}, customError.ErrUserNotFound{ErrBase: customError.ErrBase{
		OriginError:     nil,
		OriginErrorText: "",
		Comment:         fmt.Sprintf("User with id %d not found", id),
		Location:        "UserRepository.GetUserByID",
	}}
}

// GetUserSegmentsByUserID возвращает список сегментов, в которые входит пользователь с указанным `id`,
// в формате структуры, содержащей: наименование сегмента, дату добавления пользователя в сегмент,
// дату выхода пользователя из сегмента (если установлена). Возвращает только актуальные на текущий
// момент сегменты.
func (r *UserRepository) GetUserSegmentsByUserID(ctx context.Context, id int) ([]entity.UserSegmentInformation, error) {
	sql, args, _ := r.Builder.
		Select("users_segments.user_segment_id", "users.user_id", "segments.segment_id", "segments.name", "to_char(users_segments.start_date, 'HH24:MI:SS DD.MM.YYYY') as start_date",
			"coalesce(to_char(users_segments.end_date, 'HH24:MI:SS DD.MM.YYYY'), '') as end_date").
		From("segments").
		Join("users_segments on users_segments.segment_id = segments.segment_id").
		Join("users on users.user_id = users_segments.user_id").
		Where("users.user_id = ? and (users_segments.end_date >= current_timestamp or users_segments.end_date is null)", id).
		ToSql()

	var userSegments []entity.UserSegmentInformation

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to perform query to get user segments by id",
			Location:        "UserRepository.GetUserSegmentsByUserID - r.Pool.Query",
		}}
	}

	for rows.Next() {
		var segmentInfo entity.UserSegmentInformation
		err := rows.Scan(
			&segmentInfo.InfoID,
			&segmentInfo.UserID,
			&segmentInfo.SegmentID,
			&segmentInfo.Name,
			&segmentInfo.StartDate,
			&segmentInfo.EndDate,
		)
		if err != nil {
			return nil, err
		}
		userSegments = append(userSegments, segmentInfo)
	}

	return userSegments, nil
}

// AddUserToSegments добавляет пользователя с идентификатором `id` в сегменты,
// перечисленные в списке `segments`, не проверяя пользователя на существование,
// и не проверяя сегменты на существование и метку `is_deleted`.
func (r *UserRepository) AddUserToSegments(ctx context.Context, id int, segments []entity.UserSegmentInformation) error {
	query := r.Builder.Insert("users_segments").Columns("user_id", "segment_id", "end_date")
	for _, segment := range segments {
		if segment.EndDate == "" {
			query = query.Values(id, segment.SegmentID, sqlLibrary.NullString{Valid: false})
		} else {
			databaseTimeFormat, _ := time.Parse("15:04:05 02.01.2006", segment.EndDate)
			query = query.Values(id, segment.SegmentID, databaseTimeFormat)
		}
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to build query to add user (id = %d) to segments", id),
			Location:        "UserRepository.AddUserToSegments - r.Builder",
		}}
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to perform query to add user (id = %d) to segments", id),
			Location:        "UserRepository.AddUserToSegments - r.Pool.Query",
		}}
	}

	return nil
}

// DeleteUserFromSegments удаляет пользователя из указанных сегментов, не осуществляя проверки
// на существование пользователя, существование сегментов.
func (r *UserRepository) DeleteUserFromSegments(ctx context.Context, id int, segments []entity.UserSegmentInformation) error {
	var infoIDs []interface{}
	var placeholders []string
	for _, s := range segments {
		infoIDs = append(infoIDs, s.InfoID)
		placeholders = append(placeholders, "?")
	}

	query := r.Builder.Update("users_segments").
		Set("end_date", time.Now()).
		Where(fmt.Sprintf("user_segment_id in (%s)", strings.Join(placeholders, ", ")), infoIDs...)
	sql, args, err := query.ToSql()
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to build sql query for deleting user (id = %d) from segments", id),
			Location:        "UserRepository.DeleteUserFromSegments - query.ToSql",
		}}
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         fmt.Sprintf("Failed to perform sql query for deleting user (id = %d) from segments", id),
			Location:        "UserRepository.DeleteUserFromSegments - r.Pool.Exec",
		}}
	}

	return nil
}
