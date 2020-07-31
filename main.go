package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/getlantern/systray"
)

var (
	name     string = "Battery Charge Monitor"
	wg       sync.WaitGroup
	battery  *systray.MenuItem
	title    string
	conf     config
	shutdown bool
)

type config struct {
	UpdateInterval int        `json:"updateInterval"`
	Reminders      []reminder `json:"reminders"`
}

type reminder struct {
	MinutesRemaining int `json:"min"`
	item             *systray.MenuItem
	notifier         bool
}

// checkIfShutdown() returns current value of shutdown
func checkIfShutdown() bool {
	return shutdown
}

// checkIfClickQuit() checks if a certain menuItem gets clicked, then triggers a specified function
func checkIfClickQuit(menuItem *systray.MenuItem, itemFunction func()) {
Y:
	for {
		select {
		case <-menuItem.ClickedCh:
			break Y
		}
	}
	wg.Done()
	defer itemFunction()
}

// checkIfClick() checks if a certain menuItem gets clicked, then triggers a specified function with one parameter
func checkIfClick(menuItem *systray.MenuItem, itemFunction func(reminder), param reminder) {
Y:
	for {
		select {
		case <-menuItem.ClickedCh:
			if checkIfShutdown() {
				break Y
			}
			itemFunction(param)
			menuItem.Disable()
		}
	}
	defer wg.Done()
}

// disable() disables an array of *systray.menuItem
func disable(menuItems ...*systray.MenuItem) {
	for _, v := range menuItems {
		v.Disable()
	}
}

// enable() enables an array of *systray.menuItem
func enable(menuItems ...*systray.MenuItem) {
	for _, v := range menuItems {
		v.Enable()
	}
}

func openConfig(filePath string) config {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	con := config{}
	err = json.Unmarshal(content, &con)
	if err != nil {
		panic(err)
	}

	return con
}

// main() executes the systray.Run()
func main() {
	configFlag := flag.String("config", "config.json", "Path to config file (JSON)")
	flag.Parse()
	conf = openConfig(*configFlag)

	systray.Run(onReady, onExit)
}

// onReady() gets called at beginning of systray.Run() and opens updateBatteryLevel(), checkIfClick()
func onReady() {
	battery = systray.AddMenuItem("Calculating...", "")
	systray.SetTitle("...")
	battery.Disable()

	for k, v := range conf.Reminders {
		conf.Reminders[k].item = systray.AddMenuItem("Notify ("+getTitle(convMinToSpec(v.MinutesRemaining))+" remaining)", "")
		wg.Add(1)
		go checkIfClick(conf.Reminders[k].item, pushBatteryNotifyMessage, v)
	}

	wg.Add(1)
	go updateBatteryLevel()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	wg.Add(1)
	go checkIfClickQuit(mQuit, systray.Quit)

	fmt.Println(name + " started succesfully")
}

// onExit() gets called when systray finishes
func onExit() {
	fmt.Println("Waiting for goroutines to shut down...")
	shutdown = true
	for _, v := range conf.Reminders {
		v.item.ClickedCh <- struct{}{}
	}
	wg.Wait()
	fmt.Println(name + " quitted succesfully")
}
