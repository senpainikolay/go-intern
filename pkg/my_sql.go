package my_sql_db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/senpainikolay/go-tasks/models"
)

func NewDbConnection(config models.DatabaseConfig) *sql.DB {
	dsnString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)

	db, err := sql.Open("mysql", dsnString)
	if err != nil {
		panic(err)
	}

	return db
}
