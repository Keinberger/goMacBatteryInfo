# GoMacBatteryInfo

## Requirements

The program requires a recent version of MacOS and Go11 or newer.

### Configuration

Custom configurations may be made in the config.json file.

The configuration file must have the following syntax:
```json
{
  "updateInterval": 20,
  "reminders": [
    {
      "min": 90
    },
    {
      "min": 60
    },
    {
      "min": 30
    }
  ]
}
```

* `UpdateInterval` specifies the interval in second, in which the program should update the time remaining/untill full.
* `reminders` is an array of reminders, each being a button visible inside of the aplicaiton.
** `min` specifies the minutes of the reminder (e.g. 60 minutes/90 minutes) remaining, when one wants to be reminded.

The application will remind you, if the battery time drops below the mintues remaining specified in the configuration file.
One may parse the configuration file with the `-config <filePath>` flag when starting the program. If no configuration file is specified, the programm will use the default values shown in the example of the json file structure above. 

### Usage

Run `go run *.go` or `./macBatteryInfo` and you will see `...` popping up in your menu bar. After 20sec you should see the remaining time to charge, or to use your Mac.

You can set a reminder for 1h, 30min or 10min (or as you specified in the config file) of remaining battery life. After you set that reminder, you have the option to stop it.

In order to quit the app, click on it and select `Quit`. Quitting the program can take up to as many seconds as specified as the UpdateInterval.

### Important

There might still be some bugs in the program.
