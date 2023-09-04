package pgdb

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/package/postgres"
	"context"
	"time"
)

type ReportRepository struct {
	*postgres.PostgreDB
}

func NewReportRepository(pg *postgres.PostgreDB) *ReportRepository {
	return &ReportRepository{pg}
}

// MakeReport формирует новый отчёт об истории вхождения-выхождения пользователей из сегментов,
// содержащий столбцы `user_id`, `segment_name`, `start_date`, `end_date`
func (r *ReportRepository) MakeReport(ctx context.Context) (entity.Report, error) {
	sql, args, err := r.Builder.
		Select("u.user_id, s.name as segment_name, to_char(us.start_date, 'HH24:MI:SS DD.MM.YYYY'), coalesce(to_char(us.end_date, 'HH24:MI:SS DD.MM.YYYY'), '')").
		From("users u").
		Join("users_segments us on us.user_id = u.user_id").
		Join("segments s on s.segment_id = us.segment_id").
		OrderBy("us.start_date asc").
		ToSql()
	if err != nil {
		return entity.Report{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to build sql expression for fetching report data, inspect origin error text",
			Location:        "ReportRepository.MakeReport: r.Builder",
		}}
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Report{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to execute sql expression for fetching report data, inspect origin error text",
			Location:        "ReportRepository.MakeReport: r.Pool.Query",
		}}
	}

	report := entity.Report{ReportDate: time.Now().Format("15:04:05 02.01.2006")}
	for rows.Next() {
		var reportRow entity.ReportRow
		err = rows.Scan(
			&reportRow.UserID,
			&reportRow.SegmentName,
			&reportRow.StartDate,
			&reportRow.EndDate,
		)
		if err != nil {
			return entity.Report{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
				OriginError:     err,
				OriginErrorText: err.Error(),
				Comment:         "Failed to scan report row to structure, inspect origin error text",
				Location:        "ReportRepository.MakeReport: rows.Scan",
			}}
		}
		report.ReportRows = append(report.ReportRows, reportRow)
	}

	return report, nil
}
