package db

import "fmt"

type MysqlDBConfig struct {
	WRITE           DBBase `yaml:"WRITE"`           // write database
	READ            DBBase `yaml:"READ"`            // read database
	ConnMaxIdleTime int    `yaml:"ConnMaxIdleTime"` // max idle time of connection, unit: second
	MaxOpenConns    int    `yaml:"MaxOpenConns"`    // max open connections
	MaxIdleConns    int    `yaml:"MaxIdleConns"`    // max idle connections
	ConnMaxLifetime int    `yaml:"ConnMaxLifetime"`

	SlowThresholdUs int    `yaml:"SlowThresholdUs"` // slow SQL threshold, unit: us. if slow SQL log will print SQL statement
	SqlConsole      bool   `yaml:"SqlConsole"`      // sql print or not
	ConnParam       string `yaml:"ConnParam"`       // charset=utf8&parseTime=true&loc=Local
}

type DBBase struct {
	Server string `yaml:"Server"` // database address
	Port   int    `yaml:"Port"`   // database port
	DBName string `yaml:"DBName"` // database name
	User   string `yaml:"User"`   // database access account
	Psw    string `yaml:"Psw"`    // database password
}

func (db *DBBase) DBUrl(connParam string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		db.User, db.Psw, db.Server, db.Port, db.DBName, connParam)
}
