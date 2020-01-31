package main

import (
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

