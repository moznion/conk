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
  -stdin-distinct on-notified-cmd
        if this value is true, it makes the arguments for on-notified-cmd that come from STDIN distinct (i.e. makes them unique). see also: -stdin-placeholder
  -stdin-placeholder on-notified-cmd
        placeholder name that can be used in on-notified-cmd to give the command the arguments that come from STDIN.
  -version
        show version info
```

https://user-images.githubusercontent.com/1422834/136492005-e9816823-79ae-4cdb-9282-abb205039bcb.mov

### -stdin-placeholder

`-stdin-placeholder` option receives a placeholder name that can be used in `on-notified-cmd` to give the command the arguments that come from STDIN.

For example,

```
$ generator | conk -interval-sec 3 -on-notified-cmd '["echo", "{{__STDIN__}}"]' -stdin-placeholder '{{__STDIN__}}'
```

if the `generator` outputs `foo\nbar\nbuz\n` to STDOUT in 3 seconds, the conk command interpolates those texts in `{{__STDIN__}}`, like `["echo", "foo", "bar", "buz"].

### -stdin-distinct

`-stdin-distinct` option accepts the flag to instruct whether it makes the arguments for `on-notified-cmd` that come from STDIN distinct, i.e. makes them unique.

```
$ generator | conk -interval-sec 3 -on-notified-cmd '["echo", "{{__STDIN__}}"]' -stdin-placeholder '{{__STDIN__}}' -stdin-distinct
```

in these commands, if the `generator` outputs `foo\nbar\nfoo\nfoo` to STDOUT in 3 seconds, the conk command interpolates those texts in `{{__STDIN__}}` with unique arguments, like `["echo", "foo", "bar"] because `foo` is duplicated.

## Releases

https://github.com/moznion/conk/releases

## Author

moznion (<moznion@mail.moznion.net>)

