package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL Driver
)

func main() {
	//{
	//	"credentials": {
	//		"host": "192.168.1.122",
	//		"port": 8635,
	//		"name": "postgres",
	//		"username": "root",
	//		"password": "Huangwei!120521",
	//		"uri": "postgres://root:Huangwei!120521@192.168.1.122:8635/postgres?reconnect=true"
	//		}
	//}

	var address string = "192.168.1.122"
	var port int = 8635
	var dbname string= "postgres"
	var username string= "root"
	var password string= "Huangwei!120521"

	//connectionString := connectionString(address, port, dbname, username, password)
	var connectionString string = fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password='%s'", address, port, dbname, username, password)
	fmt.Println("sql-open: connection-string ", connectionString)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("sql-open: error: ", err)
		return
	}

	selectDatabaseStatement := "SELECT datname FROM pg_database WHERE datname='" + dbname + "'"
	fmt.Println("selectDatabaseStatement: ", selectDatabaseStatement)

	var dummy string
	err = db.QueryRow(selectDatabaseStatement).Scan(&dummy)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("sql-QueryRow: sql.ErrNoRows: ", err)
		db.Close()
		return
	case err != nil:
		fmt.Println("sql-QueryRow: error: ", err)
		db.Close()
		return
	}
	fmt.Println("sql-QueryRow: dummy", dummy)

	db.Close()
	return
}


//func connectionString(address string, port int64, dbname string, username string, password string) string {
//	return fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password='%s'", address, port, dbname, username, password)
//}
