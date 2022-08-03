package conk

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Run is the entry point of the application.
func Run(intervalDuration time.Duration, onNotifiedCommand []string, onNotNotifiedCommand []string, onTickedCommand []string, dryRun bool, stdinPlaceholder string, stdinDistinct bool) error {
	onNotifiedCmdRunner := makeCommandRunner(onNotifiedCommand, dryRun)
	onNotNotifiedCmdRunner := makeCommandRunner(onNotNotifiedCommand, dryRun)
	onTickedCmdRunner := makeCommandRunner(onTickedCommand, dryRun)

	var m sync.Mutex
	stdinLines := make([]string, 0)
	distinctMap := make(map[string]struct{})

	notifyCh := make(chan interface{}, 1)
	go func() {
		ticker := time.Tick(intervalDuration)

		for range ticker {
			onTickedCmdRunner("")

			select {
			case <-notifyCh:
				if stdinPlaceholder == "" {
					onNotifiedCmdRunner("")
					break
				}

				m.Lock()
				onNotifiedCmdRunner(stdinPlaceholder, stdinLines...)
				stdinLines = stdinLines[:0]
				if stdinDistinct {
					distinctMap = make(map[string]struct{})
				}
				m.Unlock()
			default:
				onNotNotifiedCmdRunner("")
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

		if stdinPlaceholder != "" {
			line := scanner.Text()

			m.Lock()
			if stdinDistinct {
				if _, exists := distinctMap[line]; exists {
					m.Unlock()
					continue
				}
				distinctMap[line] = struct{}{}
			}
			stdinLines = append(stdinLines, line)
			m.Unlock()
		}
	}

	return scanner.Err()
}

func makeCommandRunner(commands []string, dryRun bool) func(stdinPlaceholder string, stdinLines ...string) {
	if len(commands) <= 0 {
		return func(stdinPlaceholder string, stdinLines ...string) {}
	}

	if dryRun {
		return func(stdinPlaceholder string, stdinLines ...string) {
			serializedCommands := `"` + strings.Join(interpolateCommands(commands, stdinPlaceholder, stdinLines...), `" "`) + `"`
			log.Printf("[dry-run] triggered: [%s]", serializedCommands)
		}
	}

	return func(stdinPlaceholder string, stdinLines ...string) {
		cc := interpolateCommands(commands, stdinPlaceholder, stdinLines...)
		cmd := exec.Command(cc[0], cc[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go func() { // don't wait the command finish
			err := cmd.Run()
			if err != nil {
				log.Print(err)
			}
		}()
	}
}

var placeholderNumCache = -1

func interpolateCommands(commands []string, stdinPlaceholder string, stdinLines ...string) []string {
	if stdinPlaceholder == "" {
		return commands
	}

	if placeholderNumCache < 0 { // initial run
		placeholderNumCache = 0
		interpolatedCommands := make([]string, 0, len(commands))
		for _, command := range commands {
			if command == stdinPlaceholder {
				placeholderNumCache++
				interpolatedCommands = append(interpolatedCommands, stdinLines...)
				continue
			}
			interpolatedCommands = append(interpolatedCommands, command)
		}
		return interpolatedCommands
	}

	if placeholderNumCache == 0 {
		return commands
	}

	interpolatedCommands := make([]string, 0, len(commands)-placeholderNumCache+len(stdinLines)*placeholderNumCache)
	for _, command := range commands {
		if command == stdinPlaceholder {
			interpolatedCommands = append(interpolatedCommands, stdinLines...)
			continue
		}
		interpolatedCommands = append(interpolatedCommands, command)
	}
	return interpolatedCommands
}
