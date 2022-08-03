package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/moznion/conk"
	"github.com/moznion/conk/internal"
)

func main() {
	var intervalDurationSec uint64
	var onNotifiedCommand string
	var onNotNotifiedCommand string
	var onTickedCommand string
	var dryRun bool
	var shouldShowVersion bool
	var stdinPlaceholder string
	var stdinDistinct bool

	flag.Uint64Var(&intervalDurationSec, "interval-sec", 0, "[mandatory] interval duration seconds to check the bytes that come from stdin.")
	flag.StringVar(&onNotifiedCommand, "on-notified-cmd", "[]", "[semi-mandatory] command that runs on notified (i.e. when bytes come from stdin in an interval). this must be JSON string array. it requires this value and/or \"--on-not-notified-cmd\"")
	flag.StringVar(&onNotNotifiedCommand, "on-not-notified-cmd", "[]", "[semi-mandatory] command that runs on NOT notified (i.e. when bytes don't come from stdin in an interval). this must be JSON string array. it requires this value and/or \"--on-notified-cmd\"")
	flag.StringVar(&onTickedCommand, "on-ticked-cmd", "[]", "command that runs every interval. this must be JSON string array.")
	flag.BoolVar(&dryRun, "dry-run", false, "dry-run mode. if this value is true, it notifies the command was triggered, instead of executing commands.")
	flag.BoolVar(&shouldShowVersion, "version", false, "show version info")
	flag.StringVar(&stdinPlaceholder, "stdin-placeholder", "", "placeholder name that can be used in `on-notified-cmd` to give the command the arguments that come from STDIN.")
	flag.BoolVar(&stdinDistinct, "stdin-distinct", false, "if this value is true, it makes the arguments for `on-notified-cmd` that come from STDIN distinct (i.e. makes them unique). see also: -stdin-placeholder")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `A command-line tool that triggers command when the input (doesn't) comes from STDIN in an interval.
If the input comes from STDIN, it fires "--on-notified-cmd" command. Elsewise, it executes "--on-not-notified-cmd".

Usage of %s:
   %s [OPTIONS]
Options
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if shouldShowVersion {
		v, _ := json.Marshal(map[string]string{
			"revision":  internal.Revision,
			"version":   internal.Version,
			"goVersion": runtime.Version(),
		})
		fmt.Printf("%s\n", v)
		os.Exit(0)
	}

	if intervalDurationSec <= 0 {
		log.Fatal(`"--interval-sec" is mandatory parameter`)
	}

	var onNotifiedCmd []string
	err := json.Unmarshal([]byte(onNotifiedCommand), &onNotifiedCmd)
	if err != nil {
		log.Fatal(err)
	}

	var onNotNotifiedCmd []string
	err = json.Unmarshal([]byte(onNotNotifiedCommand), &onNotNotifiedCmd)
	if err != nil {
		log.Fatal(err)
	}

	if len(onNotifiedCmd) <= 0 && len(onNotNotifiedCmd) <= 0 {
		log.Fatal(`either one of "--on-notified-cmd" or "--on-not-notified-cmd" is mandatory`)
	}

	var onTickedCmd []string
	err = json.Unmarshal([]byte(onTickedCommand), &onTickedCmd)
	if err != nil {
		log.Fatal(err)
	}

	err = conk.Run(time.Duration(intervalDurationSec)*time.Second, onNotifiedCmd, onNotNotifiedCmd, onTickedCmd, dryRun, stdinPlaceholder, stdinDistinct)
	if err != nil {
		log.Fatal(err)
	}
}
