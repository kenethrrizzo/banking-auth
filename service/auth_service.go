package service

import (
	"github.com/kenethrrizzo/banking-auth/domain"
	"github.com/kenethrrizzo/banking-auth/dto"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, error)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	login, err := s.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}
	loginResponse := dto.LoginResponse{
		AccessToken: *token,
		RefreshToken: "",
	}
	return &loginResponse, nil
}

func NewAuthService(repo domain.AuthRepository) AuthService {
	return DefaultAuthService{repo}
}
