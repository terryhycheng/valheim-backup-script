package compressor

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Compressor interface {
	CreateTarGz(destinationDir, sourceDir, archiveName string) error
}

type compressor struct {
}

func New() Compressor {
	return &compressor{}
}

// CreateTarGz creates a tar.gz file of the entire source directory and saves it to the destination path, by using the tar command from Linux
//
// Examples:
//
//	CreateTarGz("/tmp", "/mnt/HC_Volume_102879414/valheim-server/config/worlds_local", "valheim-backup-2025-08-03")
func (c *compressor) CreateTarGz(destinationDir, sourceDir, archiveName string) error {
	// Add .tar.gz extension if not provided
	if filepath.Ext(archiveName) != ".gz" {
		archiveName = archiveName + ".tar.gz"
	}

	destinationPath := filepath.Join(destinationDir, archiveName)
	cmd := exec.Command("tar", "-czf", destinationPath, "-C", sourceDir, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
