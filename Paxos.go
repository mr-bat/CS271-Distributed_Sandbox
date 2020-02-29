package main

type BallotNum struct {
	num int
	id  int
}

var lastBallot BallotNum

func isGreaterBallot(bn BallotNum) bool {
	if bn.num > lastBallot.num {
		return true
	} else if bn.num == lastBallot.num && bn.id > lastBallot.id {
		return true
	}
	return false

}

func beginSync() {
}
