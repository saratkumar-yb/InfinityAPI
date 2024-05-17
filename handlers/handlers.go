package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/saratkumar-yb/infinityapi/db"
)

func jsonResponse(w http.ResponseWriter, status string, message string) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": status, "message": message}
	json.NewEncoder(w).Encode(response)
}

func InsertYbaHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var yba db.Yba
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
	var ybdb db.Ybdb
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
	var compatibility db.Compatibility
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
	var req db.CompatibilityRequest
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
