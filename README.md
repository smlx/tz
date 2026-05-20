# tz

`tz` is a timezone conversion tool.

## Features

* Location names are fuzzy matched. So e.g. you can write "Zurich" or "Zürich".
    * Only cities with a population of at least 100,000 are included.
* Common timezone shorthands and offsets are accepted (e.g. `UTC`, `UTC+8`, `+0800`, `-07:00`).

Commands have aliases: `c` for `convert` and `m` for `meeting`.

### Convert

Convert from one timezone to another.

Examples:

```bash
# print the current time and timezone in Zürich
tz c zurich
# Output:
# Zürich (Europe/Zurich CEST UTC+2)
# Wed, 06 May 2026 04:44:20 CEST

# print the time and timezone in Poznań when it is 5am local time (the next time it is 5am local time)
tz c poznan @ 5am
# print the time and timezone in Zürich when it is 5am in Sydney (the next time it is 5am in Sydney)
tz c Zurich Sydney 5am
# print the time and timezone in Tokyo when it is 5am on Friday in London (the next time it is 5am on Friday in London)
tz c Tokyo London 5am friday
# print the time and timezone local time when it is 5am in Prague (the next time it is 5am in Prague)
tz c @ prague 5am
```

The logic is:
* the CLI structure is `tz convert [TARGET] [SOURCE] [TIME SPECIFICATION]`
* `@` means "local time"
* if the source location is not provided, default to the local timezone location
* if the time specification is not provided, default to now

### Meeting Planner

Plan meetings by printing a table of times in each location specified compared to UTC, centered around the current time.
The columns of the table are: `UTC`, `Location_1 (UTC+N)`, `Location_2 (UTC+N)`, `Location_3 (UTC+n)` etc.
The current time is shown by a horizontal line.


```bash
tz m zurich london tokyo
# Output:
#                    Zürich         London         Tokyo       
# UTC    Local       Europe/Zurich  Europe/London  Asia/Tokyo  
#        AWST UTC+8  CEST UTC+2     BST UTC+1      JST UTC+9   
# 14:00  22:00       16:00          15:00          23:00       
# 15:00  23:00       17:00          16:00          00:00       
# 16:00  00:00       18:00          17:00          01:00       
# 17:00  01:00       19:00          18:00          02:00       
# 18:00  02:00       20:00          19:00          03:00       
# 19:00  03:00       21:00          20:00          04:00       
# 20:00  04:00       22:00          21:00          05:00       
# 21:00  05:00       23:00          22:00          06:00       
# 22:00  06:00       00:00          23:00          07:00       
# 23:00  07:00       01:00          00:00          08:00       
# 00:00  08:00       02:00          01:00          09:00       
# 01:00  09:00       03:00          02:00          10:00       
# 02:00  10:00       04:00          03:00          11:00       
# 02:45  10:45       04:45          03:45          11:45       <-- current time
# 03:00  11:00       05:00          04:00          12:00       
# 04:00  12:00       06:00          05:00          13:00       
# 05:00  13:00       07:00          06:00          14:00       
# 06:00  14:00       08:00          07:00          15:00       
# 07:00  15:00       09:00          08:00          16:00       
# 08:00  16:00       10:00          09:00          17:00       
# 09:00  17:00       11:00          10:00          18:00       
# 10:00  18:00       12:00          11:00          19:00       
# 11:00  19:00       13:00          12:00          20:00       
# 12:00  20:00       14:00          13:00          21:00       
# 13:00  21:00       15:00          14:00          22:00       
# 14:00  22:00       16:00          15:00          23:00       
```

This feature is inspired by [this web interface](https://www.timeanddate.com/worldclock/meetingtime.html).

## Acknowledgments

The cities timezone data comes from [GeoNames](https://www.geonames.org) and is licensed under the [Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/).
