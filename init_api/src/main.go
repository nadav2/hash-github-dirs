package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var GitRef string
var Branch string

// port to listen on
var port = "8081"

// ApiStartPoint is the structure of the json initialize object
type ApiStartPoint struct {
	Ref    string `json:"ref"`
	Branch string `json:"branch"`
}

// CheckOutRef set the git ref and branch for the api
func CheckOutRef(path string, branch string) {
	GitRef = "https://github.com/" + path + ".git"

	if branch == "" {
		branch = "master"
	}
	Branch = branch
}

// sendError send the error to the client
func sendError(c *gin.Context, err string) {
	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
}

// initApi initialize the api
func initApi(c *gin.Context) {
	var details ApiStartPoint
	if c.BindJSON(&details) == nil {
		CheckOutRef(details.Ref, details.Branch)
		c.IndentedJSON(http.StatusCreated, details)
	} else {
		sendError(c, "Invalid json")
	}
}

// getApiData returns the current git ref and branch
func getApiData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"gitRef": GitRef, "branch": Branch})
}

// the main function start the api
func main() {
	router := gin.Default()

	// create api
	router.POST("/check_out_ref", initApi)
	router.GET("/details", getApiData)

	// run the server
	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}

}
