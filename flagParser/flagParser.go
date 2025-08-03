package flagParser

import (
	"flag"
	"log"
)

type FlagParser interface {
	Destination() string
	Source() string
	ArchiveName() string
}

type flagParser struct {
	destination string
	source      string
	archiveName string
}

func New() FlagParser {
	fp := &flagParser{}
	flag.StringVar(&fp.destination, "d", "", "The destination directory to save the backup")
	flag.StringVar(&fp.source, "s", "", "The source directory to backup")
	flag.StringVar(&fp.archiveName, "a", "valheim-backup", "The name of the archive to save the backup")

	flag.Parse()

	if fp.destination == "" || fp.source == "" {
		log.Fatalf("❌ Missing required flags: -d, -s, -a")
	}

	return fp
}

func (f *flagParser) Destination() string { return f.destination }
func (f *flagParser) Source() string      { return f.source }
func (f *flagParser) ArchiveName() string { return f.archiveName }
