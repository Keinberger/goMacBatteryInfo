package main

import (
	"errors"
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
	for k, v := range m {
		disable(v)
		notifications[k] = true
	}

	var previousLoad string
	for {
		time.Sleep(interval * time.Second)
		if checkIfShutdown() {
			break
		}

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
		case strings.Contains(load, "no estimate") || strings.Contains(load, "AC attached") && !strings.Contains(load, "100%"):
			title = "..."
			battery.SetTitle("Calculating...")
			for _, v := range m {
				if !v.Disabled() {
					disable(v)
				}
			}
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
		case strings.Contains(load, "100%"):
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
		}

		systray.SetTitle(title)
		previousLoad = load
	}
	defer wg.Done()
}
