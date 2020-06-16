package main

import (
	"log"
	"strconv"
	"time"

	notif "github.com/deckarep/gosx-notifier"
	"github.com/getlantern/systray"
)

var notifications = make(map[int]bool)

// pushBatteryNotifyMessage() will trigger notify() when time remaining equals the specified minutesRemaining variable
func pushBatteryNotifyMessage(minutesRemaining int) {
	infoHour, _ := strconv.Atoi(string(title[0]))
	infoMinute, _ := strconv.Atoi(string(title[2:4]))

	minutesTillZero := infoHour*60 + infoMinute

	notifications[minutesRemaining] = false
	wg.Add(1)
	go func() {
		var hour int
		var min int
		var minTillZero int = minutesTillZero

		stop := systray.AddMenuItem("Stop Notifier ("+strconv.Itoa(minutesRemaining)+"min)", "")
		wg.Add(1)
		go checkIfClick(stop, stopNotification, minutesRemaining)

		for {
			if checkIfShutdown() {
				stop.ClickedCh <- struct{}{}
				break
			}
			if minTillZero > minutesRemaining {
				if notifications[minutesRemaining] {
					enable(m[minutesRemaining]) // <-- control that
					stop.Hide()
					break
				}
				time.Sleep(10 * time.Second)

				hour, _ = strconv.Atoi(string(title[0]))
				min, _ = strconv.Atoi(string(title[2:4]))

				minTillZero = hour*60 + min
			} else {
				stop.ClickedCh <- struct{}{}
				stop.Hide()
				err := notify("You have "+string(title[0])+"h and "+string(title[2:4])+"min of battery life remaining", "", "")
				if err != nil {
					log.Fatal("There was a problem while sending the notification", err)
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

	note.Group = "com.philippkeinberger.macBatteryInfo.batteryLevelNotification"
	note.Sender = "com.philippkeinberger.macBatteryInfo"
	note.Title = tit
	note.ContentImage = iconPath
	note.Sound = "'default'"
	note.AppIcon = "icon/battery.png"

	err := note.Push()

	return err
}

// stopNotification() changes the notifications struct so that pushBatteryNotification() will break out of the for loop
func stopNotification(minutesRemaining int) {
	notifications[minutesRemaining] = true
}
