package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

// searchKey is the name of the query string parameter which will be used to search emoji data.
const searchKey = "search"

// EmojiEntry represents the emoji data provided, per entry, in the source JSON file.
// Source: https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json.
type EmojiEntry struct {
	Emoji          string   `json:"emoji"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Aliases        []string `json:"aliases"`
	Tags           []string `json:"tags"`
	UnicodeVersion string   `json:"unicode_version"`
	IOSVersion     string   `json:"ios_version"`
}

func main() {
	lambda.Start(handler)
}

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	emojis, err := getEmojis(event.QueryStringParameters)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       strings.Join(emojis, ""),
	}, nil
}

// getEmojis searches availabe emoji tags and aliases for the given search term.
func getEmojis(qs map[string]string) ([]string, error) {
	search, ok := qs[searchKey]
	if !ok {
		return nil, errors.Errorf("%q is a required query string parameter", searchKey)
	}

	emojis, err := getData(os.Getenv("SOURCE_URL"))
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve emoji data")
	}

	se := make(map[string][]string)
	for _, e := range emojis {
		for _, alias := range e.Aliases {
			if _, ok := se[alias]; !ok {
				se[alias] = []string{e.Emoji}
			} else {
				se[alias] = append(se[alias], e.Emoji)
			}
		}

		for _, tag := range e.Tags {
			if _, ok := se[tag]; !ok {
				se[tag] = []string{e.Emoji}
			} else {
				se[tag] = append(se[tag], e.Emoji)
			}
		}
	}

	result, ok := se[search]
	if !ok {
		return nil, errors.Errorf("no results for %q", search)
	}

	return result, nil
}

// getData retrieves the emoji data from the url and unmarshals it into a slice of EmojiData.
func getData(url string) ([]EmojiEntry, error) {
	client := http.Client{Timeout: time.Second * 2}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create the request")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not make the request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read the response")
	}

	var emojis []EmojiEntry
	if err = json.Unmarshal(body, &emojis); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal the response")
	}

	return emojis, nil
}
