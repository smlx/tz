# tz

`tz` is a timezone conversion tool.

## Features

Location names are fuzzy matched. So e.g. you can write "Zurich" or "Zürich".

### Convert

Convert from one timezone to another.

Examples:

```bash
# print the current time and timezone in Zürich
tz convert Zurich
# print the time and timezone in Zürich when it is 5am local time (the next time it is 5am local time)
tz convert Zurich @ 5am
# print the time and timezone in Zürich when it is 5am in Sydney (the next time it is 5am in Sydney)
tz convert Zurich Sydney 5am
# print the time and timezone in Tokyo when it is 5am on Friday in Bangalore (the next time it is 5am on Friday in Bangalore)
tz convert Tokyo Bangalore 5am friday
# print the time and timezone local time when it is 5am in Zürich (the next time it is 5am in Zürich)
tz convert @ Zurich 5am
```

The logic is:
* `@` means "local time"
* the CLI structure is `tz convert [TARGET] [SOURCE] [TIME SPECIFICATION]
* if the source location is not provided, default to the local timezone location
* if the time specification is not provided, default to now

### Meeting Planner

Plan meetings by printing a table of times in each location specified compared to UTC, centered around the current time.
The columns of the table are: `UTC`, `Location_1 (UTC+N)`, `Location_2 (UTC+N)`, `Location_3 (UTC+n)` etc.
The current time is shown by a horizontal line.


```bash
tz meeting Zürich Bangalore Tokyo
```

This feature is inspired by [this web interface](https://www.timeanddate.com/worldclock/meetingtime.html).

## Acknowledgments

The cities timezone data comes from [GeoNames](https://www.geonames.org) and is licensed under the [Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/).
