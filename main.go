package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/zyedidia/generic/list"
)

var shouldBenchmark bool
var displayToTerminal bool
var rollCount int

func rollRaffle(raffleData *list.List[RaffleEntry], rolls int) map[string]*RolledEntry {
	raffleWinners := make(map[string]*RolledEntry)
	highRoll := 0
	for entry := raffleData.Front; entry != nil; entry = entry.Next {
		highRoll += entry.Value.Points
	}

	for i := 0; i < rolls; i++ {
		if highRoll == 0 {
			break
		}
		var raffleRoll int
		for j := 0; j < 5; j++ {
			raffleRoll = rand.Intn(highRoll) + 1
		}
		counter := raffleRoll
		var raffleWinner string

		for entry := raffleData.Front; entry != nil; entry = entry.Next {
			counter -= entry.Value.Points
			if counter <= 0 {
				raffleWinner = entry.Value.Patron
				entry.Value.Points -= 1
				if entry.Value.Points <= 0 {
					raffleData.Remove(entry)
				}
				break
			}
		}

		_, ok := raffleWinners[raffleWinner]
		if !ok {
			raffleWinners[raffleWinner] = &RolledEntry{points: 0, rolls: make([]int, 0, rolls)}
		}

		winner := raffleWinners[raffleWinner]
		winner.rolls = append(winner.rolls, raffleRoll)
		winner.points += 1

		highRoll -= 1
	}

	return raffleWinners
}

func timeIt(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

func init() {
	flag.BoolVar(&displayToTerminal, "display", false, "Display results to terminal")
	flag.BoolVar(&shouldBenchmark, "benchmark", false, "Benchmark the program")
	flag.IntVar(&rollCount, "rolls", 5, "Number of rolls to perform")
	flag.Parse()
}

func main() {
	raffleData := importDataFile()
	raffleResults := rollRaffle(raffleData, rollCount)

	printToFile(raffleResults)
	exportDataFile(raffleData)

	if displayToTerminal {
		for patron, rolledEntry := range raffleResults {
			fmt.Printf("+%d for %s %s \n", rolledEntry.points, patron, formatArray(rolledEntry.rolls))
		}
	}
}
