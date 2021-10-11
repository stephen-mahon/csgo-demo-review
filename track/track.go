package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/golang/geo/r3"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

// https://gist.github.com/markus-wa/6aad5ce8d6ef34b5a3f340db5547c74c

func main() {
	f, err := os.Open("..\\demo\\og-vs-sprout-dust2.dem")
	checkError(err)
	defer f.Close()

	p := dem.NewParser(f)

	// rounds -> time-since-plant -> player-name -> position
	postPlantPositions := make(map[int]map[time.Duration]map[string]r3.Vector)

	var currentPPPos map[time.Duration]map[string]r3.Vector
	var lastPlantTime time.Duration
	var lastSnapshotTime time.Duration
	var isBombPlantActive bool

	// snapshot when the bomb gets planted
	p.RegisterEventHandler(func(e events.BombPlanted) {
		isBombPlantActive = true
		lastPlantTime = p.CurrentTime()
		lastSnapshotTime = p.CurrentTime()

		currentPPPos = make(map[time.Duration]map[string]r3.Vector)
		currentPPPos[0] = positionSnapshot(p)
	})

	// snapshot positions every 5 seconds
	p.RegisterEventHandler(func(e events.FrameDone) {
		const snapshotFrequency = 5 * time.Second

		now := p.CurrentTime()
		if isBombPlantActive && (lastSnapshotTime+snapshotFrequency) < now {
			lastSnapshotTime = now
			currentPPPos[now-lastPlantTime] = positionSnapshot(p)
		}
	})

	// store post-plant positions at the end of the round
	p.RegisterEventHandler(func(e events.RoundEnd) {
		if !isBombPlantActive {
			return
		}

		isBombPlantActive = false
		postPlantPositions[p.GameState().TotalRoundsPlayed()] = currentPPPos
	})

	// Parse to end
	err = p.ParseToEnd()
	checkError(err)

	// just to make sure, maybe we didn't get a RoundEnd event for the final round
	if isBombPlantActive {
		postPlantPositions[p.GameState().TotalRoundsPlayed()] = currentPPPos
	}

	// sort rounds, otherwise output order is random
	rounds := make([]int, 0)
	for k := range postPlantPositions {
		rounds = append(rounds, k)
	}
	sort.Ints(rounds)

	for _, roundNr := range rounds {
		ppos := postPlantPositions[roundNr]

		fmt.Printf("Round %d:\n", roundNr)

		snapshotTimes := make([]int, 0)
		for t := range ppos {
			snapshotTimes = append(snapshotTimes, int(t))
		}
		sort.Ints(snapshotTimes)

		for _, t := range snapshotTimes {
			timeSincePlant := time.Duration(t)
			positions := ppos[timeSincePlant]

			fmt.Printf(" t=%s:\n", timeSincePlant)
			for name, position := range positions {
				fmt.Printf("  %-30s @ %5.0f %5.0f %5.0f\n", name, position.X, position.Y, position.Z)
			}
		}
	}
}

func positionSnapshot(parser dem.Parser) map[string]r3.Vector {
	snapshot := make(map[string]r3.Vector)
	for _, pl := range parser.GameState().Participants().Playing() {
		if !pl.IsAlive() {
			continue
		}

		snapshot[pl.Name] = pl.Position()
	}

	return snapshot
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
