package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

// getBatteryInfo returns a string containing the return-value of the pmset -g batt command
func getBatteryInfo() (string, error) {
	if runtime.GOOS == "windows" {
		err := errors.New("Program not executable on Windows")
		return "", err
	} else if runtime.GOOS == "linux" {
		err := errors.New("Program not executable on Linux")
		return "", err
	} else {
		out, err := exec.Command("pmset", "-g", "batt").Output()
		if err != nil {
			return "", err
		}
		return string(out[:]), nil
	}
}

// updateBatteryLevel updates the remaining battery time and the message inside of the application every 30 seconds
func updateBatteryLevel(interval time.Duration) {
	battery := systray.AddMenuItem("Calculating...", "")
	systray.SetTitle("...")
	battery.Disable()

	m[60] = systray.AddMenuItem("Notify (1hour remaining)", "")
	m[30] = systray.AddMenuItem("Notify (30min remaining)", "")
	m[10] = systray.AddMenuItem("Notify (10min remaining)", "")
	wg.Done()

	for k, v := range m {
		disable(v)
		notifications[k] = true
	}

	wg.Add(3)
	go checkIfClickNotify(m[60], pushBatteryNotifyMessage, 60, 0)
	go checkIfClickNotify(m[30], pushBatteryNotifyMessage, 30, 0)
	go checkIfClickNotify(m[10], pushBatteryNotifyMessage, 10, 0)

	var previousLoad string
	wg.Add(1)
	for {
		if checkIfShutdown() {
			break
		}
		time.Sleep(interval * time.Second)

		load, err := getBatteryInfo()
		if err != nil {
			log.Fatal("Error while updating battery", err)
			break
		}
		if load == previousLoad {
			continue
		}

		m1h := m[60]
		m3 := m[30]
		m10 := m[10]
		switch {
		case strings.Contains(load, "no estimate"):
			title = "..."
			battery.SetTitle("Calculating...")
			for _, v := range m {
				if !v.Disabled() {
					disable(v)
				}
			}
			fmt.Println("I was at no 'estimate'")
		case strings.Contains(load, "discharging"):
			title = load[83:89] // [84:89]
			if strings.Contains(title, "r") {
				title = strings.Trim(title, "r")
			}
			if strings.Contains(title, ";") {
				title = strings.Trim(title, ";")
			}
			title = strings.TrimSpace(title)
			battery.SetTitle(title + " remaining")
			for k, v := range m {
				if v.Disabled() && notifications[k] {
					enable(v)
				}
			}
			h, _ := strconv.Atoi(string(title[0]))
			m, _ := strconv.Atoi(string(title[2:4]))
			if h <= 1 {
				if m == 0 || h < 1 {
					disable(m1h)
				}
				if h < 1 {
					if m <= 30 {
						disable(m3)
					}
					if m <= 10 {
						disable(m10)
					}
				}
			}
		case strings.Contains(load, "AC attached") || strings.Contains(load, "100%"):
			title = "∞"
			battery.SetTitle("Battery is charged")
			for _, v := range m {
				if !v.Disabled() {
					disable(v)
				}
			}
		case strings.Contains(load, "charging"):
			title = load[75:81] // [75:81]
			if strings.Contains(title, "r") {
				title = strings.Trim(title, "r")
			}
			if strings.Contains(title, ";") {
				title = strings.Trim(title, ";")
			}
			title = strings.TrimSpace(title)
			battery.SetTitle(title + " until charged")
			for _, v := range m {
				if !v.Disabled() {
					disable(v)
				}
			}
			fmt.Println("I was at no 'charging'")
		}

		systray.SetTitle(title)
		previousLoad = load

		fmt.Println("I was at the end of it")
	}
	fmt.Println("Update battery level has shutdown")
	defer wg.Done()
}
