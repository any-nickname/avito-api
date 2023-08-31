package repository

import (
	"avito-rest-api/internal/entity"
	"avito-rest-api/internal/repository/pgdb"
	"avito-rest-api/package/postgres"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserSegmentsByUserID(ctx context.Context, id int) ([]entity.UserSegmentInformation, error)
	AddUserToSegments(ctx context.Context, id int, segments []entity.UserSegmentInformation) error
	DeleteUserFromSegments(ctx context.Context, id int, segments []entity.UserSegmentInformation) error
}

type Segment interface {
	CreateSegment(ctx context.Context, segment entity.Segment) (string, error)
	GetAllSegments(ctx context.Context, sTypes int) ([]entity.Segment, error)
	GetSegmentByName(ctx context.Context, name string) (entity.Segment, error)
	DeleteSegment(ctx context.Context, name string) error
	RecoverSegment(ctx context.Context, name string) (string, error)
	AddUsersToSegmentByRandomPercent(ctx context.Context, name string, percent int) error
}

type Repositories struct {
	User
	Segment
}

func NewRepositories(pg *postgres.PostgreDB) *Repositories {
	return &Repositories{
		User:    pgdb.NewUserRepository(pg),
		Segment: pgdb.NewSegmentRepository(pg),
	}
}
