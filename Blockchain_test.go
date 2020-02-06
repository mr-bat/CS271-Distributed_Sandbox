package main

import (
	"reflect"
	"testing"
)


func CheckBalanceCorrectness(users []string, balances []int) bool {
	for i, user := range users{
		bal := getBalance(user)

		if bal != balances[i] {
			return false
		}
	}

	return true
}

var blocks = []Block{
	{"1", "2", 1, 1},
	{"1", "3", 1, 2},
	{"3", "2", 1, 3},
}

func TestBlockchain(t *testing.T) {
	users := []string {
		"1",
		"2",
		"3",
		"4",
	}

	expectedBalances := [][]int{
		{9, 11, 10, 10},
		{8, 11, 11, 10},
		{8, 12, 10, 10},
	}

	initializeTimetable(4)
	for i, block := range blocks {
		addBlock(block.sender, block.receiver, block.amount)
		if !CheckBalanceCorrectness(users, expectedBalances[i]) {
			for i := range users {
				t.Errorf("Balance %s was incorrect, got: %d, want: %d.", users[i], getBalance(users[i]), expectedBalances[i])
			}
		}
	}
}

func TestBlockchainParser(t *testing.T) {
	if "1&2&1&1\n" != blocks[0].toString() {
		t.Errorf("Incorrect conversion from single block to string: %s", blocks[0].toString())
	}
	if blocks[0] != parseBlock("1&2&1&1") {
		t.Error("Incorrect conversion from string to single block")
	}

	convertedBlocks := rangeToString(0)
	if convertedBlocks != blocks[0].toString() + blocks[1].toString() + blocks[2].toString() {
		t.Error("Incorrect conversion from range of blocks to string")
		t.Error(convertedBlocks)
		t.Error(blocks)
	}
	if !reflect.DeepEqual(blocks, parseRange(convertedBlocks)) {
		t.Error("Incorrect back and forth conversion of range of blocks")
		t.Error("\tblocks:")
		t.Error(blocks)
		t.Error("\tparsed blocks:")
		t.Error(parseRange(convertedBlocks))
	}
}
