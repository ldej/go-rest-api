package user

import (
	"context"
	"errors"

	"github.com/ldej/go-rest-example/internal/api"
	"github.com/ldej/go-rest-example/pkg/encrypt"
	"github.com/ldej/go-rest-example/pkg/uid"
)

type Service struct {
	uidGenerator uid.Generator
	encryptor    encrypt.Encryptor
	userRepo     Repository
}

func NewService(userRepo Repository, uidGenerator uid.Generator, encryptor encrypt.Encryptor) *Service {
	return &Service{userRepo: userRepo, uidGenerator: uidGenerator, encryptor: encryptor}
}

func (s *Service) Register(c context.Context, registerUser api.RegisterUserRequest) (*api.User, error) {
	user := api.User{
		UID:               s.uidGenerator.NewUUID(),
		Name:              registerUser.Name,
		EmailAddress:      registerUser.EmailAddress,
		EncryptedPassword: s.encryptor.Encrypt(registerUser.Password),
	}
	return s.userRepo.UserCreate(c, &user)
}

func (s *Service) GetByUID(c context.Context, uid string) (*api.User, error) {
	return s.userRepo.UserGetByUID(c, uid)
}

func (s *Service) Login(c context.Context, request api.LoginRequest) (*api.User, error) {
	user, err := s.userRepo.UserGetByEmailAddress(c, request.EmailAddress)
	if err != nil {
		return nil, err
	}
	if !s.encryptor.IsValid(user.EncryptedPassword, request.Password) {
		return nil, errors.New("password invalid")
	}
	return user, nil
}
