package my_sql_db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DbName = "NikolayInternDB"
	DbUser = "root"
	DbPass = "password"
	DbHost = "localhost"
	DbPort = "3306"
)

func NewDbConnection() *sql.DB {
	dsnString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", DbUser, DbPass, DbHost, DbPort, DbName)

	db, err := sql.Open("mysql", dsnString)
	if err != nil {
		panic(err)
	}

	return db
}
