package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	notif "github.com/deckarep/gosx-notifier"
	"github.com/getlantern/systray"
)

var notifications = make(map[int]bool)

// pushBatteryNotifyMessage() will trigger notify() when time remaining equals the specified minutesRemaining variable
func pushBatteryNotifyMessage(minutesRemaining, charge int) {
	infoHour, _ := strconv.Atoi(string(title[0]))
	infoMinute, _ := strconv.Atoi(string(title[2:4]))

	minutesTillZero := infoHour*60 + infoMinute

	if charge > 0 && minutesRemaining == 0 {
		message, err := getBatteryInfo()
		if err != nil {
			panic(err)
		}

		percentage, _ := strconv.Atoi(message[61:63])

		remainingChargeTillNotif := percentage - charge
		minutesTillCharge := (minutesTillZero / percentage) * remainingChargeTillNotif

		time.Sleep(time.Duration(minutesTillCharge) * time.Minute)

		chargeString := strconv.Itoa(charge)
		notify("Battery charge is at "+chargeString+"%", "You should consider charging your battery", "")
	} else {
		notifications[minutesRemaining] = false
		wg.Add(1)
		go func() {
			var hour int
			var min int
			var minTillZero int = minutesTillZero

			stop := systray.AddMenuItem("Stop Notifier ("+strconv.Itoa(minutesRemaining)+"min)", "")
			wg.Add(1)
			go checkIfClickStop(stop, stopNotification, minutesRemaining)

			for {
				if checkIfShutdown() {
					break
				}
				if minTillZero > minutesRemaining {
					if notifications[minutesRemaining] {
						notifications[minutesRemaining] = true
						enable(m[minutesRemaining]) // <-- control that
						stop.Hide()
						break
					}
					time.Sleep(10 * time.Second)

					hour, _ = strconv.Atoi(string(title[0]))
					min, _ = strconv.Atoi(string(title[2:4]))

					minTillZero = hour*60 + min
				} else {
					<-stop.ClickedCh
					stop.Hide()
					err := notify("You have "+string(title[0])+"h and "+string(title[2:4])+"min of battery life remaining", "", "")
					if err != nil {
						log.Fatal("There was a problem while sending the notification", err)
					}
					enable(m[minutesRemaining])
					break
				}
			}
			fmt.Println("pushBatteryNotifyMessage has shut down")
			defer wg.Done()
		}()
	}
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
