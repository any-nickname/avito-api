package service

import (
	"avito-rest-api/internal/entity"
	customError "avito-rest-api/internal/error"
	"avito-rest-api/internal/repository"
	"avito-rest-api/internal/webapi"
	"context"
	"fmt"
	"github.com/gocarina/gocsv"
	"time"
)

type ReportService struct {
	reportRepository repository.Report
	gDrive           webapi.GDrive
}

func NewReportService(reportRepository repository.Report, gDrive webapi.GDrive) *ReportService {
	return &ReportService{
		reportRepository: reportRepository,
		gDrive:           gDrive,
	}
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

	var reportContent string
	if rs.gDrive.IsSet() {
		// Загрузка отчёта на google drive в случае, если установлен путь
		// до credentials в файле конфигураций, и возврат ссылки
		now := time.Now()
		reportContent, err = rs.gDrive.UploadCSVFile(
			ctx,
			fmt.Sprintf(
				"report_%d-%d_%d.%d.%d.csv",
				now.Hour(),
				now.Minute(),
				now.Day(),
				now.Month(),
				now.Year(),
			),
			[]byte(reportCSV),
		)
		if err != nil {
			return entity.ReportCSV{}, customError.ErrInternalServerError{ErrBase: customError.ErrBase{
				OriginError:     err,
				OriginErrorText: err.Error(),
				Comment:         "Failed to upload report.csv to the google drive, please inspect origin error text",
				Location:        "ReportService.MakeReport - gDrive.UploadCSVFile",
			}}
		}
	} else { // Возврат отчёта прямо в теле запроса
		reportContent = reportCSV
	}

	return entity.ReportCSV{
		ReportDate: report.ReportDate,
		Report:     reportContent,
	}, nil
}
