# SQLDump

`SQLDump` is a Go-based utility for creating database file backups (MySQL only) and uploading them to a remote storage box via SFTP (or any file, actually). It is designed to automate the process of database backup and storage, ensuring data safety and accessibility.

## Features

- **MySQL Backup**: Generate MySQL database dumps using `go-mysqldump`.
- **SFTP Upload**: Upload backup files to a remote storage box via SFTP.
- **Scheduled Uploads**: Automate periodic uploads with configurable intervals.
- **File Rotation**: Manage the number of backup files stored remotely by rotating filenames.

## Requirements

- Go 1.25.1 or higher
- Remote storage box with SFTP access and external access enabled

## Installation

   ```bash
   go get github.com/KaiserWerk/sqldump
   ```

## Upload to Storage Box

Create an `Uploader` instance and use the `ScheduleUpload` method to automate uploads:

```go
uploader := NewUploader("username", "password", "host:22")
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

uploader.ScheduleUpload(ctx, GetMySQLBackupFile("username", "password", "hostname", "3306", "dbname"), "/remote/dir", time.Hour*24, 7)
```

This will upload the file returned by the funcition filenameFunc to the remote directory every 24 hours and rotate through seven differently appended file names.
