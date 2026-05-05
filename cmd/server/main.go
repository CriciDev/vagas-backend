package main

import (
	"fmt"
	"net/http"

	"github.com/CriciumaDevJobs/backend/infra"
	"github.com/CriciumaDevJobs/backend/internal"
)

func main() {

	db := infra.InitDB()

	app := internal.StartAppContext(db)

	fmt.Println("Iniciando servidor")

	fmt.Println("Servidor iniciado")
	http.ListenAndServe(":8080", app.Router)
}
