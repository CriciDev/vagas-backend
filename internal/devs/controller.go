package devs

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CriciumaDevJobs/backend/handlers"
)

type DevController struct {
	Usecase *DevUseCase
}

func NewDevController(usecase *DevUseCase) *DevController {
	controller := DevController{
		Usecase: usecase,
	}

	return &controller
}

func (controller *DevController) CreateDev(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")

	var dev = Dev{}

	json_err := json.NewDecoder(request.Body).Decode(&dev)

	if json_err != nil {
		handlers.ResponseWithHttpError(writer, http.StatusBadRequest, "JSON Enviado não segue a estrutura esperada!")
		return
	}

	resp, http_err := controller.Usecase.CreateDev(request.Context(), &dev)

	if http_err != nil {
		handlers.ResponseWithHttpError(writer, http_err.Code, http_err.Message)
		return
	}

	json.NewEncoder(writer).Encode(resp)
}

func (controller *DevController) FindDevProfile(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")

	devID, err := extractUserIDFromRequestToken(request)

	if err != nil {
		handlers.ResponseWithHttpError(writer, err.Code, err.Message)
		return
	}

	resp, err := controller.Usecase.FindDevByID(request.Context(), *devID)

	if err != nil {
		handlers.ResponseWithHttpError(writer, err.Code, err.Message)
		return
	}

	json.NewEncoder(writer).Encode(resp)
}

func extractUserIDFromRequestToken(request *http.Request) (*int32, *handlers.ErrorResponse) {

	value := request.Context().Value("user_id")

	val, ok := value.(float64)

	if !ok {
		log.Printf("ERRO: Falha ao extrair ID do usuário no token JWT!")
		return nil, handlers.NewError(http.StatusInternalServerError, "Erro Interno!")
	}

	var id = int32(val)

	return &id, nil

}
