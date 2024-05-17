package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/saratkumar-yb/infinityapi/db"
)

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

func jsonResponse(w http.ResponseWriter, status string, message string) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": status, "message": message}
	json.NewEncoder(w).Encode(response)
}

func InsertYbaHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var yba Yba
	err := json.NewDecoder(r.Body).Decode(&yba)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	err = db.InsertYba(yba)
	if err != nil {
		log.Printf("Failed to insert into yba: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}

	jsonResponse(w, "successful", "")
}

func InsertYbdbHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ybdb Ybdb
	err := json.NewDecoder(r.Body).Decode(&ybdb)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	err = db.InsertYbdb(ybdb)
	if err != nil {
		log.Printf("Failed to insert into ybdb: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}

	jsonResponse(w, "successful", "")
}

func InsertCompatibilityHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var compatibility Compatibility
	err := json.NewDecoder(r.Body).Decode(&compatibility)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	err = db.InsertCompatibility(compatibility)
	if err != nil {
		log.Printf("Failed to insert into yba_ybdb_compatibility: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}

	jsonResponse(w, "successful", "")
}

func GetCompatibleYbdbHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CompatibilityRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonResponse(w, "failed", err.Error())
		return
	}

	ybdbs, err := db.GetCompatibleYbdb(req.YbaVersion)
	if err != nil {
		log.Printf("Failed to query compatible ybdb versions: %v", err)
		jsonResponse(w, "failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ybdbs)
}
