package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/helper"
	"github.com/mirzahilmi/hackathon/internal/pkg/pool"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
	"github.com/mirzahilmi/hackathon/internal/repository"
	"github.com/sirupsen/logrus"
	storage_go "github.com/supabase-community/storage-go"
)

type UserUsecaseItf interface {
	GetUserProfile(ctx context.Context) (model.UserCleanResource, error)
	ResetPassword(ctx context.Context, current, new string) error
	UpdateUserProfile(ctx context.Context, user *model.UserCleanResource) error
	ChangePicture(ctx context.Context, picture *multipart.FileHeader) error
	DeletePicture(ctx context.Context) error
}

type userUsecase struct {
	userRepo    repository.UserRepositoryItf
	supabaseImg *storage_go.Client
	log         *logrus.Logger
}

func NewUserUsecase(
	userRepo repository.UserRepositoryItf,
	supabaseImg *storage_go.Client,
	log *logrus.Logger,
) UserUsecaseItf {
	return &userUsecase{userRepo, supabaseImg, log}
}

func (u *userUsecase) GetUserProfile(ctx context.Context) (model.UserCleanResource, error) {
	client, err := u.userRepo.NewClient(false, nil)
	if err != nil {
		return model.UserCleanResource{}, err
	}
	user, err := client.GetByParam(ctx, "ID", ctx.Value(ClientID).(int64))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.UserCleanResource{}, response.NewError(fiber.StatusNotFound, ErrUserNotExist)
		default:
			return model.UserCleanResource{}, err
		}
	}

	userClean, err := user.Clean(u.supabaseImg)
	if err != nil {
		return model.UserCleanResource{}, err
	}
	return userClean, nil
}

func (u *userUsecase) ResetPassword(ctx context.Context, current, new string) error {
	id := ctx.Value(ClientID).(int64)
	client, err := u.userRepo.NewClient(false, nil)
	if err != nil {
		return err
	}
	user, err := client.GetByParam(ctx, "ID", id)
	if err != nil {
		return response.NewError(fiber.StatusNotFound, ErrUserNotExist)
	}
	if err := helper.BcryptCompare(user.Password, current); err != nil {
		return response.NewError(fiber.StatusUnauthorized, ErrWrongPassword)
	}
	hashed, err := helper.BcryptHash(new)
	if err != nil {
		return err
	}
	if err := client.UpdatePassword(ctx, int64(id), hashed); err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) UpdateUserProfile(ctx context.Context, userClean *model.UserCleanResource) error {
	user := model.UserResource{
		ID:   ctx.Value(ClientID).(int64),
		Name: userClean.Name,
	}
	client, err := u.userRepo.NewClient(false, nil)
	if err != nil {
		return err
	}
	mysqlErr := pool.MySQLErr.Get().(*mysql.MySQLError)
	defer pool.MySQLErr.Put(mysqlErr)
	if err := client.Update(ctx, &user); err != nil {
		switch {
		case errors.As(err, &mysqlErr) && mysqlErr.Number == 1062:
			return response.NewCustomError(fiber.StatusConflict, fiber.Map{
				"name": ErrNameExist.Error(),
			})
		default:
			return err
		}
	}
	return nil
}

func (u *userUsecase) ChangePicture(ctx context.Context, picture *multipart.FileHeader) error {
	pictureFile, err := picture.Open()
	defer func() {
		if err := pictureFile.Close(); err != nil {
			u.log.Warn(err)
		}
	}()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%d-%s-%s", time.Now().UnixMilli(), ctx.Value(ClientName), picture.Filename)
	if _, err := u.supabaseImg.UploadFile(os.Getenv("SUPABASE_BUCKET_ID"), filename, pictureFile); err != nil {
		return err
	}

	client, err := u.userRepo.NewClient(false, nil)
	if err != nil {
		return err
	}
	if err := client.UpdatePicture(ctx, ctx.Value(ClientID).(int64), filename); err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) DeletePicture(ctx context.Context) error {
	client, err := u.userRepo.NewClient(false, nil)
	if err != nil {
		return err
	}

	// TODO: Also delete image in supabase, client lib aren't viable for now
	// see https://github.com/supabase-community/storage-go/issues/28

	if err := client.UpdatePicture(ctx, ctx.Value(ClientID).(int64), ""); err != nil {
		return err
	}
	return nil
}
