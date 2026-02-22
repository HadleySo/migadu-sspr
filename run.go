package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"golang.hadleyso.com/msspr/src/config"
	"golang.hadleyso.com/msspr/src/models"
	"golang.hadleyso.com/msspr/src/routes"
)

func main() {
	viper.SetConfigName("msspr")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./data")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(&config.C); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Register struct
	gob.Register(&models.UserInfo{})

	// Register Routes
	routes.Main()

	// Listen
	port := viper.GetString("SERVER_PORT")
	log.Println("Listening to localhost:" + port)
	http.ListenAndServe(fmt.Sprintf("localhost:%v", port), routes.Router)
}
