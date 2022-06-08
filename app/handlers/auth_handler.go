package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kenethrrizzo/banking-auth/dto"
	"github.com/kenethrrizzo/banking-auth/service"
	log "github.com/sirupsen/logrus"
)

type AuthHandler struct {
	Service service.AuthService
}

func (h AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	log.Debug("[GET] login")
	var LoginRequest dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&LoginRequest)
	if err != nil {
		log.Error("Error while decoding login request: ", err.Error())
		writeResponse(rw, http.StatusBadRequest, dto.BadResponse{Error: err.Error()})
		return
	}
	auth_response, err := h.Service.Login(LoginRequest)
	if err != nil {
		log.Error("Error while authenticating user: ", err.Error())
		writeResponse(rw, http.StatusUnauthorized, dto.BadResponse{Error: err.Error()})
		return
	}
	writeResponse(rw, http.StatusOK, auth_response)
}

func (h AuthHandler) Verify(rw http.ResponseWriter, r *http.Request) {
	log.Debug("[GET] verify")
	urlParams := make(map[string]string)
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}
	if urlParams["token"] == "" {
		writeResponse(rw, http.StatusForbidden, dto.BadResponse{Error: errors.New("Missing tolen").Error()})
		return
	}
	isAuthorized, err := h.Service.Verify(urlParams)
	if err != nil {
		writeResponse(rw, http.StatusForbidden, dto.BadResponse{Error: errors.New("Missing tolen").Error()})
		return
	}
	if !isAuthorized {
		writeResponse(rw, http.StatusForbidden, dto.BadResponse{Error: errors.New("Missing tolen").Error()})
		return
	}
	writeResponse(rw, http.StatusOK, map[string]bool{"is_authorized": true})
}

func writeResponse(rw http.ResponseWriter, code int, data interface{}) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		log.Error(err.Error())
	}
} 
