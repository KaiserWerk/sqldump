package sqldump

func GetSQLiteBackupFile(file string) func() (string, error) {
	return func() (string, error) {
		return file, nil
	}
}
