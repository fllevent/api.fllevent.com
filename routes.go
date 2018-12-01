package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Event struct {
	EventID   int
	EventName string
	Owner     int
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

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type newuser struct {
	UserName string
	Level    int
	Password string
}

var identityKey = "id"

func helloHandler(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	UserRows, UserErr := DB.Query("SELECT * FROM users WHERE username= ?", user.(*User).UserName)
	handleErr(400, UserErr, c)
	for UserRows.Next() {
		var userID int
		var username string
		var level int
		var password string
		UserErr = UserRows.Scan(&userID, &username, &level, &password)
		c.JSON(200, gin.H{
			"userID":   userID,
			"userName": username,
			"text":     "Hello.",
		})
	}
	// c.JSON(200, gin.H{
	// 	"error": user,
	// })
}

var healthOK = fmt.Sprintf("{\"message\":\"ok\", \"version\":\"%s\"}\n", versionNumber)

func healthcheck() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, healthOK)
	}

	return gin.HandlerFunc(fn)
}

func newUser(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newuserInternal newuser
		if c.ShouldBind(&newuserInternal) == nil {
			newUserStmt, newUserStmtErr := db.Prepare("INSERT users (username, level, password) VALUES (?, ?, ?)")
			handleErr(400, newUserStmtErr, c)

			pass, passerr := HashPassword(newuserInternal.Password)
			handleErr(400, passerr, c)

			newUserRes, newuserErr := newUserStmt.Exec(newuserInternal.UserName, newuserInternal.Level, pass)
			handleErr(400, newuserErr, c)

			id, err := newUserRes.LastInsertId()
			handleErr(400, err, c)

			c.JSON(200, gin.H{
				"id":       id,
				"UserName": newuserInternal.UserName,
				"Level":    newuserInternal.Level,
			})

		}
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
			stmt, err := db.Prepare("INSERT events (eventname, eventOwner) VALUES (?, ?)")
			handleErr(400, err, c)
			if (event.EventName != "") && (event.Owner != 0) {

				res, err := stmt.Exec(event.EventName, event.Owner)
				handleErr(400, err, c)

				id, err := res.LastInsertId()
				handleErr(400, err, c)

				c.JSON(200, gin.H{
					"sucess":      "sucess",
					"id":          id,
					"Event Name":  event.EventName,
					"Event Owner": event.Owner,
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
			var Owner int
			var eventName string
			EventErr = EventRows.Scan(&eventID, &Owner, &eventName)
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
				Owner:     Owner,
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
			var Owner int
			var eventName string
			EventErr = EventRows.Scan(&eventID, &Owner, &eventName)
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
				Owner:     Owner,
				Match:     matchesArray,
			}
			arrayEvent = append(arrayEvent, b)
		}
		c.JSON(200, arrayEvent)
	}
	return gin.HandlerFunc(fn)
}

func getteam(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		teamnumber := c.Param("number")
		teamRows, teamErr := db.Query("SELECT * FROM matches WHERE teamNumber = ?", teamnumber)
		handleErr(400, teamErr, c)
		var matchesArray []Matches
		for teamRows.Next() {
			var matchID int
			var teamName string
			var teamNumber int
			var eventName string
			var matchScoreOne int
			var matchScoreTwo int
			var matchScoreThree int
			var year int
			teamErr = teamRows.Scan(&matchID, &teamName, &teamNumber, &eventName, &matchScoreOne, &matchScoreTwo, &matchScoreThree, &year)
			handleErr(400, teamErr, c)
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
		c.JSON(200, matchesArray)
	}
	return gin.HandlerFunc(fn)
}

func getteams(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		teamRows, teamErr := db.Query("SELECT * FROM matches")
		handleErr(400, teamErr, c)
		var matchesArray []Matches
		for teamRows.Next() {
			var matchID int
			var teamName string
			var teamNumber int
			var eventName string
			var matchScoreOne int
			var matchScoreTwo int
			var matchScoreThree int
			var year int
			teamErr = teamRows.Scan(&matchID, &teamName, &teamNumber, &eventName, &matchScoreOne, &matchScoreTwo, &matchScoreThree, &year)
			handleErr(400, teamErr, c)
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
		c.JSON(200, matchesArray)
	}
	return gin.HandlerFunc(fn)
}

func handleErr(errorCode int, err error, c *gin.Context) {
	if err != nil {
		c.JSON(errorCode, err)
	}
}
