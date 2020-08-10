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
	Y:
		for {
			for i := 0; i < conf.UpdateInterval*1000; i++ {
				if checkIfShutdown() {
					stop.ClickedCh <- struct{}{}
					break Y
				}
				time.Sleep(1 * time.Millisecond)
			}

			info, err := getBatteryInfo()
			logError("", err)
			minTillZero := convTimeSpecToMin(info.timeRemaining)
			if minTillZero > notifier.MinutesRemaining {
				if notifier.notifier {
					stop.Hide()
					break
				}
			} else {
				stop.ClickedCh <- struct{}{}
				stop.Hide()
				message := "You have " + strconv.Itoa(info.timeRemaining.hours) + "h and " + strconv.Itoa(info.timeRemaining.mins) + "min of battery life remaining"
				err = notify(message, "", "")
				logError("There was a problem while sending the notification", err)
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
