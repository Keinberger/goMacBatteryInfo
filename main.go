package main

import (
	"fmt"
	"sync"

	"github.com/getlantern/systray"
)

var (
	name     string = "Battery Charge Monitor"
	wg       sync.WaitGroup
	battery  *systray.MenuItem
	title    string
	m        = make(map[int]*systray.MenuItem)
	shutdown bool
)

// returns current value of shutdown
func checkIfShutdown() bool {
	return shutdown
}

// checkIfClick() checks if a certain menuItem gets clicked, then triggers a specified function
func checkIfClick(menuItem *systray.MenuItem, itemFunction func()) {
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

// checkIfClickNotify() checks if a certain menuItem gets clicked, then triggers a specified function with two parameters
func checkIfClickNotify(menuItem *systray.MenuItem, itemFunction func(int, int), param ...int) {
Y:
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction(param[0], param[1])
			menuItem.Disable()
		default:
			if checkIfShutdown() {
				break Y
			}
		}
	}
	defer wg.Done()
}

// checkIfClickStop() checks if a certain menuItem gets clicked, then triggers a specified function with one parameter
func checkIfClickStop(menuItem *systray.MenuItem, itemFunction func(int), param ...int) {
Y:
	for {
		select {
		case <-menuItem.ClickedCh:
			itemFunction(param[0])
			menuItem.Disable()
		default:
			if checkIfShutdown() {
				break Y
			}
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

// main() executes the systray.Run()
func main() {
	systray.Run(onReady, onExit)
}

// onReady() gets called at beginning of systray.Run() and opens updateBatteryLevel(), checkIfClick()
func onReady() {
	battery = systray.AddMenuItem("Calculating...", "")
	systray.SetTitle("...")
	battery.Disable()

	m[60] = systray.AddMenuItem("Notify (1hour remaining)", "")
	m[30] = systray.AddMenuItem("Notify (30min remaining)", "")
	m[10] = systray.AddMenuItem("Notify (10min remaining)", "")

	go updateBatteryLevel(20)

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	wg.Add(1)
	go checkIfClick(mQuit, systray.Quit)

	fmt.Println(name + " started succesfully")
}

// onExit() gets called when systray finishes
func onExit() {
	fmt.Println("Waiting for goroutines to shut down...")
	shutdown = true
	wg.Wait()
	fmt.Println(name + " quitted succesfully")
}
