package services

import (
	"errors"

	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/repositories"
	"github.com/fahmiarz/project-management/utils"
	"github.com/google/uuid"
)

//user service ini bisa apa aja: contoh register dll
type UserService interface {
	Register (user *models.User) error
}

type userService struct{
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService)Register(user *models.User) error {
	//mengecek email sudah terdaftar atau belum
	existingUser , _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID !=0 {
		return errors.New("Email already registered")
	}
	//hasing password
	hased, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hased
	//set role
	user.Role = "user"
	user.PublicID = uuid.New()
	//simpan user
	return s.repo.Create(user)
}