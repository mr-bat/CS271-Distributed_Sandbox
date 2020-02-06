package main

import (
	"github.com/adam-lavrik/go-imath/ix"
	"strconv"
	"strings"
)

var id int
var timeTable = make([][]int, 10)

func setId(_id int) { id = _id }
func getId() int { return id }
func getTime() int { return timeTable[getId()][getId()]}
func incTime() int {
	timeTable[getId()][getId()] += 1
	return timeTable[getId()][getId()]
}

func initializeTimetable(length int) {
	rows := make([]int, length * length)
	for i := 0; i < length; i++ {
		timeTable[i] = rows[i*length : (i+1)*length]
	}
}

func updateTimeTable(_timeTable [][]int, _id int) {
	for i := 0; i < len(timeTable); i++ {
		for j := 0; j < len(timeTable); j++ {
			timeTable[i][j] = ix.Max(timeTable[i][j], _timeTable[i][j])
		}
		timeTable[getId()][i] = ix.Max(timeTable[getId()][i], timeTable[_id][i])
	}
}

func convertTtToString() string {
	result := ""
	for i := 0; i < len(timeTable); i++ {
		for j := 0; j < len(timeTable); j++ {
			if j > 0 {
				result += "&"
			}
			result += string(timeTable[i][j])
		}
		if i + 1 < len(timeTable) {
			result += "\n"
		}
	}

	return result
}

func parseTt(convertedTt string) [][]int {
	Tt := make([][]int, len(timeTable))

	rows := strings.Split(convertedTt, "\n")

	for i := range rows {
		Tt[i] = make([]int, len(timeTable))
		entries := strings.Split(rows[i], "&")
		for j := range rows {
			Tt[i][j], _ = strconv.Atoi(entries[j])
		}
	}

	return Tt
}