package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/saratkumar-yb/infinityapi/config"
	"github.com/saratkumar-yb/infinityapi/db"
	"github.com/saratkumar-yb/infinityapi/router"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func setupRouter() *httprouter.Router {
	return router.NewRouter()
}

func setupDatabase() (*sql.DB, error) {
	config.LoadConfig()
	return db.Connect()
}

func clearDatabase(db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE yba_ybdb_compatibility CASCADE")
	if err != nil {
		log.Fatalf("Failed to clear yba_ybdb_compatibility table: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE yba CASCADE")
	if err != nil {
		log.Fatalf("Failed to clear yba table: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE ybdb CASCADE")
	if err != nil {
		log.Fatalf("Failed to clear ybdb table: %v", err)
	}
}

func TestInsertYbaHandler(t *testing.T) {
	db, err := setupDatabase()
	assert.NoError(t, err)
	defer db.Close()

	clearDatabase(db)

	router := setupRouter()

	payload := []byte(`{
		"version": "1.0",
		"type": "type1",
		"architecture": "arch1",
		"platform": "platform1",
		"commit": "commit1",
		"branch": "branch1"
	}`)

	req, err := http.NewRequest("POST", "/yba", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response Response
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "successful", response.Status)
}

func TestInsertYbdbHandler(t *testing.T) {
	db, err := setupDatabase()
	assert.NoError(t, err)
	defer db.Close()

	clearDatabase(db)

	router := setupRouter()

	payload := []byte(`{
		"version": "1.0",
		"type": "type1",
		"architecture": "arch1",
		"platform": "platform1",
		"download_url": "http://example.com/download",
		"commit": "commit1",
		"branch": "branch1"
	}`)

	req, err := http.NewRequest("POST", "/ybdb", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response Response
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "successful", response.Status)
}

func TestInsertCompatibilityHandler(t *testing.T) {
	db, err := setupDatabase()
	assert.NoError(t, err)
	defer db.Close()

	clearDatabase(db)

	router := setupRouter()

	// Insert test data into yba and ybdb tables
	_, err = db.Exec(`INSERT INTO yba (version, type, architecture, platform, commit, branch) VALUES ('1.0', 'type1', 'arch1', 'platform1', 'commit1', 'branch1')`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO ybdb (version, type, architecture, platform, download_url, commit, branch) VALUES ('1.0', 'type1', 'arch1', 'platform1', 'http://example.com/download', 'commit1', 'branch1')`)
	assert.NoError(t, err)

	payload := []byte(`{
		"yba_versions": ["1.0"],
		"ybdb_versions": ["1.0"]
	}`)

	req, err := http.NewRequest("POST", "/compatibility", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response Response
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "successful", response.Status)
}

func TestGetCompatibleYbdbHandler(t *testing.T) {
	router := setupRouter()

	// Ensure there is data in the database for this test
	db, err := setupDatabase()
	assert.NoError(t, err)
	defer db.Close()

	clearDatabase(db)

	// Insert test data into yba and ybdb tables
	_, err = db.Exec(`INSERT INTO yba (version, type, architecture, platform, commit, branch) VALUES ('1.0', 'type1', 'arch1', 'platform1', 'commit1', 'branch1')`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO ybdb (version, type, architecture, platform, download_url, commit, branch) VALUES ('1.0', 'type1', 'arch1', 'platform1', 'http://example.com/download', 'commit1', 'branch1')`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO yba_ybdb_compatibility (yba_version, ybdb_version) VALUES ('1.0', '1.0')`)
	assert.NoError(t, err)

	payload := []byte(`{
		"yba_version": "1.0"
	}`)

	req, err := http.NewRequest("POST", "/compatibility_list", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []db.Ybdb
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, "1.0", response[0].Version)
	assert.Equal(t, "type1", response[0].Type)
	assert.Equal(t, "arch1", response[0].Architecture)
	assert.Equal(t, "platform1", response[0].Platform)
	assert.Equal(t, "http://example.com/download", response[0].DownloadURL)
	assert.Equal(t, "commit1", response[0].Commit)
	assert.Equal(t, "branch1", response[0].Branch)
}
