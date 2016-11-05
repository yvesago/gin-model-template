package models

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
	//"log"
	"regexp"
	"strconv"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/User/NewStruct/g
  :%s/user/newinst/g

**/

// XXX custom struct name and fields
type User struct {
	Id         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"mail"`
	Status     string    `db:"status" json:"status"`
	Comment     string    `db:"name:comment, size:16384" json:"comment"`
	Pass string    `db:"pass" json:"pass"`
	Created    time.Time `db:"created" json:"created"` // or int64
	Updated    time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func (a *User) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

func (a *User) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	return nil
}

// REST handlers

func GetUsers(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	query := "SELECT * FROM user"

	// Parse query string
	//  receive : map[_filters:[{"q":"wx"}] _sortField:[id] ...
	q := c.Request.URL.Query()
	//log.Println(q)
	if q["_filters"] != nil {
		re := regexp.MustCompile("{\"([a-zA-Z0-9_]+?)\":\"([a-zA-Z0-9_. ]+?)\"}")
		r := re.FindStringSubmatch(q["_filters"][0])
		// TODO: special col name for all fields via reflections
		col := r[1]
		search := r[2]
		if col != "" && search != "" {
			query = query + " WHERE " + col + " LIKE \"%" + search + "%\" "
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
		// _page, _perPage, _sortDir, _sortField
		if sortField != "" {
			query = query + " ORDER BY " + sortField + " " + sortOrder
		}
	}
	//log.Println(" -- " + query)

	var users []User
	_, err := dbmap.Select(&users, query)

	if err == nil {
		c.JSON(200, users)
	} else {
		c.JSON(404, gin.H{"error": "no user(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/users
}

func GetUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i http://localhost:8080/api/v1/users/1
}

func PostUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var user User
	c.Bind(&user)

	//log.Println(user)

	if user.Name != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&user)
		if err == nil {
			c.JSON(201, user)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func UpdateUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)
	if err == nil {
		var json User
		c.Bind(&json)

		//log.Println(json)
		user_id, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		user := User{
			Id:         user_id,
			Pass:         json.Pass,
			Name:       json.Name,
			Email:       json.Email,
			Status:     json.Status,
			Comment:     json.Comment,
			Created:    user.Created, //user read from previous select
		}

		if user.Name != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&user)
			if err == nil {
				c.JSON(200, user)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func DeleteUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&user)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}
