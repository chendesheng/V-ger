package util

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func KillProcess(path string) {
	ps := exec.Command("ps", "-e", "-opid,comm")
	output, _ := ps.Output()

	for i, s := range strings.Split(string(output), "\n") {
		if i == 0 || len(s) == 0 {
			continue
		}

		f := strings.Fields(s)
		pid, _ := strconv.Atoi(f[0])
		processPath := f[1]

		if strings.Index(processPath, path) != -1 {
			log.Print("Kill process: " + processPath)

			p, _ := os.FindProcess(pid)
			p.Kill()
			break
		}
	}
}
