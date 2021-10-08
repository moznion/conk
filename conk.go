package conk

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Run is the entry point of the application.
func Run(intervalDuration time.Duration, onNotifiedCommand []string, onNotNotifiedCommand []string, onTickedCommand []string, dryRun bool) error {
	onNotifiedCmdRunner := makeCommandRunner(onNotifiedCommand, dryRun)
	onNotNotifiedCmdRunner := makeCommandRunner(onNotNotifiedCommand, dryRun)
	onTickedCmdRunner := makeCommandRunner(onTickedCommand, dryRun)

	notifyCh := make(chan interface{}, 1)
	go func() {
		ticker := time.Tick(intervalDuration)

		for range ticker {
			onTickedCmdRunner()

			select {
			case <-notifyCh:
				onNotifiedCmdRunner()
			default:
				onNotNotifiedCmdRunner()
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		select {
		case notifyCh <- struct{}{}:
		default:
		}
	}

	return scanner.Err()
}

func makeCommandRunner(commands []string, dryRun bool) func() {
	if len(commands) <= 0 {
		return func() {}
	}

	if dryRun {
		serializedCommands := strings.Join(commands, " ")
		return func() {
			log.Printf("[dry-run] triggered: `%s`", serializedCommands)
		}
	}

	return func() {
		cmd := exec.Command(commands[0], commands[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start() // don't wait the command finish
		if err != nil {
			log.Print(err)
		}
	}
}
