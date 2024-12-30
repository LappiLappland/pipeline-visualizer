package main

import (
	"bufio"
	"math/rand/v2"
	"os"
	"time"
)

var randomSentences []string

func initCreateLog() {
	randomSentences, _ = readLines("randomSentences.txt")
}

type LogTotalResult struct {
	Total uint
}

// Creates mock logs
func createLog(stageWorkerId uint, shouldFail bool) []StageLog {

	times := uint(1 + rand.IntN(5))

	var result LogTotalResult
	db.Model(&StageLog{}).Where("stage_worker_id = ?", stageWorkerId).Select("count(id) as total").Group("stage_worker_id").Scan(&result)

	createdAt := time.Now()

	logs := make([]StageLog, times)
	for i := uint(0); i < times; i++ {

		usedType := logTypeBase
		if shouldFail && i == times-1 {
			usedType = logTypeError
		} else {
			luck := rand.IntN(100)
			if luck > 80 {
				usedType = logTypeWarn
			} else if luck > 60 {
				usedType = logTypeInfo
			}
		}

		logs[i] = StageLog{
			StageWorkerId: stageWorkerId,
			Text:          randomSentences[rand.IntN(len(randomSentences))],
			Line:          result.Total + i,
			LogTypeId:     uint(usedType),
			CreatedAt:     createdAt,
		}
	}

	db.Create(logs)

	return logs
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
