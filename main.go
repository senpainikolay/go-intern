package main

import (
	controller "github.com/senpainikolay/go-tasks/http-controller"
	"github.com/senpainikolay/go-tasks/models"
	my_sql_db "github.com/senpainikolay/go-tasks/pkg"
	"github.com/senpainikolay/go-tasks/repository"
)

func main() {
	databaseConfig := models.DatabaseConfig{
		DbName: "NikolayInternDB",
		DbUser: "root",
		DbPass: "password",
		DbHost: "localhost",
		DbPort: "3306",
	}

	db := my_sql_db.NewDbConnection(databaseConfig)
	defer db.Close()

	generalRepository := repository.NewGeneralRepository(db)
	err := generalRepository.TryCreate()
	if err != nil {
		panic(err)
	}
	err = generalRepository.PopulateRandomDB()
	if err != nil {
		panic(err)
	}

	generalController := controller.NewController(generalRepository)

	controller.Serve(generalController, "8080")

}
