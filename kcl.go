package main

// ref. https://github.com/harlow/kinesis-consumer

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	consumer "github.com/harlow/kinesis-consumer"
)

// struct for Logfiles comes from API servers
type Log struct {
	Message string `json:"message"`
}

func main() {
	var stream = flag.String("stream", "", "Stream name")
	var outputFile = flag.String("outputFile", "", "File to send the IPS Service")
	flag.Parse()

	var logs []Log
	var line []byte

	// consumer
	c, err := consumer.New(*stream)
	if err != nil {
		log.Fatalf("consumer error: %v", err)
	}

	// open logfile to send ips
	file, err := os.OpenFile(*outputFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()

	// start
	err = c.Scan(context.TODO(), func(r *consumer.Record) bool {
		// decode json
		//		line = append(line, '[')
		line = []byte("[")
		line = append(line, string(r.Data)...)
		line = append(line, ']')

		if err := json.Unmarshal(line, &logs); err != nil {
			// fmt.Println(string(r.Data))
			log.Fatal(err)
		}
		// デコードしたデータを表示
		for _, p := range logs {
			// fmt.Printf("%s\n", p.Message)
			fmt.Fprintln(file, p.Message)
		}
		return true // continue scanning
	})
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}

	// ddb checkpoint
	// var appName = "appName"
	// var tableName = "checkpoint"
	// ck, err := checkpoint.New(appName, tableName)
	// if err != nil {
	// 	log.Fatalf("new checkpoint error: %v", err)
	// }

	// Note: If you need to aggregate based on a specific shard the `ScanShard`
	// method should be leverged instead.
}
