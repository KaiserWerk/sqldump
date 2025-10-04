package sqldump

import (
	"database/sql"
	"fmt"

	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLBackupFile(username, password, hostname, port, dbname, filebasename string) func() (string, error) {
	return func() (string, error) {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8", username, password, hostname, port, dbname))
		if err != nil {
			return "", err
		}

		dumper, err := mysqldump.Register(db, "./", filebasename)
		if err != nil {
			return "", err
		}

		resultFilename, err := dumper.Dump()
		if err != nil {
			return "", err
		}

		return resultFilename, nil
	}
}

func GetMySQLBackupFileByDSN(dsn string, filebasename string) func() (string, error) {
	return func() (string, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return "", err
		}

		dumper, err := mysqldump.Register(db, "./", filebasename)
		if err != nil {
			return "", err
		}

		resultFilename, err := dumper.Dump()
		if err != nil {
			return "", err
		}

		return resultFilename, nil
	}
}
