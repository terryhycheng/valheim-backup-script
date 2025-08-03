package logs

import (
	"fmt"
	"os"
	"os/exec"
)

func DeleteOldLogs(folder string, daysOld int) error {
	cmd := exec.Command("find", folder, "-name", "backup-*.log", "-type", "f", "-mtime", fmt.Sprintf("+%d", daysOld), "-delete")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
