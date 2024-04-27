package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/repository"
	storage_go "github.com/supabase-community/storage-go"
)

type ReportUsecaseItf interface {
	CreateReport(ctx context.Context, req model.ReportRequest) error
	GetReports(ctx context.Context) ([]model.ReportResource, error)
	UpdateReport(ctx context.Context, id int64, action string) error
}

type ReportUsecase struct {
	reportRepo repository.ReportRepositoryItf
	supabase   *storage_go.Client
}

func NewReportUsecase(reportRepo repository.ReportRepositoryItf, supabase *storage_go.Client) ReportUsecaseItf {
	return &ReportUsecase{reportRepo, supabase}
}

func (u *ReportUsecase) CreateReport(ctx context.Context, req model.ReportRequest) error {
	client, err := u.reportRepo.NewClient(true, nil)
	if err != nil {
		return err
	}
	defer client.Rollback()

	req.Picture.Filename = fmt.Sprintf("%d_%s", time.Now().UnixMilli(), req.Picture.Filename)
	pict, err := req.Picture.Open()
	if err != nil {
		return err
	}

	_, err = u.supabase.UploadFile(os.Getenv("SUPABASE_BUCKET_ID"), req.Picture.Filename, pict)
	if err != nil {
		return err
	}

	pictUrl := u.supabase.GetPublicUrl(os.Getenv("SUPABASE_BUCKET_ID"), req.Picture.Filename)
	report := model.ReportResource{
		Picture:     pictUrl.SignedURL,
		Description: req.Description,
		Location:    req.Location,
	}
	log.Println(report)
	if err := client.CreateReport(ctx, &report); err != nil {
		return err
	}

	return client.Commit()
}

func (u *ReportUsecase) GetReports(ctx context.Context) ([]model.ReportResource, error) {
	client, err := u.reportRepo.NewClient(false, nil)
	if err != nil {
		return nil, err
	}

	reports, err := client.GetReports(ctx)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (u *ReportUsecase) UpdateReport(ctx context.Context, id int64, action string) error {
	client, err := u.reportRepo.NewClient(true, nil)
	if err != nil {
		return err
	}

	report, err := client.GetReportByID(ctx, id)
	if err != nil {
		return err
	}

	report.Action = action
	if err := client.UpdateActionReport(ctx, &report); err != nil {
		return err
	}

	return client.Commit()
}
