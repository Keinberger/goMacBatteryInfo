package main

import (
	"errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

type batteryInfo struct {
	calculating   bool
	charging      bool
	fullyCharged  bool
	timeOnBattery timeSpec
	timeRemaining timeSpec
}

type timeSpec struct {
	hours int
	mins  int
}

func convMinToSpec(min int) timeSpec {
	return timeSpec{min / 60, min % 60}
}

func convTimeSpecToMin(t timeSpec) int {
	return t.hours*60 + t.mins
}

func getTitle(t timeSpec) (title string) {
	title = strconv.Itoa(t.hours) + ":"
	if t.mins < 10 {
		title += "0"
	}
	title += strconv.Itoa(t.mins)
	return
}

// getBatteryInfo returns a string containing the return-value of the 'pmset -g batt' command
func getBatteryInfo() (info batteryInfo, err error) {
	switch runtime.GOOS {
	case "windows":
		err = errors.New("Program not executable on Windows")
	case "linux":
		// var batteryPath, entireString []byte
		// batteryPath, err = exec.Command("upower", "-e", "--enumerate").Output() // not sure if --enumerate is necessary
		// if err != nil {
		// 	return
		// }
		// entireString, err = exec.Command("upower", "-i", string(batteryPath), "grep", "-E", `"state|to\ full|percentage"`)
		// if err != nil {
		// 	return
		// }

		// sampleOutput := `
		// state:               charging
		// time to full:        57.3 minutes
		// percentage:          42.5469%
		// `

		// entireFormatted := strings.Split(sampleOutput, ":")

		// info.calculating = false
		// info.charging = entireFormatted[1] == "charging"
		// info.fullyCharged = entireFormatted[1] == "charged"

		// if !info.calculating && !info.fullyCharged {
		// }

		err = errors.New("Program not executable on Linux")
	default:
		var out []byte
		out, err = exec.Command("pmset", "-g", "batt").Output()
		if err != nil {
			return
		}
		entireString := strings.Split(strings.Join(strings.Split(strings.Join(strings.Split(string(out[:]), "\n"), ""), " "), ""), "	")
		if len(entireString) < 2 {
			err = errors.New("Could not retrieve battery info")
			break
		}
		entireFormatted := strings.Split(entireString[1], ";")

		info.calculating = strings.Contains(entireFormatted[2], "no estimate")
		info.charging = entireFormatted[1] == "charging"
		info.fullyCharged = entireFormatted[1] == "charged"

		if !info.calculating && !info.fullyCharged {
			remaining := strings.Split(entireFormatted[2][:4], ":")
			if len(remaining) < 2 {
				info.calculating = true
				break
			}
			info.timeRemaining.hours, _ = strconv.Atoi(remaining[0])
			info.timeRemaining.mins, _ = strconv.Atoi(remaining[1])
		}
	}

	return info, nil
}

// updateBatteryLevel updates the remaining battery time and the message inside of the application every 30 seconds
func updateBatteryLevel() {
	for k, v := range conf.Reminders {
		disable(v.item)
		conf.Reminders[k].notifier = true
	}

	var previousInfo batteryInfo
	for {
		time.Sleep(time.Duration(conf.UpdateInterval) * time.Second)
		if checkIfShutdown() {
			break
		}

		batteryInfo, err := getBatteryInfo()
		if err != nil {
			logError("Error while updating battery info", err)
			break
		}
		if batteryInfo == previousInfo {
			continue
		}
		switch {
		case batteryInfo.calculating:
			title = "..."
			battery.SetTitle("Calculating...")
			for _, v := range conf.Reminders {
				if !v.item.Disabled() {
					disable(v.item)
				}
			}
		case batteryInfo.fullyCharged:
			title = "∞"
			battery.SetTitle("Battery is charged")
			for _, v := range conf.Reminders {
				if !v.item.Disabled() {
					disable(v.item)
				}
			}
		case batteryInfo.charging: // charging battery
			title = getTitle(batteryInfo.timeRemaining)
			battery.SetTitle(title + " until charged")
			for _, v := range conf.Reminders {
				if !v.item.Disabled() {
					disable(v.item)
				}
			}
		default: // discharging battery (bc none of the other cases fit)
			title = getTitle(batteryInfo.timeRemaining)
			battery.SetTitle(title + " remaining")
			for _, v := range conf.Reminders {
				if convTimeSpecToMin(batteryInfo.timeRemaining) < v.MinutesRemaining {
					disable(v.item)
				} else if v.item.Disabled() && v.notifier {
					enable(v.item)
				}
			}
		}
		systray.SetTitle(title)
		previousInfo = batteryInfo
	}
	defer wg.Done()
}
