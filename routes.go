package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Event struct {
	EventID   int
	EventName string
}

var healthOK = fmt.Sprintf("{\"message\":\"ok\", \"version\":\"%s\"}\n", versionNumber)

func healthcheck() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, "healthOK")
	}

	return gin.HandlerFunc(fn)
}

func removeevent(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var event Event
		if c.ShouldBind(&event) == nil {
			//delete
			stmt, err := db.Prepare("DELETE FROM events where eventname =?")
			handleErr(400, err, c)

			res, err := stmt.Exec(event.EventName)
			handleErr(400, err, c)

			affect, err := res.RowsAffected()
			handleErr(400, err, c)
			if affect == 0 {
				c.JSON(409, "failed")
			} else {
				c.JSON(200, gin.H{
					"sucess": "event removed",
					"affect": affect,
				})
			}
		}
	}

	return gin.HandlerFunc(fn)
}

func addevent(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var event Event
		if c.ShouldBind(&event) == nil {
			//insert
			stmt, err := db.Prepare("INSERT events SET eventname = ?")
			handleErr(400, err, c)
			if event.EventName != "" {

				res, err := stmt.Exec(event.EventName)
				handleErr(400, err, c)

				id, err := res.LastInsertId()
				handleErr(400, err, c)

				c.JSON(200, gin.H{
					"sucess":     "sucess",
					"id":         id,
					"Event Name": event.EventName,
				})
			} else {
				c.JSON(400, "No name Given")
			}
		}
	}
	return gin.HandlerFunc(fn)

}

func getsingleevent(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		name := c.Param("name")
		rows, err := db.Query("SELECT * FROM events WHERE eventname = ?", name)
		handleErr(400, err, c)
		var arrayEvent []Event
		for rows.Next() {
			var eventID int
			var eventName string
			err = rows.Scan(&eventID, &eventName)
			handleErr(400, err, c)
			b := Event{
				EventID:   eventID,
				EventName: eventName,
			}
			arrayEvent = append(arrayEvent, b)
		}
		c.JSON(200, arrayEvent)
	}
	return gin.HandlerFunc(fn)

}

func getallevents(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM events")
		handleErr(400, err, c)
		var arrayEvent []Event
		for rows.Next() {
			var eventID int
			var eventName string
			err = rows.Scan(&eventID, &eventName)
			handleErr(400, err, c)
			b := Event{
				EventID:   eventID,
				EventName: eventName,
			}
			arrayEvent = append(arrayEvent, b)
		}
		c.JSON(200, arrayEvent)
	}
	return gin.HandlerFunc(fn)

}

func handleErr(errorCode int, err error, c *gin.Context) {
	if err != nil {
		c.JSON(errorCode, err)
	}
}
