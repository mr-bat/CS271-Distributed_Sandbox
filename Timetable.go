package main

import (
	"github.com/adam-lavrik/go-imath/ix"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var id int
var timetable [][]int

func setId(_id int) { id = _id }
func getId() int { return id }

func initializeTimetable(length int) {
	Logger.WithField("length", length).Info("initializing timetable")
	timetable = make([][]int, length)
	rows := make([]int, length * length)
	for i := 0; i < length; i++ {
		timetable[i] = rows[i*length : (i+1)*length]
	}
}

func updateTimetable(_timeTable [][]int, updaterId, updateeId int) {
	Logger.WithFields(logrus.Fields{
		"timetable" : timetable,
		"id": updaterId,
	}).Info("updating timetable")

	for i := 0; i < len(timetable); i++ {
		for j := 0; j < len(timetable); j++ {
			timetable[i][j] = ix.Max(timetable[i][j], _timeTable[i][j])
		}
		timetable[updateeId][i] = ix.Max(timetable[updateeId][i], _timeTable[updaterId][i])
	}

	Logger.WithFields(logrus.Fields{
		"timetable" : timetable,
	}).Info("updated timetable")
}

func pickNewBlocks(blocks []Block, rId int) []Block {
	Logger.WithFields(logrus.Fields{
		"receiver-id": rId,
	}).Info("picking blocks")

	var picked []Block
	for _, block := range blocks {
		Logger.WithFields(logrus.Fields{
			"block": block,
		}).Info("analyzing block")

		bId, _ := strconv.Atoi(block.sender)
		if 0 > timetable[rId][bId] {
			picked = append(picked, block)
			println("picked")
		}
	}
	return picked
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