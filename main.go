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

type BreakServer struct {
	ServerName   string
	BreakTime    time.Time
	RecoveryTime time.Time
	BreakCount   int32
	IsBreak      bool
}

type Result struct {
	BreakHost string
	BreakSpan time.Duration
}

func main() {
	f, err := os.Open("log.txt")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)

	var bs []*BreakServer
	// var res []Result

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		/* 故障サーバを抽出・格納 */
		breakHost := record[1]
		if record[2] == "-" {
			for _, s := range bs {
				if s.ServerName == breakHost && s.IsBreak {
					s.BreakCount += 1
				}
			}
			bs = append(bs, &BreakServer{
				ServerName: record[1],
				BreakTime:  stringToTime(record[0]),
				IsBreak:    true,
				BreakCount: 1,
			})
		} else {
			for _, a := range bs {
				if a.ServerName == breakHost && a.IsBreak {
					a.RecoveryTime = stringToTime(record[0])
				}
			}
		}
		// fmt.Printf("%#v\n", record)
	}

	/* 故障サーバ名、故障期間を抽出 */
	for _, s := range bs {
		if s.RecoveryTime.IsZero() {
			continue
		}
		bt := s.RecoveryTime.Sub(s.BreakTime)
		fmt.Println("故障サーバー:", s.ServerName, "故障期間: ", bt)
		// res = append(res, Result{
		// 	BreakHost: s.ServerName,
		// 	BreakSpan: bt,
		// })
	}

	// fmt.Println(res)
}

/* 文字列をTime型に変換する */
func stringToTime(str string) time.Time {
	t, _ := time.Parse(dateFormat, str)
	return t
}
