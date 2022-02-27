package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const dateFormat = "20060102150405"

var failedServers []*FailedServer

type FailedServer struct {
	ServerName   string
	BreakTime    time.Time
	RecoveryTime time.Time
	IsBreak      bool
}

func main() {
	f, err := os.Open("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	/* CSVリーダー生成 */
	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		confirmTime, _ := time.Parse(dateFormat, record[0])
		serverName := record[1]
		responseResult := record[2]

		/* 故障サーバを抽出・格納 */
		if responseResult == "-" {
			for _, s := range failedServers {
				if s.ServerName == serverName && s.IsBreak {
					continue
				}
			}
			failedServers = append(failedServers, &FailedServer{
				ServerName: serverName,
				BreakTime:  confirmTime,
				IsBreak:    true,
			})
		} else {
			for _, s := range failedServers {
				if s.ServerName == serverName && s.IsBreak {
					s.RecoveryTime = confirmTime
				}
			}
		}
	}

	/* 復旧したサーバの故障サーバ名、故障期間を抽出 */
	for _, s := range failedServers {
		failurePeriod := s.RecoveryTime.Sub(s.BreakTime)
		fmt.Printf("故障サーバーIP: %s 故障期間: %s\n", s.ServerName, failurePeriod)
	}
}
