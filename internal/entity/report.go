package entity

type ReportCSV struct {
	ReportDate string `json:"report_date"`
	Report     string `json:"report"`
}

type Report struct {
	ReportDate string      `json:"report_date"`
	ReportRows []ReportRow `json:"report_rows"`
}

type ReportRow struct {
	UserID      int    `json:"user_id" csv:"user_id"`
	SegmentName string `json:"segment_name" csv:"segment_name"`
	StartDate   string `json:"start_date" csv:"start_date"` // Дата добавления пользователя в сегмент SegmentName
	EndDate     string `json:"end_date" csv:"end_date"`     // Дата выхода пользователя из сегмента, может быть пустой
}
