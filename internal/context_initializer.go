package internal

import (
	"database/sql"
	"net/http"

	"github.com/CriciumaDevJobs/backend/internal/auth"
	"github.com/CriciumaDevJobs/backend/internal/devs"
)

type AppContext struct {
	Router *http.ServeMux
	DB     *sql.DB
}

func StartAppContext(db *sql.DB) *AppContext {
	app := &AppContext{
		Router: http.NewServeMux(),
		DB:     db,
	}

	app.setupContext()

	return app
}

func (app *AppContext) setupContext() {

	devRepo := devs.New(app.DB)

	devUseCase := devs.NewDevUseCase(devRepo)
	authUseCase := auth.NewAuthenticationUseCase(devUseCase)

	devController := devs.NewDevController(devUseCase)
	authController := auth.NewAuthenticationController(authUseCase)

	app.registerDevRoutes(devController)
	app.registerAuthRoutes(authController)
}

func (app *AppContext) registerDevRoutes(controller *devs.DevController) {

	app.Router.HandleFunc("POST /devs/register", controller.CreateDev)

	app.Router.HandleFunc("GET /devs/me", auth.AuthenticationMiddleware(controller.FindDevProfile))
}

func (app *AppContext) registerAuthRoutes(controller *auth.AuthenticationController) {
	app.Router.HandleFunc("POST /auth/login", controller.AuthenticateUser)
}
