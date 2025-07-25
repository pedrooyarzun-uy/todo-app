package service

import (
	"errors"
	"os"
	"time"
	"todo-app/internal/domain"
	"todo-app/internal/dto"
	"todo-app/internal/helpers"
	"todo-app/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(usr dto.CreateUserRequest) error
	VerifyUser(token string) error
	SignIn(usr dto.SignInRequest) (string, error)
}

type userService struct {
	ur repository.UserRepository
	vr repository.VerificationRepository
}

func NewUserService(
	ur repository.UserRepository,
	vr repository.VerificationRepository,
) UserService {
	return &userService{
		ur, vr,
	}
}

func (s *userService) CreateUser(usr dto.CreateUserRequest) error {
	res, err := s.ur.GetUserByEmail(usr.Email)

	if err != nil {
		return err
	}

	if res != nil {
		err := errors.New("el usuario ya existe en el sistema")
		return err
	}

	password, _ := helpers.HashPassword(usr.Password)

	user := domain.User{
		Name:      usr.Name,
		Email:     usr.Email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	err = s.ur.CreateUser(user)

	res, err = s.ur.GetUserByEmail(usr.Email)

	token := uuid.New().String()
	ev := domain.EmailVerification{
		Token:     token,
		UserId:    res.Id,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(5)),
		Used:      false,
	}

	err = s.vr.Save(&ev)

	link := "<a href='" + os.Getenv("ALLOWED_ORIGINS") + "/verify-account?token=" + token + "'>" + "Click aqui" + "</a>"
	err = helpers.SendMail(user.Email, "Verifica tu cuenta!", "Necesitamos que verifiques tu cuenta! Ingresa a: "+link)

	return err
}

func (s *userService) VerifyUser(token string) error {
	return s.vr.MarkAsUsed(token)
}

func (s *userService) SignIn(req dto.SignInRequest) (string, error) {
	usr, err := s.ur.GetUserByEmail(req.Email)

	if err != nil || usr == nil || !helpers.VerifyPassword(req.Password, usr.Password) {
		return "", errors.New("Incorrect email or password")
	}

	token, err := helpers.GenerateJWT(usr.Id)

	if err != nil {
		return "", err
	}

	return token, nil
}
