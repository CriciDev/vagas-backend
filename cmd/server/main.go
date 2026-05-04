package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CriciumaDevJobs/backend/infra"
	"github.com/CriciumaDevJobs/backend/internal/devs"
)

func main() {

	db := infra.InitDB()

	devController := devs.InitializeDevController(db)

	fmt.Println("Iniciando servidor")

	http.HandleFunc("/devs", func(w http.ResponseWriter, r *http.Request) {
		devController.CreateDev(context.Background(), w, r)
	})

	fmt.Println("Servidor iniciado")
	http.ListenAndServe(":8080", nil)
}
