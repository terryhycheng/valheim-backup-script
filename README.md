# Valheim Backup Script

I am hosting a dedicated server for Valheim on a VPS (Linux based), and I need a backup script that is reliable and easy to amend. Golang is the language of choice for this project as it's easy to build and deploy.

## What does this script do?

- Deletes old backups (more than 90 days)
- Create a `.tar.gz` file of the backup
- Upload the backup to a bucket on AWS S3 (currently `valheim-backup` on DigitalOcean Spaces)

## How to use this script?

- Clone the repository or download the release from [here](https://github.com/terryhycheng/valheim-backup-script/releases)
- Copy the `.env.example` file to `.env` and fill in the values
- Run `go run main.go` or `./valheim-backup-script` to run the script
- The script will ask for the backup directory and the remote server details.

### Run the script from the source code

```bash
$ git clone https://github.com/terryhycheng/valheim-backup-script.git
$ cd valheim-backup-script

# Install dependencies
$ go mod download

# Run the script
$ go run main.go
```

### Download the release

```bash
# Download the release
$ wget https://github.com/terryhycheng/valheim-backup-script/releases/download/v1.0.0/valheim-backup-script

# Make the script executable
$ chmod +x valheim-backup-script

# Run the script
$ ./valheim-backup-script
```

## Roadmap

- [ ] Cover the script with tests
