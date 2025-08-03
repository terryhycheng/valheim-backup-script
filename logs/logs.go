package logs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DeleteOldLogs deletes old logs from the specified folder.
//
// Takes `folderName` as the name of the folder to delete logs from (this folder should be in the user's home directory).
// Takes `daysOld` as the number of days to keep logs for.
func DeleteOldLogs(folderName string, daysOld int) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	folder := filepath.Join(home, folderName)
	cmd := exec.Command("find", folder, "-name", "backup-*.log", "-type", "f", "-mtime", fmt.Sprintf("+%d", daysOld), "-delete")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
