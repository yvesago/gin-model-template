package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v2"
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
// Agent db and json type
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

// PreInsert set created an updated time before insert in db
func (a *Agent) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

// PreUpdate set updated time before insert in db
func (a *Agent) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	return nil
}

// REST handlers

// GetAgents return all agents filtered by URL query
func GetAgents(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM agent"

	// Parse query string
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM user  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM user")
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	if verbose == true {
		fmt.Println(q)
		fmt.Println("query: " + query)
	}

	var agents []Agent
	_, err := dbmap.Select(&agents, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, agents)
	} else {
		c.JSON(404, gin.H{"error": "no agent(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/agents
}

// GetAgent return one agent by id
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

// PostAgent create and return agent
func PostAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)

	var agent Agent
	c.Bind(&agent)

	if verbose == true {
		fmt.Println(agent)
	}

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

// UpdateAgent by id
func UpdateAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	id := c.Params.ByName("id")

	var agent Agent
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE id=?", id)
	if err == nil {
		var json Agent
		c.Bind(&json)
		if verbose == true {
			fmt.Println(json)
		}
		agentId, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		agent := Agent{
			Id:         agentId,
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

// DeleteAgent by id
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
