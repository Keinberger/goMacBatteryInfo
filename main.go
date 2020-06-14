package main

import (
	"fmt"
	"sync"

	"github.com/getlantern/systray"
)

var (
	name     string = "Battery Charge Monitor"
	wg       sync.WaitGroup
	title    string
	m             = make(map[int]*systray.MenuItem)
	shutdown bool = false
)

// returns current value of shutdown
func checkIfShutdown() bool {
	return shutdown
}

func checkIfClick(menuItem *systray.MenuItem, itemFunction func(), param ...int) {
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction()
		}
	}
}

func checkIfClickNotify(menuItem *systray.MenuItem, itemFunction func(int, int), param ...int) {
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction(param[0], param[1])
			menuItem.Disable()
		}
	}
}

func checkIfClickStop(menuItem *systray.MenuItem, itemFunction func(int), param ...int) {
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction(param[0])
			menuItem.Disable()
		}
	}
}

func disable(menuItems ...*systray.MenuItem) {
	for _, v := range menuItems {
		v.Disable()
	}
}

func enable(menuItems ...*systray.MenuItem) {
	for _, v := range menuItems {
		v.Enable()
	}
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	wg.Add(1)
	go updateBatteryLevel(20)
	wg.Wait()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	go checkIfClick(mQuit, systray.Quit)

	fmt.Println(name + " started succesfully")
}

func onExit() {
	shutdown = true
	fmt.Println("Waiting for goroutines to shut down...")
	fmt.Println(name + " quitted succesfully")
}
