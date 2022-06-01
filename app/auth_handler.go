package app

import (
	"encoding/json"
	"net/http"

	"github.com/kenethrrizzo/banking-auth/dto"
	"github.com/kenethrrizzo/banking-auth/service"
	log "github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service service.AuthService
}

func (h AuthHandler) login(rw http.ResponseWriter, r *http.Request) {
	var LoginRequest dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&LoginRequest)
	if err != nil {
		log.Error("Error while decoding login request: ", err.Error())
		writeResponse(rw, http.StatusBadRequest, dto.BadResponse{Error: err.Error()})
		return
	}
	auth_response, err := h.service.Login(LoginRequest)
	if err != nil {
		log.Error("Error while authenticating user: ", err.Error())
		writeResponse(rw, http.StatusUnauthorized, dto.BadResponse{Error: err.Error()})
		return
	}
	writeResponse(rw, http.StatusOK, auth_response)
}

func writeResponse(rw http.ResponseWriter, code int, data interface{}) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		log.Error(err.Error())
	}
}
