package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/email"
	"github.com/mirzahilmi/hackathon/internal/pkg/helper"
	"github.com/mirzahilmi/hackathon/internal/pkg/pasetok"
	"github.com/mirzahilmi/hackathon/internal/pkg/pool"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
	"github.com/mirzahilmi/hackathon/internal/repository"
	"github.com/sirupsen/logrus"
)

type AuthUsecaseItf interface {
	RegisterUser(ctx context.Context, user model.UserSignupRequest) error
	LogUserIn(ctx context.Context, attempt model.UserLoginRequest) (string, error)
	VerifyUser(ctx context.Context, verifReq model.UserVerifRequest) error
}

type authUsecase struct {
	userRepo  repository.UserRepositoryItf
	verifRepo repository.VerifRepositoryItf
	mailer    email.VerificationMailer
	log       *logrus.Logger
}

func NewAuthUsecase(
	userRepo repository.UserRepositoryItf,
	verifRepo repository.VerifRepositoryItf,
	mailer email.VerificationMailer,
	log *logrus.Logger,
) AuthUsecaseItf {
	return &authUsecase{userRepo, verifRepo, mailer, log}
}

func (u *authUsecase) RegisterUser(ctx context.Context, userReq model.UserSignupRequest) error {
	hashed, err := helper.BcryptHash(userReq.Password)
	if err != nil {
		return err
	}
	user := model.UserResource{
		Email:    userReq.Email,
		Password: hashed,
		Name:     userReq.Name,
	}

	userClient, err := u.userRepo.NewClient(true, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := userClient.Rollback(); err != nil {
			u.log.Warn(err)
		}
	}()
	mysqlErr := pool.MySQLErr.Get().(*mysql.MySQLError)
	defer pool.MySQLErr.Put(mysqlErr)
	userID, err := userClient.Create(ctx, &user)
	if err != nil {
		switch {
		case errors.As(err, &mysqlErr) && mysqlErr.Number == 1062:
			if strings.Contains(strings.ToLower(mysqlErr.Message), "name") {
				return response.NewCustomError(fiber.StatusConflict, fiber.Map{
					"name": ErrNameExist.Error(),
				})
			}
			return response.NewCustomError(fiber.StatusConflict, fiber.Map{
				"email": ErrEmailExist.Error(),
			})
		default:
			return err
		}
	}

	verif := model.UserVerificationResource{
		UserID: userID,
		Token:  helper.RandString(32),
	}
	verifClient, err := u.verifRepo.NewClient(true, userClient.Ext())
	if err != nil {
		return err
	}
	verifID, err := verifClient.Create(ctx, &verif)
	if err != nil {
		return err
	}

	emailProps := map[string]string{
		"AppName": os.Getenv("APP_NAME"),
		"URL":     fmt.Sprintf("%s/verify-email?t=%d.%s", os.Getenv("WEB_DOMAIN"), verifID, verif.Token),
		"Expiry":  os.Getenv("MFA_EMAIL_TTL"),
	}
	if err := u.mailer.SendMail(user.Email, "Account Verification", email.VerificationView, emailProps); err != nil {
		return err
	}

	if err := userClient.Commit(); err != nil {
		return err
	}
	return nil
}

func (u *authUsecase) LogUserIn(ctx context.Context, attempt model.UserLoginRequest) (string, error) {
	client, err := u.userRepo.NewClient(false, nil)
	if err != nil {
		return "", err
	}
	user, err := client.GetByParam(ctx, "Email", attempt.Email)
	if err != nil {
		return "", response.NewError(fiber.StatusNotFound, ErrUserNotExist)
	}

	if !user.Active {
		return "", response.NewError(fiber.StatusUnauthorized, ErrUserNotActive)
	}

	if err := helper.BcryptCompare(user.Password, attempt.Password); err != nil {
		return "", response.NewError(fiber.StatusUnauthorized, ErrWrongPassword)
	}

	token := paseto.NewToken()
	token.SetAudience("*")
	token.SetIssuer(os.Getenv("APP_HOST"))
	token.SetSubject(fmt.Sprint(user.ID))
	token.SetString("name", user.Name)
	ttl, err := strconv.Atoi(os.Getenv("PASETO_TTL"))
	if err != nil {
		return "", err
	}
	token.SetExpiration(time.Now().Add(time.Duration(ttl) * time.Minute))
	token.SetNotBefore(time.Now())
	token.SetIssuedAt(time.Now())

	signed, err := pasetok.Encode(token)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (u *authUsecase) VerifyUser(ctx context.Context, verifReq model.UserVerifRequest) error {
	verifID, err := strconv.Atoi(verifReq.ID)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrIDNotNumeric)
	}
	verifClient, err := u.verifRepo.NewClient(true, nil)
	if err != nil {
		return err
	}
	verif, err := verifClient.GetByIDAndToken(ctx, int64(verifID), verifReq.Token)
	if err != nil {
		return response.NewError(fiber.StatusNotFound, ErrVerificationNotExist)
	}
	defer func() {
		if err := verifClient.Rollback(); err != nil {
			u.log.Warn(err)
		}
	}()
	if err := verifClient.UpdateSucceedStatus(ctx, verif.ID); err != nil {
		return err
	}
	userClient, err := u.userRepo.NewClient(true, verifClient.Ext())
	if err != nil {
		return err
	}
	if err := userClient.UpdateActiveStatus(ctx, verif.UserID); err != nil {
		return err
	}
	if err := verifClient.Commit(); err != nil {
		return err
	}
	return nil
}
