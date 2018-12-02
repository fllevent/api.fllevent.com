package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const versionNumber = "1.0.0.0"

var (
	help       bool
	version    bool
	portNumber int

	DB *sql.DB
)

// User
type User struct {
	UserName string
	ID       int
}

func init() {
	flag.BoolVar(&help, "help", false, "Prints out available comands")
	flag.BoolVar(&version, "version", false, "Prints the version number")
	flag.IntVar(&portNumber, "port", 8000, "Set custome port number")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := dbconnect(dbhost(), dbname(), dbusername(), dbpassword())
	DB = db
	if err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()

	if help {
		printHelp()
	}
	if version {
		printVersion()
	}

	webServerPort := fmt.Sprintf(":%d", portNumber)

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userName := loginVals.Username
			Password := loginVals.Password

			UserRows, UserErr := DB.Query("SELECT * FROM users WHERE username= ?", userName)
			handleErr(400, UserErr, c)
			for UserRows.Next() {
				var userID int
				var username string
				var level int
				var password string
				UserErr = UserRows.Scan(&userID, &username, &level, &password)
				if (username == userName) || (CheckPasswordHash(Password, password)) {
					return &User{
						UserName: username,
						ID:       userID,
					}, nil
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok {
				UserRows, UserErr := DB.Query("SELECT * FROM users WHERE username= ?", v.UserName)
				handleErr(400, UserErr, c)
				for UserRows.Next() {
					var userID int
					var username string
					var level int
					var password string
					UserErr = UserRows.Scan(&userID, &username, &level, &password)

					if level >= 1 {
						return true
					}
					return false
				}
				return false
			}

			return false

		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.POST("/login", authMiddleware.LoginHandler)

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	authRouter := router.Group("/api/v1/auth")

	// Refresh time can be longer than token timeout
	authRouter.GET("/refresh_token", authMiddleware.RefreshHandler)
	authRouter.Use(authMiddleware.MiddlewareFunc())
	{
		authRouter.GET("/hello", helloHandler)
		authRouter.POST("/event/addevent", addevent(DB))
		authRouter.POST("/event/removeevent", removeevent(DB))
		authRouter.POST("/matches/addmatch", addMatch(DB))
		authRouter.POST("/matches/removematch", removeMatch(DB))

	}

	router.GET("/api/v1/healthcheck", healthcheck())
	router.POST("/api/v1/event/allevents", getallevents(DB))
	router.GET("/api/v1/event/allevents", getallevents(DB))
	router.POST("/api/v1/event/singleevent/:name", getsingleevent(DB))
	router.GET("/api/v1/event/singleevent/:name", getsingleevent(DB))
	router.GET("/api/v1/team/allteams", getteams(DB))
	router.POST("/api/v1/team/allteams", getteams(DB))
	router.GET("/api/v1/team/singleteam/:number", getteam(DB))
	router.POST("/api/v1/team/singleteam/:number", getteam(DB))
	router.POST("/api/v1/user/newuser", newUser(DB))

	fmt.Println("Starting Server on port" + webServerPort)
	router.Run(webServerPort) // listen and serve on port
}
