package sqldump

import (
	"database/sql"
	"fmt"

	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLBackupFile(username, password, hostname, port, dbname string) (string, error) {
	dumpFilenameFormat := dbname
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8", username, password, hostname, port, dbname))
	if err != nil {
		return "", err
	}

	dumper, err := mysqldump.Register(db, "./", dumpFilenameFormat)
	if err != nil {
		return "", err
	}

	resultFilename, err := dumper.Dump()
	if err != nil {
		return "", err
	}

	return resultFilename, nil
}
