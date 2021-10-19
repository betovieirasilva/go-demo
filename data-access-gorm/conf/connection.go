package conf

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func SetupDatabaseSqlDB() *sql.DB {
	// Create connection
	//From => https://pkg.go.dev/github.com/go-sql-driver/mysql#Config
	cfg := mysql.Config{
		// @see export properties, get with os.Getenv("DBUSER"), etc
		//$ export DBUSER=root
		//$ export DBPASS=1
		//User:   os.Getenv("DBUSER"),
		//Passwd: os.Getenv("DBPASS"),
		User:   "root",
		Passwd: "1",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
		// connector.go:95: could not use requested auth plugin 'mysql_native_password': this user requires mysql native password authentication.
		// 2021/10/14 09:40:00 this user requires mysql native password authentication.
		AllowNativePasswords: true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err) //exit
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr) //exit
	}
	fmt.Println("Connected!")
	return db
}
