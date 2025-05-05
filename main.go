package main

import (
	"github.com/Killazius/yadro/config"
	"github.com/Killazius/yadro/event"
	"github.com/Killazius/yadro/proc"
	"os"
)

// eventPath задает путь к файлу с событиями гонки
const eventPath = "sunny_5_skiers/events"

func main() {
	cfg := config.MustLoad()
	f, err := os.OpenFile(eventPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	events, err := event.Load(f)
	if err != nil {
		panic(err)
	}
	p := proc.New(cfg, os.Stdout)
	p.ProcessEvents(events)

}
