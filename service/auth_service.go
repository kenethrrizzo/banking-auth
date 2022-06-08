package service

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kenethrrizzo/banking-auth/dto"

	"github.com/kenethrrizzo/banking-auth/domain"
	repo "github.com/kenethrrizzo/banking-auth/domain/repositories"
	log "github.com/sirupsen/logrus"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, error)
	Verify(map[string]string) (bool, error)
}

type DefaultAuthService struct {
	repo            repo.AuthRepository
	rolePermissions domain.RolePermissions
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
		log.Error("Error while converting JWT from string")
		return false, err
	}
	// Checking the validity of the token, expiry time and signature
	if !jwtToken.Valid {
		log.Error("JWT is not valid")
		return false, errors.New("Invalid token")
	}
	// Casting the token claims to jwt.MapClaims
	mapClaims := jwtToken.Claims.(*domain.AccessTokenClaim)
	// If role is user, check if the account_id or customer_id belongs to the same token
	if mapClaims.IsUserRole() {
		if !mapClaims.IsRequestVerifiedWithTokenClaims(urlParams) {
			return false, nil
		}
	}
	// Verify if the user has permissions to the routes
	isAuthorized := s.rolePermissions.IsAuthorizedFor(mapClaims.Role, urlParams["route-name"])
	return isAuthorized, nil
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	// Secret key
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaim{}, keyFunc)
	if err != nil {
		log.Error("Error while parsing token: ", err.Error())
		return nil, err
	}
	return token, nil
}

func NewAuthService(repo repo.AuthRepository, permissions domain.RolePermissions) AuthService {
	return DefaultAuthService{
		repo,
		permissions,
	}
}
