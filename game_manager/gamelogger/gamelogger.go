package gamelogger

import (
	"fmt"

	scribble "github.com/nanobox-io/golang-scribble"
)

type Game struct {
	name string
}

func Test() {
	db, err := scribble.New("/tmp/db", nil)
	if err != nil {
		fmt.Println(err)
	}
	game := Game{"blekota"}
	err = db.Write("games", "ala", game)
	fmt.Println(err)
}
