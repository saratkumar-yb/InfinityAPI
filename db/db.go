package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/saratkumar-yb/infinityapi/config"
)

var DB *sql.DB

type Yba struct {
	Version      string `json:"version"`
	Type         string `json:"type"`
	Architecture string `json:"architecture"`
	Platform     string `json:"platform"`
	Commit       string `json:"commit"`
	Branch       string `json:"branch"`
}

type Ybdb struct {
	Version      string `json:"version"`
	Type         string `json:"type"`
	Architecture string `json:"architecture"`
	Platform     string `json:"platform"`
	DownloadURL  string `json:"download_url"`
	Commit       string `json:"commit"`
	Branch       string `json:"branch"`
}

type Compatibility struct {
	YbaVersions  []string `json:"yba_versions"`
	YbdbVersions []string `json:"ybdb_versions"`
}

func Connect() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBUser, config.AppConfig.DBPassword, config.AppConfig.DBName, config.AppConfig.DBSSLMode)
	var err error
	DB, err = sql.Open("postgres", connStr)
	return err
}

func Migrate() error {
	err := Connect()
	if err != nil {
		return err
	}
	defer DB.Close()

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}

func Insert(sqlStatement string, args ...interface{}) error {
	_, err := DB.Exec(sqlStatement, args...)
	return err
}

func Query(sqlStatement string, args ...interface{}) (*sql.Rows, error) {
	return DB.Query(sqlStatement, args...)
}

func InsertYba(yba Yba) error {
	sqlStatement := "INSERT INTO yba (version, type, architecture, platform, commit, branch) VALUES ($1, $2, $3, $4, $5, $6)"
	return Insert(sqlStatement, yba.Version, yba.Type, yba.Architecture, yba.Platform, yba.Commit, yba.Branch)
}

func InsertYbdb(ybdb Ybdb) error {
	sqlStatement := "INSERT INTO ybdb (version, type, architecture, platform, download_url, commit, branch) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	return Insert(sqlStatement, ybdb.Version, ybdb.Type, ybdb.Architecture, ybdb.Platform, ybdb.DownloadURL, ybdb.Commit, ybdb.Branch)
}

func InsertCompatibility(compatibility Compatibility) error {
	for _, ybaVersion := range compatibility.YbaVersions {
		for _, ybdbVersion := range compatibility.YbdbVersions {
			sqlStatement := "INSERT INTO yba_ybdb_compatibility (yba_version, ybdb_version) VALUES ($1, $2)"
			err := Insert(sqlStatement, ybaVersion, ybdbVersion)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetCompatibleYbdb(ybaVersion string) ([]Ybdb, error) {
	sqlStatement := `
		SELECT ybdb.version, ybdb.type, ybdb.architecture, ybdb.platform, ybdb.download_url, ybdb.commit, ybdb.branch
		FROM ybdb
		INNER JOIN yba_ybdb_compatibility ON ybdb.version = yba_ybdb_compatibility.ybdb_version
		WHERE yba_ybdb_compatibility.yba_version = $1`
	rows, err := Query(sqlStatement, ybaVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ybdbs []Ybdb
	for rows.Next() {
		var ybdb Ybdb
		err := rows.Scan(&ybdb.Version, &ybdb.Type, &ybdb.Architecture, &ybdb.Platform, &ybdb.DownloadURL, &ybdb.Commit, &ybdb.Branch)
		if err != nil {
			return nil, err
		}
		ybdbs = append(ybdbs, ybdb)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ybdbs) == 0 {
		ybdbs = []Ybdb{}
	}

	return ybdbs, nil
}
