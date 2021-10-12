package main

import (
	"encoding/csv"
	"os"
	"strconv"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

type Output struct {
	Frame   int
	Events  interface{}
	Players [][]string
}

var fileName = "og-vs-sprout-dust2"

func main() {
	f, err := os.Open("..\\demo\\" + fileName + ".dem")
	checkError(err)
	defer f.Close()

	p := dem.NewParser(f)

	var data []Output

	// parse frame by frame
	for ok := true; ok; ok, err = p.ParseNextFrame() {
		checkError(err)

		gs := p.GameState()
		frame := p.CurrentFrame()

		var players [][]string

		for _, player := range gs.Participants().Playing() {
			players = append(players, extractPlayerData(frame, player))
		}

		o := Output{
			Frame:   frame,
			Players: players,
		}

		data = append(data, o)
	}

	err = csvExport(data)
	checkError(err)
}

func extractPlayerData(frame int, player *common.Player) []string {
	return []string{
		strconv.Itoa(frame),
		player.Name,
		strconv.FormatUint(player.SteamID64, 10),
		strconv.FormatFloat(player.Position().X, 'G', -1, 64),
		strconv.FormatFloat(player.Position().Y, 'G', -1, 64),
		strconv.FormatFloat(player.Position().Z, 'G', -1, 64),

		strconv.FormatFloat(player.LastAlivePosition.X, 'G', -1, 64),
		strconv.FormatFloat(player.LastAlivePosition.Y, 'G', -1, 64),
		strconv.FormatFloat(player.LastAlivePosition.Z, 'G', -1, 64),

		strconv.FormatFloat(player.Velocity().X, 'G', -1, 64),
		strconv.FormatFloat(player.Velocity().Y, 'G', -1, 64),
		strconv.FormatFloat(player.Velocity().Z, 'G', -1, 64),

		strconv.FormatFloat(float64(player.ViewDirectionX()), 'G', -1, 64),
		strconv.FormatFloat(float64(player.ViewDirectionY()), 'G', -1, 64),

		strconv.Itoa(player.Health()),
		strconv.Itoa(player.Armor()),
		strconv.Itoa(player.Money()),
		strconv.Itoa(player.EquipmentValueCurrent()),
		strconv.Itoa(player.EquipmentValueCurrent()),
		strconv.Itoa(player.EquipmentValueRoundStart()),
		strconv.FormatBool(player.IsDucking()),
		strconv.FormatBool(player.HasDefuseKit()),
		strconv.FormatBool(player.HasHelmet()),
		strconv.Itoa(player.Kills()),
		strconv.Itoa(player.Deaths()),
		strconv.Itoa(player.Assists()),
		strconv.Itoa(player.Score()),
		strconv.Itoa(player.MVPs()),
		strconv.Itoa(player.MoneySpentTotal()),
		strconv.Itoa(player.MoneySpentThisRound()),
	}
}

func csvExport(data []Output) error {
	file, err := os.OpenFile(fileName+".csv", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// header
	header := []string{
		"Frame", "Name", "SteamID", "Position_X", "Position_Y", "Position_Z", "LastAlivePosition_X", "LastAlivePosition_Y", "LastAlivePosition_Z",
		"Velocity_X", "Velocity_Y", "Velocity_Z", "ViewDirectionX", "ViewDirectionY", "Hp", "Armor", "Money",
		"CurrentEquipmentValue", "FreezetimeEndEquipmentValue", "RoundStartEquipmentValue", "IsDucking", "HasDefuseKit",
		"HasHelmet", "Kills", "Deaths", "Assists", "Score", "MVPs", "TotalCashSpent", "CashSpentThisRound",
	}
	if err := writer.Write(header); err != nil {
		return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
	}

	// data
	for _, frameData := range data {
		for _, player := range frameData.Players {
			err := writer.Write(player)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
