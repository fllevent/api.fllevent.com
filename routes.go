package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Event struct {
	EventID   int
	EventName string
	Match     []Matches
}

type Matches struct {
	MatchID         int
	TeamName        string
	TeamNumber      int
	EventName       string
	MatchScoreOne   int
	MatchScoreTwo   int
	MatchScoreThree int
	Year            int
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
		Name := c.Param("name")
		EventRows, EventErr := db.Query("SELECT * FROM events WHERE eventName= ?", Name)
		handleErr(400, EventErr, c)
		var arrayEvent []Event
		for EventRows.Next() {
			var eventID int
			var eventName string
			EventErr = EventRows.Scan(&eventID, &eventName)
			handleErr(400, EventErr, c)
			MatchRows, MatchErr := db.Query("SELECT * FROM matches WHERE eventName= ?", eventName)
			handleErr(400, MatchErr, c)
			var matchesArray []Matches
			for MatchRows.Next() {
				var matchID int
				var teamName string
				var teamNumber int
				var eventName string
				var matchScoreOne int
				var matchScoreTwo int
				var matchScoreThree int
				var year int
				MatchErr = MatchRows.Scan(&matchID, &teamName, &teamNumber, &eventName, &matchScoreOne, &matchScoreTwo, &matchScoreThree, &year)
				handleErr(400, MatchErr, c)
				f := Matches{
					MatchID:         matchID,
					TeamName:        teamName,
					TeamNumber:      teamNumber,
					EventName:       eventName,
					MatchScoreOne:   matchScoreOne,
					MatchScoreTwo:   matchScoreTwo,
					MatchScoreThree: matchScoreThree,
					Year:            year,
				}
				matchesArray = append(matchesArray, f)
			}
			b := Event{
				EventID:   eventID,
				EventName: eventName,
				Match:     matchesArray,
			}
			arrayEvent = append(arrayEvent, b)
		}
		c.JSON(200, arrayEvent)
	}
	return gin.HandlerFunc(fn)
}

func getallevents(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		EventRows, EventErr := db.Query("SELECT * FROM events")
		handleErr(400, EventErr, c)
		var arrayEvent []Event
		for EventRows.Next() {
			var eventID int
			var eventName string
			EventErr = EventRows.Scan(&eventID, &eventName)
			handleErr(400, EventErr, c)
			MatchRows, MatchErr := db.Query("SELECT * FROM matches WHERE eventName= ?", eventName)
			handleErr(400, MatchErr, c)
			var matchesArray []Matches
			for MatchRows.Next() {
				var matchID int
				var teamName string
				var teamNumber int
				var eventName string
				var matchScoreOne int
				var matchScoreTwo int
				var matchScoreThree int
				var year int
				MatchErr = MatchRows.Scan(&matchID, &teamName, &teamNumber, &eventName, &matchScoreOne, &matchScoreTwo, &matchScoreThree, &year)
				handleErr(400, MatchErr, c)
				f := Matches{
					MatchID:         matchID,
					TeamName:        teamName,
					TeamNumber:      teamNumber,
					EventName:       eventName,
					MatchScoreOne:   matchScoreOne,
					MatchScoreTwo:   matchScoreTwo,
					MatchScoreThree: matchScoreThree,
					Year:            year,
				}
				matchesArray = append(matchesArray, f)
			}
			b := Event{
				EventID:   eventID,
				EventName: eventName,
				Match:     matchesArray,
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
