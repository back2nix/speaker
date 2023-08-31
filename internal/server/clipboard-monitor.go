package server

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func clipboardMonitor(ch chan<- string) {
	cmd := exec.Command("./clipboard-monitor.sh")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Clipboard contents changed:") {
			ch <- strings.TrimPrefix(line, "Clipboard contents changed: ")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
	}

	cmd.Wait()
}
