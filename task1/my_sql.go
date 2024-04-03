package task1

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

func NewDbConnection() (*sql.DB, func()) {
	dsnString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", DbUser, DbPass, DbHost, DbPort, DbName)

	db, err := sql.Open("mysql", dsnString)
	checkError(err)

	return db, func() {
		_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
		checkError(err)
		_, err = db.Exec("truncate table sources;")
		checkError(err)
		_, err = db.Exec("truncate table campaigns;")
		checkError(err)
		_, err = db.Exec("truncate table sources_campaigns;")
		checkError(err)
		_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
		checkError(err)
	}
}

func InitDB() {
	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/", DbUser, DbPass, DbHost, DbPort)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)

	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + DbName)
	checkError(err)

	_, err = db.Exec("USE " + DbName)
	checkError(err)

	_, err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS sources (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(20) UNIQUE NOT NULL);
		`)
	checkError(err)

	_, err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS campaigns (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(20) UNIQUE NOT NULL);
		`)
	checkError(err)

	_, err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS sources_campaigns (
		source_id INT NOT NULL, 
		campaign_id INT NOT NULL, 
		PRIMARY KEY (source_id, campaign_id),
		FOREIGN KEY (source_id) REFERENCES sources(id),
		FOREIGN KEY (campaign_id) REFERENCES campaigns(id) );
		`)
	checkError(err)

}
