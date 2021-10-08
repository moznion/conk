# conk

A command-line tool that triggers command when the input (doesn't) comes from STDIN in an interval.

## Usage

```
$ conk -h
A command-line tool that triggers command when the input (doesn't) comes from STDIN in an interval.
If the input comes from STDIN, it fires "--on-notified-cmd" command. Elsewise, it executes "--on-not-notified-cmd".

Usage of ./conk:
   ./conk [OPTIONS]
Options
  -dry-run
        dry-run mode. if this value is true, it notifies the command was triggered, instead of executing commands.
  -interval-sec uint
        [mandatory] interval duration seconds to check the bytes that come from stdin.
  -on-not-notified-cmd string
        [semi-mandatory] command that runs on NOT notified (i.e. when bytes don't come from stdin in an interval). this must be JSON string array. it requires this value and/or "--on-notified-cmd" (default "[]")
  -on-notified-cmd string
        [semi-mandatory] command that runs on notified (i.e. when bytes come from stdin in an interval). this must be JSON string array. it requires this value and/or "--on-not-notified-cmd" (default "[]")
  -on-ticked-cmd string
        command that runs every interval. this must be JSON string array. (default "[]")
  -version
        show version info
```

https://user-images.githubusercontent.com/1422834/136492005-e9816823-79ae-4cdb-9282-abb205039bcb.mov

## Releases

https://github.com/moznion/conk/releases

## Author

moznion (<moznion@mail.moznion.net>)

