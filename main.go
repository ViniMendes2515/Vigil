package main

import (
	"log"
	"os"
	"vigil/cmd"
	"vigil/database"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatal("Erro ao inicializar o banco de dados: ", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal("Erro ao executar CLI: ", err)
		os.Exit(1)
	}

	defer database.DB.Close()
}
