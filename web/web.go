package web

import (
	"fmt"
	"net/http"

	"wiremelt/wiremelt"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func InitHTTPServer(include interface{}) {
	fmt.Println(color.HiBlueString(fmt.Sprintf("\n~ INIT_HTTPServer_API: %v\n", include)))

	router := gin.Default()
	router.GET("/sessions", getSessions)
	router.POST("/sessions", postSessions)

	router.Run("localhost:4444")
}

// getSessions responds with the list of all sessions as JSON.
func getSessions(c *gin.Context) {
	sessions := wiremelt.GetSessions()
	c.IndentedJSON(http.StatusOK, sessions)
}

// postSessions adds an session from JSON received in the request body.
func postSessions(c *gin.Context) {
	var newSession wiremelt.SessionConfiguration

	// Call BindJSON to bind the received JSON to
	// newSession.
	if err := c.BindJSON(&newSession); err != nil {
		return
	}

	sessions := wiremelt.GetSessions()
	// Add the new session to the slice.
	sessions = append(sessions, newSession)
	c.IndentedJSON(http.StatusCreated, newSession)
}
