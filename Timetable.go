package main

import (
	"github.com/adam-lavrik/go-imath/ix"
	"strconv"
	"strings"
)

var id int
var timetable [][]int

func setId(_id int) { id = _id }
func getId() int { return id }
func getTime() int { return timetable[getId()][getId()]}
func incTime() int {
	timetable[getId()][getId()] += 1
	return timetable[getId()][getId()]
}

func initializeTimetable(length int) {
	timetable = make([][]int, length)
	rows := make([]int, length * length)
	for i := 0; i < length; i++ {
		timetable[i] = rows[i*length : (i+1)*length]
	}
}

func updateTimetable(_timeTable [][]int, _id int) {
	for i := 0; i < len(timetable); i++ {
		for j := 0; j < len(timetable); j++ {
			timetable[i][j] = ix.Max(timetable[i][j], _timeTable[i][j])
		}
		timetable[getId()][i] = ix.Max(timetable[getId()][i], timetable[_id][i])
	}
}

func convertTtToString() string {
	result := ""
	for i := 0; i < len(timetable); i++ {
		for j := 0; j < len(timetable); j++ {
			if j > 0 {
				result += "&"
			}
			result += strconv.Itoa(timetable[i][j])
		}
		if i + 1 < len(timetable) {
			result += "\n"
		}
	}

	return result
}

func parseTt(convertedTt string) [][]int {
	Tt := make([][]int, len(timetable))

	rows := strings.Split(convertedTt, "\n")

	for i := range rows {
		Tt[i] = make([]int, len(timetable))
		entries := strings.Split(rows[i], "&")
		for j := range rows {
			Tt[i][j], _ = strconv.Atoi(entries[j])
		}
	}

	return Tt
}