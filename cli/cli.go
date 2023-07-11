package cli

import (
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, strChan chan string, wg *sync.WaitGroup) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", translateUrl, nil)

	if err != nil {
		log.Fatalf("1. There was a problem: %s", err)
	}

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("2. There was a problem: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		strChan <- "You have been rate limited, try again later."
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(res.Body)
	if err != nil {
		log.Fatalf("3. There was a problem - %s", err)
	}

	nestOne, err := parsedJson.ArrayElement(0)
	if err != nil {
		log.Fatalf("4. There was a problem - %s", err)
	}

	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("5. There was a problem - %s", err)
	}

	translatedStr, err := nestTwo.ArrayElement(0)
	if err != nil {
		log.Fatalf("6. There was a problem - %s", err)
	}

	strChan <- translatedStr.Data().(string)
	wg.Done()

}
