package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "strings"
import "gopkg.in/russross/blackfriday.v2"
import "context"
import "github.com/aws/aws-lambda-go/lambda"
import "os"

type Response struct {
	Files map[string]File `json:"files"`
}
type File struct {
	Content string `json:"content"`
}

type GistID struct {
	ID string `json:"id"`
}

func HandleRequest(ctx context.Context, gist GistID) ([]string, error) {
	gistId := gist.ID
	github_token := os.Getenv("GITHUB_TOKEN")
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistId)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", github_token)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var responseObject Response
	json.Unmarshal(body, &responseObject)

	keys := make([]string, 0, len(responseObject.Files))
	for k := range responseObject.Files {
		keys = append(keys, k)
	}

	presentation := responseObject.Files[keys[0]].Content

	mdSlice := strings.Split(presentation, "---")

	var htmlSlice []string

	for i := 0; i < len(mdSlice); i++ {
		slide := string(blackfriday.Run([]byte(mdSlice[i])))
		htmlSlice = append(htmlSlice, string(slide))
	}

	return htmlSlice, nil
}

func main() {
	lambda.Start(HandleRequest)
}
