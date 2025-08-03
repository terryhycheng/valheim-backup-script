package backup

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/terryhycheng/valheim-backup-script/compressor"
	"github.com/terryhycheng/valheim-backup-script/s3"
)

type Backup interface {
	CreateNewBackup() error
	DeleteOldBackup() error
}

type BackupConfig struct {
	DestinationDir string
	SourceDir      string
	ArchiveName    string
	S3             s3.S3Helper
	Compressor     compressor.Compressor
	DaysToKeep     int
}

func New(config *BackupConfig) Backup {
	return &BackupConfig{
		DestinationDir: config.DestinationDir,
		SourceDir:      config.SourceDir,
		ArchiveName:    config.ArchiveName,
		S3:             config.S3,
		Compressor:     config.Compressor,
		DaysToKeep:     config.DaysToKeep,
	}
}

func (b *BackupConfig) CreateNewBackup() error {
	// variables
	archiveName := fmt.Sprintf("%s-%s.tar.gz", b.ArchiveName, time.Now().Format("2006-01-02_15:04:05"))
	filePath := filepath.Join(b.DestinationDir, archiveName)

	// Create new backup
	fmt.Println("📦 Creating backup...")

	err := b.Compressor.CreateTarGz(b.DestinationDir, b.SourceDir, archiveName)
	if err != nil {
		return fmt.Errorf("❌ Failed to create backup: %w", err)
	}

	// Upload to S3
	fmt.Printf("☁️ Uploading %s to S3 bucket %s...\n", archiveName, b.S3.Bucket())
	err = b.S3.Upload(filePath)
	if err != nil {
		return fmt.Errorf("❌ Failed to upload backup: %w", err)
	}

	fmt.Println("✅ Backup created successfully")

	return nil
}

func (b *BackupConfig) DeleteOldBackup() error {
	// Delete old backups
	fmt.Println("🧹 Checking for backups older than 90 days to delete...")
	oldBackups, err := b.S3.ListOldObjects(b.DaysToKeep)
	if err != nil {
		return fmt.Errorf("❌ Failed to list old backups: %w", err)
	}
	if len(oldBackups) > 0 {
		fmt.Printf("Found %d old backups, deleting...\n", len(oldBackups))
		err = b.S3.Delete(oldBackups)
		if err != nil {
			return fmt.Errorf("❌ Failed to delete old backups: %w", err)
		}
		fmt.Println("✅ Old backups deleted successfully")
	} else {
		fmt.Printf("🟢 No old backups (older than %d days) found, skipping deletion.\n", 90)
	}
	return nil
}
