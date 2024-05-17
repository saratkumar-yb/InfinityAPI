package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-ini/ini"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost       string
	DBPort       int
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	HTTPListener string
	HTTPPort     int
}

var config Config

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

type CompatibilityRequest struct {
	YbaVersion string `json:"yba_version"`
}

func main() {
	loadConfig()

	migrateCmd := flag.Bool("migrate", false, "Migrate the database")
	startServerCmd := flag.Bool("startserver", false, "Start the API server")
	flag.Parse()

	if *migrateCmd {
		err := migrate()
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration successful")
	} else if *startServerCmd {
		startServer()
	} else {
		log.Println("No command provided. Use -migrate or -startserver.")
	}
}

func loadConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	dbSection := cfg.Section("db")
	serverSection := cfg.Section("server")

	config = Config{
		DBHost:       dbSection.Key("host").String(),
		DBPort:       dbSection.Key("port").MustInt(),
		DBUser:       dbSection.Key("user").String(),
		DBPassword:   dbSection.Key("password").String(),
		DBName:       dbSection.Key("dbname").String(),
		DBSSLMode:    dbSection.Key("sslmode").String(),
		HTTPListener: serverSection.Key("http_listener").String(),
		HTTPPort:     serverSection.Key("http_port").MustInt(),
	}
}

func dbConnect() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode)
	return sql.Open("postgres", connStr)
}

func migrate() error {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}

func startServer() {
	router := httprouter.New()
	router.POST("/yba", insertYbaHandler)
	router.POST("/ybdb", insertYbdbHandler)
	router.POST("/compatibility", insertCompatibilityHandler)
	router.POST("/compatibility_list", getCompatibleYbdbHandler)

	address := fmt.Sprintf("%s:%d", config.HTTPListener, config.HTTPPort)
	log.Fatal(http.ListenAndServe(address, router))
}

func jsonResponse(w http.ResponseWriter, status string, message string) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": status, "message": message}
	json.NewEncoder(w).Encode(response)
}

func insertYbaHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var yba Yba
	err := json.NewDecoder(r.Body).Decode(&yba)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	db, err := dbConnect()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO yba (version, type, architecture, platform, commit, branch) VALUES ($1, $2, $3, $4, $5, $6)",
		yba.Version, yba.Type, yba.Architecture, yba.Platform, yba.Commit, yba.Branch)
	if err != nil {
		log.Printf("Failed to insert into yba: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}

	jsonResponse(w, "successful", "")
}

func insertYbdbHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ybdb Ybdb
	err := json.NewDecoder(r.Body).Decode(&ybdb)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	db, err := dbConnect()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO ybdb (version, type, architecture, platform, download_url, commit, branch) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		ybdb.Version, ybdb.Type, ybdb.Architecture, ybdb.Platform, ybdb.DownloadURL, ybdb.Commit, ybdb.Branch)
	if err != nil {
		log.Printf("Failed to insert into ybdb: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}

	jsonResponse(w, "successful", "")
}

func insertCompatibilityHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var compatibility Compatibility
	err := json.NewDecoder(r.Body).Decode(&compatibility)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	db, err := dbConnect()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}
	defer db.Close()

	for _, ybaVersion := range compatibility.YbaVersions {
		for _, ybdbVersion := range compatibility.YbdbVersions {
			_, err = db.Exec("INSERT INTO yba_ybdb_compatibility (yba_version, ybdb_version) VALUES ($1, $2)",
				ybaVersion, ybdbVersion)
			if err != nil {
				log.Printf("Failed to insert into yba_ybdb_compatibility: %v", err)
				jsonResponse(w, "failed", err.Error())
				return
			}
		}
	}

	jsonResponse(w, "successful", "")
}

func getCompatibleYbdbHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CompatibilityRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	db, err := dbConnect()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT ybdb.version, ybdb.type, ybdb.architecture, ybdb.platform, ybdb.download_url, ybdb.commit, ybdb.branch
		FROM ybdb
		INNER JOIN yba_ybdb_compatibility ON ybdb.version = yba_ybdb_compatibility.ybdb_version
		WHERE yba_ybdb_compatibility.yba_version = $1`, req.YbaVersion)
	if err != nil {
		log.Printf("Failed to query compatible ybdb versions: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}
	defer rows.Close()

	var ybdbs []Ybdb
	for rows.Next() {
		var ybdb Ybdb
		err := rows.Scan(&ybdb.Version, &ybdb.Type, &ybdb.Architecture, &ybdb.Platform, &ybdb.DownloadURL, &ybdb.Commit, &ybdb.Branch)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			jsonResponse(w, "failed", err.Error())
			return
		}
		ybdbs = append(ybdbs, ybdb)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}

	if len(ybdbs) == 0 {
		ybdbs = []Ybdb{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ybdbs)
}
