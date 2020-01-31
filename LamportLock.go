package main

import (
	"context"
	"fmt"
	"github.com/wangjia184/sortedset"
	"golang.org/x/sync/semaphore"
	"strconv"
	"strings"
	"time"
)

const bufferSize = 128

var lamportOrderedQueue *sortedset.SortedSet
var orderedQueueSemaphore, clockSemaphore *semaphore.Weighted
var lamportTime sortedset.SCORE
var permittedEntry chan string
var busyWait chan string

func ParseString(command string) (string, int, string) {
	parsed := strings.Split(command, "&")
	time, _ := strconv.Atoi(parsed[2])
	return parsed[1], time, parsed[0]
}

func ParseAddress(addr string) (string, string) {
	parsed := strings.Split(addr, ":")
	return parsed[0], parsed[1]
}

func HandleMessage(message string) string {
	println("handling: " + message)
	address, time, command := ParseString(message)

	if command == "REQ" {
		AddToQueueCustom("", address, sortedset.SCORE(time))
		IncrementTime()

		ip, port := ParseAddress(address)
		sendClient(Addr{IP: ip, Port: port}, fmt.Sprintf("REP&%v&%v", getAddress(), time))
		IncrementTime()
	} else if command == "REP" {
		permittedEntry <- message
	} else if command == "REL" {
		lamportOrderedQueue.Remove(address)
		IncrementTime()
		busyWait <- "TRY"
	}

	return ""
}

func init() {
	lamportOrderedQueue = sortedset.New()
	orderedQueueSemaphore = semaphore.NewWeighted(int64(1))
	clockSemaphore = semaphore.NewWeighted(int64(1))
	lamportTime = 0
	permittedEntry = make(chan string, bufferSize)
	busyWait = make(chan string, bufferSize)
}

func IncrementTime() {
	clockSemaphore.Acquire(context.Background(), 1)
	lamportTime = sortedset.SCORE(int(lamportTime) + 1)
	println(fmt.Sprintf("incremented: %v", lamportTime))
	clockSemaphore.Release(1)
}

func AddToQueueCustom(value, address string, lTime sortedset.SCORE) {
	orderedQueueSemaphore.Acquire(context.Background(), 1)
	lamportOrderedQueue.AddOrUpdate(address, lTime, value)
	orderedQueueSemaphore.Release(1)
}

func AddToQueue() string {
	orderedQueueSemaphore.Acquire(context.Background(), 1)
	value := fmt.Sprintf("%v:%v&%v", getAddress(), time.Now().UTC().UnixNano(), lamportTime)
	lamportOrderedQueue.AddOrUpdate(getAddress(), lamportTime, value)
	orderedQueueSemaphore.Release(1)
	return value
}

var CurrentTimestamp string
func AcquireLock() {
	println("acquiring")
	CurrentTimestamp = AddToQueue()
	println("added")
	sendToClients("REQ&" + CurrentTimestamp)
	println("sent")
	IncrementTime()

	for i := 0; i < GetNumberOfClients(); i++ {
		<-permittedEntry
		println("permitted once")
	}
	IncrementTime()

	for lamportOrderedQueue.PeekMin().Value != CurrentTimestamp {
		println("busy waiting")
		<-busyWait
	}
}

func ReleaseLock() {
	address, _, _ := ParseString("TMP&" + CurrentTimestamp)
	lamportOrderedQueue.Remove(address)
	sendToClients("REL&" + CurrentTimestamp)
	IncrementTime()
}
