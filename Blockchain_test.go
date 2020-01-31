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

func TestBlockchain(t *testing.T) {
	blocks := []Block{
		{"1", "2", 1},
		{"1", "3", 1},
		{"3", "2", 1},
	}

	users := []string {
		"1",
		"2",
		"3",
		"a",
	}

	expectedBalances := [][]int{
		{9, 11, 10, 10},
		{8, 11, 11, 10},
		{8, 12, 10, 10},
	}

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
	blocks := []Block{
		{"1", "2", 1},
		{"1", "3", 1},
	}

	if "1&2&1" != blocks[0].toString() {
		t.Error("Incorrect conversion from single block to string")
	}
	if blocks[0] != parseBlock("1&2&1") {
		t.Error("Incorrect conversion from string to single block")
	}

	addBlock("1", "2", 1)
	addBlock("1", "3", 1)

	convertedBlocks := rangeToString(0)
	if convertedBlocks != blocks[0].toString() + blocks[1].toString() {
		t.Error("Incorrect conversion from range of blocks to string")
	}
	if reflect.DeepEqual(blocks, parseRange(convertedBlocks)) {
		t.Error("Incorrect back and forth conversion of range of blocks")
	}
}

