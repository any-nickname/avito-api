package error

import "fmt"

type ErrBase struct {
	OriginError     error `json:"-"`
	OriginErrorText string
	Title           string
	Comment         string
	Location        string
}

func (e ErrBase) Error() string {
	return fmt.Sprintf("%s(%s)", e.Comment, e.Location)
}

func (e ErrBase) Unwrap() error {
	return e.OriginError
}

// ErrInternalServerError обозначает
// непредвиденную внутреннюю ошибку сервера.
type ErrInternalServerError struct {
	ErrBase
}

// ErrUserValidationError обозначает
// ошибку валидации данных, передаваемых
// при совершении запроса к сущности `user`.
type ErrUserValidationError struct {
	ErrBase
}

// ErrUserNotFound обозначает ошибку,
// возникающую, когда происходит попытка
// обратиться к несуществующему пользователя.
type ErrUserNotFound struct {
	ErrBase
}

// ErrUserDeleted обозначает ошибку,
// возникающую, когда происходит попытка
// обратиться к удалённому пользователю.
type ErrUserDeleted struct {
	ErrBase
}

// ErrSegmentValidationError обозначает ошибку
// валидации данных сегмента.
type ErrSegmentValidationError struct {
	ErrBase
}

// ErrSegmentNotFound обозначает ошибку
// при обращении к несуществующему сегменту.
type ErrSegmentNotFound struct {
	ErrBase
}

// ErrSegmentDeleted обозначает ошибку
// при обращении к сегменту, который
// помечен как удалён.
type ErrSegmentDeleted struct {
	ErrBase
}

// ErrSegmentAlreadyExists используется, когда
// происходит попытка создать сегмент, который
// уже существует в базе данных.
type ErrSegmentAlreadyExists struct {
	ErrBase
}

// ErrReportValidationError используется,
// когда запрос на создание отчёта содержит
// некорректные данные.
type ErrReportValidationError struct {
	ErrBase
}
