package entity

type User struct {
	ID        int    `json:"user_id" example:"16"`
	Name      string `json:"name" example:"Михаил"`
	Lastname  string `json:"lastname" example:"Иванов"`
	Sex       int    `json:"sex" example:"0" enums:"0,1"` // Пол, 0 - мужской, 1 - женский
	SexText   string `json:"sex_text" example:"мужской" enums:"мужской,женский"`
	Age       int    `json:"age" example:"26"`
	IsDeleted bool   `json:"is_deleted" example:"false"`
}

type UserSegmentInformation struct {
	InfoID    int    `json:"information_id" example:"179"`             // ID записи в ассоциативной таблице, связывающей пользователей с сегментами
	UserID    int    `json:"user_id" example:"16"`                     // ID пользователя, который входит в сегмент
	SegmentID int    `json:"segment_id" example:"43"`                  // ID сегмента
	Name      string `json:"name" example:"AVITO_MUSIC_SERVICE"`       // Наименование сегмента
	StartDate string `json:"start_date" example:"15:27:32 01.09.2023"` // Дата добавления пользователя в сегмент
	EndDate   string `json:"end_date" example:""`                      // Дата выхода пользователя из сегмента
}

type UserWithSegments struct {
	User     User                     `json:"user"`
	Segments []UserSegmentInformation `json:"segments"`
}
