package main

import (
	"log"
	"strconv"

	"github.com/terryhycheng/valheim-backup-script/backup"
	"github.com/terryhycheng/valheim-backup-script/compressor"
	"github.com/terryhycheng/valheim-backup-script/config"
	"github.com/terryhycheng/valheim-backup-script/flagParser"
	"github.com/terryhycheng/valheim-backup-script/logs"
	"github.com/terryhycheng/valheim-backup-script/s3"
)

func main() {
	daysToKeep, err := strconv.Atoi(config.Envs["DAYS_TO_KEEP"])
	if err != nil {
		log.Fatalf("❌ Failed to convert DAYS_TO_KEEP to int: %v", err)
	}
	flagParser := flagParser.New()
	s3 := s3.New(config.Envs["S3_BUCKET_NAME"], config.Envs["S3_FOLDER_NAME"])
	compressor := compressor.New()

	backup := backup.New(&backup.BackupConfig{
		DestinationDir: flagParser.Destination(),
		SourceDir:      flagParser.Source(),
		ArchiveName:    flagParser.ArchiveName(),
		S3:             s3,
		Compressor:     compressor,
		DaysToKeep:     daysToKeep,
	})

	err = backup.DeleteOldBackup()
	if err != nil {
		log.Fatalf("❌ Failed to delete old backup: %v", err)
	}

	err = backup.CreateNewBackup()
	if err != nil {
		log.Fatalf("❌ Failed to create backup: %v", err)
	}

	err = logs.DeleteOldLogs(config.Envs["LOGS_FOLDER"], daysToKeep)
	if err != nil {
		log.Fatalf("❌ Failed to delete old logs: %v", err)
	}

}
