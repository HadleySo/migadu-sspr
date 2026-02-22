package routes

import (
	"golang.hadleyso.com/msspr/src/auth"
	"golang.hadleyso.com/msspr/src/handlers"
)

func userSettings() {
	authedRouter := Router.PathPrefix("/my").Subrouter()
	authedRouter.Use(auth.MiddleValidateSession)

	authedRouter.HandleFunc("/", handlers.GetInfo).Methods("GET")
	authedRouter.HandleFunc("/passwd", handlers.SetPasswd).Methods("POST")
}
