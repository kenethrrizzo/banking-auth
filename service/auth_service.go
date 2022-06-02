package service

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kenethrrizzo/banking-auth/domain"
	"github.com/kenethrrizzo/banking-auth/dto"

	log "github.com/sirupsen/logrus"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, error)
	Verify(map[string]string) (bool, error)
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
		AccessToken:  *token,
		RefreshToken: "",
	}
	return &loginResponse, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) (bool, error) {
	jwtToken, err := jwtTokenFromString(urlParams["token"])
	if err != nil {
		return false, err
	}
	// Checking the validity of the token, expiry time and signature
	if !jwtToken.Valid {
		return false, errors.New("Invalid token")
	}
	// Casting the token claims to jwt.MapClaims
	//mapClaims := jwtToken.Claims.(jwt.MapClaims)

	//TODO
	return false, nil
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	// Secret key
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		log.Error("Error while parsing token: ", err.Error())
		return nil, err
	}
	return token, nil
}

func NewAuthService(repo domain.AuthRepository) AuthService {
	return DefaultAuthService{repo}
}
