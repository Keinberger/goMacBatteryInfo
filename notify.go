package main

import (
	"os"
	"strconv"
	"time"

	notif "github.com/deckarep/gosx-notifier"
	"github.com/getlantern/systray"
)

// checkIfExists checks if a file/folder exists
func checkIfExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// pushBatteryNotifyMessage() will trigger notify() when time remaining equals the specified minutesRemaining variable
func pushBatteryNotifyMessage(notifier *reminder) {
	notifier.notifier = false
	wg.Add(1)
	go func() {
		stop := systray.AddMenuItem("Stop Notifier (at "+getTitle(convMinToSpec(notifier.MinutesRemaining))+")", "")

		wg.Add(1)
		go checkIfClick(stop, stopNotification, notifier)
		for {
			info, err := getBatteryInfo()
			if err != nil {
				logError("", err)
			}
			minTillZero := convTimeSpecToMin(info.timeRemaining)
			if checkIfShutdown() {
				stop.ClickedCh <- struct{}{}
				break
			}
			if minTillZero > notifier.MinutesRemaining {
				if notifier.notifier {
					stop.Hide()
					break
				}
				time.Sleep(time.Duration(conf.UpdateInterval) * time.Second)
				minTillZero = convTimeSpecToMin(info.timeRemaining)
			} else {
				stop.ClickedCh <- struct{}{}
				stop.Hide()
				inf, err := getBatteryInfo()
				if !logError("", err) {
					message := "You have " + strconv.Itoa(inf.timeRemaining.hours) + "h and " + strconv.Itoa(inf.timeRemaining.mins) + "min of battery life remaining"
					err = notify(message, "", "")
					logError("There was a problem while sending the notification", err)
				}
				break
			}
		}
		defer wg.Done()
	}()
}

// notify() sends a message
func notify(msg, tit, iconPath string) error {
	note := notif.NewNotification(msg)

	note.Title = tit
	note.ContentImage = iconPath
	note.Sound = "'default'"

	if checkIfExists(conf.AppIcon) {
		note.AppIcon = conf.AppIcon
	}

	return note.Push()
}

// stopNotification() changes the notifications struct so that pushBatteryNotification() will break out of the for loop
func stopNotification(notifier *reminder) {
	for k, v := range conf.Reminders {
		if v.MinutesRemaining == notifier.MinutesRemaining {
			enable(conf.Reminders[k].item)
			conf.Reminders[k].notifier = true
		}
	}
}
