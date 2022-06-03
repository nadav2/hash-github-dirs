package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var port = "8080"

// ApiGetFile is the structure of the json getFile object
type ApiGetFile struct {
	FileName string `json:"fileName"`
}

// ApiHashFiles is the structure of the json hash object
type ApiHashFiles struct {
	Files []string `json:"files"`
}

type returnObject struct {
	value string
	error string
}

func handleError(error string) returnObject {
	return returnObject{
		value: "",
		error: error,
	}
}

func sendError(c *gin.Context, err string) {
	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
}


func getFileContent(fileName string) returnObject {
	apiDetails := request("http://init_api:8081/details/")
	if apiDetails.error != "" {
		return handleError(apiDetails.error)
	}

	// convert the json string to a map
	result, err := stringToJson(apiDetails.value)
	if err != "" {
		return handleError(err)
	}

	GitRef := result["gitRef"]
	Branch := result["branch"]

	url := parseUrl(GitRef, Branch, fileName)
	if url.error != "" {
		return handleError(url.error)
	}

	return request(url.value)
}

func stringToJson(jsonString string) (map[string]string, string) {
	// Declared an empty map interface
	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return map[string]string{}, "The data is invalid"
	}

	return result, ""
}

// parseUrl parse the arguments and return the url to the row file
func parseUrl(girRef string, branch string, fileName string) returnObject {
	if girRef == "" || branch == "" {
		return handleError("Error: the gitRef or branch is not initialized")
	}

	details := strings.Split(girRef, "/")
	if len(details) < 5 {
		return handleError("The github ref is invalid")
	}

	repoName := details[3]
	branchName := details[4][:len(details[4])-4]
	url := "https://raw.githubusercontent.com/" + repoName + "/" + branchName + "/" + branch + "/" + fileName

	return returnObject{value: url}
}

// return the request from a web page
func request(url string) returnObject {
	// get text content from web page with request
	resp, err := http.Get(url)
	if err != nil {
		return handleError("The request failed")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) == "404: Not Found" {
		return handleError("Error: file not found")
	}

	if err != nil {
		return handleError("An error occurred")
	}

	return returnObject{value: string(body)}
}

// hash text by using the sha256 algorithm
func hash(text string) string {
	shaObject := sha256.New()
	shaObject.Write([]byte(text))
	return fmt.Sprintf("%x", shaObject.Sum(nil))
}

// hash list of files
func hashFiles(listOfFiles []string) returnObject {
	bigHush := ""

	for _, file := range listOfFiles {
		content := getFileContent(file)
		if content.error != "" {
			return handleError(content.error)
		}

		bigHush += hash(content.value)
	}
	return returnObject{value: hash(bigHush)}
}

// api to get file content from GitHub directory
func getFileApi(c *gin.Context) {
	var details ApiGetFile

	if c.BindJSON(&details) != nil {
		sendError(c, "bad request")
		return
	}

	fileContent := getFileContent(details.FileName)
	if fileContent.error != "" {
		sendError(c, fileContent.error)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"fileContent": fileContent.value})
}

// api to hash the content of the given list of files
func hashFilesApi(c *gin.Context) {
	var details ApiHashFiles

	if c.BindJSON(&details) != nil {
		sendError(c, "bad request")
		return
	}
	if len(details.Files) == 0 {
		sendError(c, "no files to hash")
		return
	}

	sha := hashFiles(details.Files)
	if sha.error != "" {
		sendError(c, sha.error)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"hash": sha.value})
}

// the main function start the api
func main() {
	router := gin.Default()

	// create api
	router.POST("/get_file_content", getFileApi)
	router.POST("/hash_files", hashFilesApi)

	// run the server
	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
