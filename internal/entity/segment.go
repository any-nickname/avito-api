package entity

type Segment struct {
	ID        int    `json:"segment_id" example:"43"`
	Name      string `json:"name" example:"AVITO_MUSIC_SERVICE"`
	IsDeleted bool   `json:"is_deleted" example:"false"`
}
