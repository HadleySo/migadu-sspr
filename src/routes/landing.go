package routes

import (
	"golang.hadleyso.com/msspr/src/handlers"
)

func landing() {
	Router.HandleFunc("/", handlers.Landing).Methods("GET")
}
