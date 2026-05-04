package devs

import (
	"context"
	"encoding/json"
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

func (controller *DevController) CreateDev(ctx context.Context, writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		handlers.ResponseWithHttpError(writer, http.StatusMethodNotAllowed, "Method not Allowed")
		return
	}

	var dev = Dev{}

	err := json.NewDecoder(request.Body).Decode(&dev)

	if err != nil {
		handlers.ResponseWithHttpError(writer, http.StatusBadRequest, "JSON Enviado não segue a estrutura esperada!")
		return
	}

	resp, err := controller.Usecase.CreateDev(ctx, &dev)

	if err != nil {
		handlers.ResponseWithHttpError(writer, http.StatusUnprocessableEntity, err.Error())
		return
	}

	json.NewEncoder(writer).Encode(resp)
}
