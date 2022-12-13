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

	router.Run("localhost:4444")
}

// getSessions responds with the list of all sessions as JSON.
func getSessions(c *gin.Context) {
	sessions := wiremelt.GetSessions()
	c.IndentedJSON(http.StatusOK, sessions)
}
