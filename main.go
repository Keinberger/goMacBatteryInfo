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

// checkIfClick() checks if a certain menuItem gets clicked, then triggers a specified function
func checkIfClick(menuItem *systray.MenuItem, itemFunction func()) {
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction()
		}
	}
}

// checkIfClickNotify() checks if a certain menuItem gets clicked, then triggers a specified function with two parameters
func checkIfClickNotify(menuItem *systray.MenuItem, itemFunction func(int, int), param ...int) {
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction(param[0], param[1])
			menuItem.Disable()
		}
	}
}

// checkIfClickStop() checks if a certain menuItem gets clicked, then triggers a specified function with one parameter
func checkIfClickStop(menuItem *systray.MenuItem, itemFunction func(int), param ...int) {
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction(param[0])
			menuItem.Disable()
		}
	}
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

// main() executes the systray.Run()
func main() {
	systray.Run(onReady, onExit)
}

// onReady() gets called at beginning of systray.Run() and opens updateBatteryLevel(), checkIfClick()
func onReady() {
	wg.Add(1)
	go updateBatteryLevel(20)
	wg.Wait()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	go checkIfClick(mQuit, systray.Quit)

	fmt.Println(name + " started succesfully")
}

// onExit() gets called when systray finishes
func onExit() {
	shutdown = true
	fmt.Println("Waiting for goroutines to shut down...")
	fmt.Println(name + " quitted succesfully")
}
