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
  :%s/Agent/NewStruct/g
  :%s/agent/newinst/g

**/

// XXX custom struct name and fields
type Agent struct {
	Id         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	IP         string    `db:"ip" json:"ip"`
	FileSurvey string    `db:"filesurvey" json:"filesurvey"`
	Role       string    `db:"role" json:"role"`
	Status     string    `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"` // or int64
	Updated    time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func (a *Agent) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

func (a *Agent) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	return nil
}

// REST handlers

func GetAgents(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	query := "SELECT * FROM agent"

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

	var agents []Agent
	_, err := dbmap.Select(&agents, query)

	if err == nil {
		c.JSON(200, agents)
	} else {
		c.JSON(404, gin.H{"error": "no agent(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/agents
}

func GetAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var agent Agent
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, agent)
	} else {
		c.JSON(404, gin.H{"error": "agent not found"})
	}

	// curl -i http://localhost:8080/api/v1/agents/1
}

func PostAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var agent Agent
	c.Bind(&agent)

	//log.Println(agent)

	if agent.Name != "" && agent.IP != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&agent)
		if err == nil {
			c.JSON(201, agent)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/agents
}

func UpdateAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var agent Agent
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE id=?", id)
	if err == nil {
		var json Agent
		c.Bind(&json)

		//log.Println(json)
		agent_id, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		agent := Agent{
			Id:         agent_id,
			IP:         json.IP,
			Name:       json.Name,
			Role:       json.Role,
			FileSurvey: json.FileSurvey,
			Status:     json.Status,
			Created:    agent.Created, //agent read from previous select
		}

		if agent.Name != "" && agent.IP != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&agent)
			if err == nil {
				c.JSON(200, agent)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "agent not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/agents/1
}

func DeleteAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var agent Agent
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&agent)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "agent not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/agents/1
}
