package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	elasticsearchV8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

const mapping = `
{
    "mappings": {
        "properties": {
            "word": {
                "type": "completion"
            }
        }
    }
}`

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	ctx := context.Background()
	es := getESClient()

	createIndexRequest := esapi.IndicesCreateRequest{
		Index: "words",
		Body:  bytes.NewReader([]byte(mapping)),
	}

	resCreateIndex, err := createIndexRequest.Do(ctx, es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer resCreateIndex.Body.Close()

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         "words",          // The default index name
		Client:        es,               // The Elasticsearch client
		NumWorkers:    50,               // The number of worker goroutines
		FlushBytes:    int(5e+6),        // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	var increment int = 1
	var countSuccessful uint64

	start := time.Now().UTC()

	for word, _ := range getWordsList() {
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action:     "index",
				DocumentID: strconv.Itoa(increment + 1),
				Body:       bytes.NewReader([]byte(`{"word" : "` + word + `"}`)),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}

		increment++
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	biStats := bi.Stats()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}

	<-done
	fmt.Println("exiting")
}

func getESClient() *elasticsearchV8.Client {
	cfg := elasticsearchV8.Config{
		Addresses: []string{
			"http://elastic01:9200",
		},
	}
	es, err := elasticsearchV8.NewClient(cfg)

	if err != nil {
		log.Fatalln("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalln("Error getting response: %s", err)
	}

	defer res.Body.Close()
	log.Println(res)

	return es
}

func getWordsList() map[string]interface{} {
	jsonFile, err := ioutil.ReadFile("words_dictionary.json")

	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Successfully opened json")

	var wordsJson map[string]interface{}
	unmarshalErr := json.Unmarshal(jsonFile, &wordsJson)
	if unmarshalErr != nil {
		fmt.Printf(unmarshalErr.Error())
	}

	return wordsJson
}
