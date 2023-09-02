package service

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/internal/repository"
	"context"
	"github.com/gocarina/gocsv"
)

type ReportService struct {
	reportRepository repository.Report
}

func NewReportService(reportRepository repository.Report) *ReportService {
	return &ReportService{reportRepository: reportRepository}
}

func (rs *ReportService) MakeReport(ctx context.Context) (entity.ReportCSV, error) {
	report, err := rs.reportRepository.MakeReport(ctx)
	if err != nil {
		return entity.ReportCSV{}, err
	}

	reportCSV, err := gocsv.MarshalString(&(report.ReportRows))
	if err != nil {
		return entity.ReportCSV{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
			OriginError:     err,
			OriginErrorText: err.Error(),
			Comment:         "Failed to marshal report to CSV, inspect origin error text",
			Location:        "ReportService.MakeReport - gocsv.MarshalString",
		}}
	}

	return entity.ReportCSV{
		ReportDate: report.ReportDate,
		Report:     reportCSV,
	}, nil
}
