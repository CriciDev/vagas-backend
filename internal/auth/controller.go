package auth

import (
	"encoding/json"
	"net/http"

	"github.com/CriciumaDevJobs/backend/handlers"
)

type AuthenticationController struct {
	AuthUseCase *AuthenticationUseCase
}

func NewAuthenticationController(usecase *AuthenticationUseCase) *AuthenticationController {
	auth := AuthenticationController{
		AuthUseCase: usecase,
	}

	return &auth
}

func (controller *AuthenticationController) AuthenticateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req = AuthenticationRequest{}

	json_err := json.NewDecoder(r.Body).Decode(&req)

	if json_err != nil {
		handlers.ResponseWithHttpError(w, http.StatusBadRequest, "JSON Enviado não segue o padrão esperado")
		return
	}

	resp, http_err := controller.AuthUseCase.AuthenticateUser(r.Context(), req.Email, req.Password)

	if http_err != nil {
		handlers.ResponseWithHttpError(w, http_err.Code, http_err.Message)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
