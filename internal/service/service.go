package service

import (
	"avito-rest-api/internal/entity"
	"avito-rest-api/internal/repository"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, input UserCreateInput) (int, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetAllUsersWithSegments(ctx context.Context) ([]entity.UserWithSegments, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserSegmentsByUserID(ctx context.Context, id int) ([]entity.UserSegmentInformation, error)
	GetUserWithSegmentsByUserID(ctx context.Context, id int) (entity.UserWithSegments, error)
	AddUserToSegments(ctx context.Context, id int, segments []entity.UserSegmentInformation) error
	DeleteUserFromSegments(ctx context.Context, id int, segments []string) error
}

type Segment interface {
	CreateSegment(ctx context.Context, input SegmentCreateInput) (string, error)
	GetAllSegments(ctx context.Context, sType int) ([]entity.Segment, error)
	GetSegmentByName(ctx context.Context, name string) (entity.Segment, error)
	DeleteSegment(ctx context.Context, name string) error
}

type Report interface {
	MakeReport(ctx context.Context) (entity.ReportCSV, error)
}

type Services struct {
	User    User
	Segment Segment
	Report  Report
}

type ServicesDependencies struct {
	Repositories *repository.Repositories
}

func NewService(dependencies ServicesDependencies) *Services {
	return &Services{
		User:    NewUserService(dependencies.Repositories.User, dependencies.Repositories.Segment),
		Segment: NewSegmentService(dependencies.Repositories.Segment),
		Report:  NewReportService(dependencies.Repositories.Report),
	}
}
