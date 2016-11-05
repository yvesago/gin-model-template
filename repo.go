package models

import (
	"database/sql"
     "github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"log"
)

// gin Middlware to select database
func Database(connString string) gin.HandlerFunc {
	dbmap := InitDb(connString)
	return func(c *gin.Context) {
		c.Set("DBmap", dbmap)
		c.Next()
	}
}

func InitDb(dbName string) *gorp.DbMap {
    // XXX fix database type
	db, err := sql.Open("sqlite3", dbName)
	checkErr(err, "sql.Open failed")
	//dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    // XXX fix tables names
	dbmap.AddTableWithName(Agent{}, "Agent").SetKeys(true, "Id")
	dbmap.AddTableWithName(User{}, "User").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
