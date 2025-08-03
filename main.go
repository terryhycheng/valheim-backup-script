package main

import (
	"flag"
	"log"

	"github.com/terryhycheng/valheim-backup-script/backup"
	"github.com/terryhycheng/valheim-backup-script/compressor"
	_ "github.com/terryhycheng/valheim-backup-script/config"
	"github.com/terryhycheng/valheim-backup-script/s3"
)

func main() {
	destination := flag.String("d", "", "destination directory")
	source := flag.String("s", "", "source directory")
	archiveName := flag.String("a", "", "archive name")

	if *destination == "" || *source == "" || *archiveName == "" {
		log.Fatalf("❌ Missing required flags: -d, -s, -a")
	}

	flag.Parse()

	s3 := s3.New("valheim-backup", "valheim-backup")
	compressor := compressor.New()

	backup := backup.New(&backup.BackupConfig{
		DestinationDir: *destination,
		SourceDir:      *source,
		ArchiveName:    *archiveName,
		S3:             s3,
		Compressor:     compressor,
	})

	err := backup.DeleteOldBackup()
	if err != nil {
		log.Fatalf("❌ Failed to delete old backup: %v", err)
	}

	err = backup.CreateNewBackup()
	if err != nil {
		log.Fatalf("❌ Failed to create backup: %v", err)
	}

}
