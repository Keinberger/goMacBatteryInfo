package main

import (
	"log"
	"strconv"
	"time"

	notif "github.com/deckarep/gosx-notifier"
	"github.com/getlantern/systray"
)

// pushBatteryNotifyMessage() will trigger notify() when time remaining equals the specified minutesRemaining variable
func pushBatteryNotifyMessage(notifier reminder) {
	info, err := getBatteryInfo()
	if err != nil {
		log.Fatal(err)
	}

	minutesTillZero := convTimeSpecToMin(info.timeRemaining)
	notifier.notifier = false
	wg.Add(1)
	go func() {
		var minTillZero int = minutesTillZero
		stop := systray.AddMenuItem("Stop Notifier (at "+getTitle(convMinToSpec(notifier.MinutesRemaining))+")", "")

		wg.Add(1)
		go checkIfClick(stop, stopNotification, notifier)
		for {
			if checkIfShutdown() {
				stop.ClickedCh <- struct{}{}
				break
			}
			if minTillZero > notifier.MinutesRemaining {
				if notifier.notifier {
					for k, v := range conf.Reminders {
						if v == notifier {
							enable(conf.Reminders[k].item)
						}
					}
					stop.Hide()
					break
				}
				time.Sleep(10 * time.Second)

				minTillZero = convTimeSpecToMin(info.timeRemaining)
			} else {
				stop.ClickedCh <- struct{}{}
				stop.Hide()
				inf, err := getBatteryInfo()
				logError("", err)
				message := "You have " + strconv.Itoa(inf.timeOnBattery.hours) + "h and " + strconv.Itoa(inf.timeOnBattery.mins) + "min of battery life remaining"
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
	note.AppIcon = "icon/battery.png"

	err := note.Push()

	return err
}

// stopNotification() changes the notifications struct so that pushBatteryNotification() will break out of the for loop
func stopNotification(notifier reminder) {
	for k, v := range conf.Reminders {
		if v == notifier {
			conf.Reminders[k].notifier = true
		}
	}
}
