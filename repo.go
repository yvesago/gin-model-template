package models

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"log"
	"regexp"
	"strings"
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

func ParseQuery(q map[string][]string) string {
	query := " "
	if q["_filters"] != nil {
		data := make(map[string]string)
		err := json.Unmarshal([]byte(q["_filters"][0]), &data)
		if err == nil {
			query = query + " WHERE "
			var searches []string
			for col, search := range data {
				valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
				if col != "" && search != "" && valid.MatchString(col) && valid.MatchString(search) {
					searches = append(searches, col+" LIKE \"%"+search+"%\"")
				}
			}
			query = query + strings.Join(searches, " AND ") // TODO join with OR for same keys
		}
	}
	if q["_sortField"] != nil && q["_sortDir"] != nil {
		sortField := q["_sortField"][0]
		// prevent SQLi
		valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
		if !valid.MatchString(sortField) {
			sortField = ""
		}
		if sortField == "created" || sortField == "updated" { // XXX trick for sqlite
			sortField = "datetime(" + sortField + ")"
		}
		sortOrder := q["_sortDir"][0]
		if sortOrder != "ASC" {
			sortOrder = "DESC"
		}
		if sortField != "" {
			query = query + " ORDER BY " + sortField + " " + sortOrder
		}
	}
	// _page, _perPage : LIMIT + OFFSET
	if q["_perPage"] != nil {
		perPage := q["_perPage"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		if valid.MatchString(perPage) {
			query = query + " LIMIT " + perPage
		}
	}
	if q["_page"] != nil {
		page := q["_page"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		if valid.MatchString(page) {
			query = query + " OFFSET " + page
		}
	}
	return query
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
