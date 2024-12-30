package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

type PipelineRunStatsDataHeaders struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

type PipelineRunStatsData struct {
	Data    []PipelineWorkerStatistics    `json:"data"`
	Headers []PipelineRunStatsDataHeaders `json:"headers"`
}

func getPipelineRunStats(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	idAsUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Println(err)
		return
	}

	statsData, err := getPipelineRunStatsData(uint(idAsUint))
	if err != nil {
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(statsData)
	if err != nil {
		log.Println(err)
	}
}

func getPipelineRunStatsData(pipelineWorkerId uint) (*PipelineRunStatsData, error) {
	var statsData PipelineRunStatsData
	jobsStatistics, err := fetchPipelineWorkerJobsStatistics(pipelineWorkerId)
	if err != nil {
		return nil, err
	}

	var pipelineStatistics PipelineWorkerStatistics
	err = fetchPipelineWorkerStatistics(&pipelineStatistics, pipelineWorkerId)
	if err != nil {
		return nil, err
	}
	jobsStatistics = append(jobsStatistics, pipelineStatistics)

	statsData.Data = jobsStatistics
	statsData.Headers = getPipelineRunStatsHeaders()

	return &statsData, nil
}

func getPipelineRunStatsHeaders() []PipelineRunStatsDataHeaders {
	headers := []PipelineRunStatsDataHeaders{
		{ID: "id", Name: "ID"},
		{ID: "name", Name: "Name"},
		{ID: "status", Name: "Status", Type: "status"},
		{ID: "user", Name: "Initializer"},
		{ID: "numberOfErrors", Name: "Total errors"},
		{ID: "startedAt", Name: "Started at", Type: "date"},
		{ID: "finishedAt", Name: "Finished at", Type: "date"},
		{ID: "totalDuration", Name: "Total duration"},
		{ID: "averageDuration", Name: "Average duration"},
	}

	return headers
}

func getPipelineRunStatsExcel(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	idAsUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Println(err)
		return
	}

	statsData, err := getPipelineRunStatsData(uint(idAsUint))
	if err != nil {
		log.Println(err)
		return
	}

	f, err := createExcelFile(statsData)
	if err != nil {
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	fileName := fmt.Sprintf("filename=\"Pipeline %s statistics.xlsx\"", statsData.Data[len(statsData.Data)-1].Name)
	writer.Header().Set("Content-Disposition", "attachment; "+fileName)

	_, err = f.WriteTo(writer)
	if err != nil {
		log.Println(err)
	}
}

func createExcelFile(statsData *PipelineRunStatsData) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	sheet := "Sheet1"

	// styles
	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "0369A1", Style: 1},
			{Type: "right", Color: "0369A1", Style: 1},
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"0369A1"}, Pattern: 1},
		Font: &excelize.Font{Color: "FFFFFF", Bold: true},
	})
	if err != nil {
		return nil, err
	}

	column, _ := excelize.ColumnNumberToName(len(statsData.Headers))
	f.SetCellStyle(sheet, "A1", column+"1", styleHeader)

	stylePipeline, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "6366F1", Style: 1},
			{Type: "right", Color: "6366F1", Style: 1},
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"818CF8"}, Pattern: 1},
		Font: &excelize.Font{Color: "FFFFFF", Bold: true},
	})
	if err != nil {
		return nil, err
	}

	row := strconv.FormatInt(int64(len(statsData.Data)+1), 10)
	f.SetCellStyle(sheet, "A"+row, column+row, stylePipeline)

	styleData, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	row2 := strconv.FormatInt(int64(len(statsData.Data)), 10)

	f.SetCellStyle(sheet, "A2", column+row2, styleData)

	styleStatusFailed, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FEE2E2"}, Pattern: 1},
	})
	styleStatusCompleted, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"DCFCE7"}, Pattern: 1},
	})
	styleStatusCancelled, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"E5E7EB"}, Pattern: 1},
	})
	styleStatusPlanned, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FDE68A"}, Pattern: 1},
	})
	styleStatusRunning, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"DBEAFE"}, Pattern: 1},
	})
	styleStatusPipelineFailed, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "FEE2E2", Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"818CF8"}, Pattern: 1},
	})
	styleStatusPipelineCompleted, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "DCFCE7", Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"818CF8"}, Pattern: 1},
	})
	styleStatusPipelineCancelled, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "E5E7EB", Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"818CF8"}, Pattern: 1},
	})
	styleStatusPipelinePlanned, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "FDE68A", Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"818CF8"}, Pattern: 1},
	})
	styleStatusPipelineRunning, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "DBEAFE", Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"818CF8"}, Pattern: 1},
	})

	// Headers
	for i, header := range statsData.Headers {
		column, _ := excelize.ColumnNumberToName(i + 1)
		cell := column + "1"
		f.SetCellValue(sheet, cell, header.Name)
	}

	// Data
	pipelineRowI := len(statsData.Data) - 1
	for i, row := range statsData.Data {
		cellRow := strconv.FormatInt(int64(i+2), 10)

		f.SetCellValue(sheet, "A"+cellRow, row.ID)
		f.SetCellValue(sheet, "B"+cellRow, row.Name)
		f.SetCellValue(sheet, "C"+cellRow, row.StatusId.RussianString())
		if row.User != nil {
			f.SetCellValue(sheet, "D"+cellRow, *row.User)
		}
		if row.NumberOfErrors != nil {
			f.SetCellValue(sheet, "E"+cellRow, *row.NumberOfErrors)
		}
		if row.StartedAt != nil {
			f.SetCellValue(sheet, "F"+cellRow, *row.StartedAt)
		}
		if row.FinishedAt != nil {
			f.SetCellValue(sheet, "G"+cellRow, *row.FinishedAt)
		}
		if row.TotalDuration != nil {
			f.SetCellValue(sheet, "H"+cellRow, *row.TotalDuration)
		}
		if row.AverageDuration != nil {
			f.SetCellValue(sheet, "I"+cellRow, *row.AverageDuration)
		}

		statusStyle := 0
		switch row.Status {
		case "running":
			if i != pipelineRowI {
				statusStyle = styleStatusRunning
			} else {
				statusStyle = styleStatusPipelineRunning
			}
		case "planned":
			if i != pipelineRowI {
				statusStyle = styleStatusPlanned
			} else {
				statusStyle = styleStatusPipelinePlanned
			}
		case "failed":
			if i != pipelineRowI {
				statusStyle = styleStatusFailed
			} else {
				statusStyle = styleStatusPipelineFailed
			}
		case "completed":
			if i != pipelineRowI {
				statusStyle = styleStatusCompleted
			} else {
				statusStyle = styleStatusPipelineCompleted
			}
		case "cancelled":
			if i != pipelineRowI {
				statusStyle = styleStatusCancelled
			} else {
				statusStyle = styleStatusPipelineCancelled
			}
		default:
		}

		if statusStyle != 0 {
			f.SetCellStyle(sheet, "C"+cellRow, "C"+cellRow, statusStyle)
		}

	}

	// Autofit all columns
	cols, err := f.GetCols(sheet)
	if err != nil {
		return nil, err
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return nil, err
		}
		f.SetColWidth(sheet, name, name, float64(largestWidth))
	}

	//err = f.SaveAs("Book1.xlsx")

	return f, nil
}
