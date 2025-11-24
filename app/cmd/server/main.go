package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"pyramid_puzzle/internal/db"
	"pyramid_puzzle/internal/puzzle"

	"github.com/joho/godotenv"

	"github.com/google/uuid"
)

type WordResult struct {
	Word       string `json:"word"`
	Length     int    `json:"length"`
	LetterMask int64  `json:"letter_mask"`
}

func main() {
	godotenv.Load(".env")

	puzzleDB, err := db.OpenSQLite("./data/pyramid_puzzle.db")
	if err != nil {
		log.Fatalf("❌ Failed to open SQLite DB: %v", err)
	}

	userDB, err := db.OpenSQLite("./secure/data/userdata.db")
	if err != nil {
		log.Fatalf("❌ Failed to open SQLite DB: %v", err)
	}

	svc, err := puzzle.NewService(puzzleDB, userDB)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/v1.0/validate-guess", svc.HandleValidateGuessEndpoint)
	http.HandleFunc("/api/v1.0/get-today-remaining-guesses", svc.HandleGetTodayRemainingGuessesEndpoint)
	http.HandleFunc("/api/v1.0/load-today-game", svc.HandleLoadTodayGameEndpoint)
	http.HandleFunc("/api/v1.0/get-next-word-hint", svc.HandleGetNextWordHint)
	http.HandleFunc("/api/v1.0/get-user", handleUser)
	// err = svc.PrintFullPuzzle()
	if err != nil {
		log.Fatalf("❌ Failed to print puzzle: %v", err)
	}

	log.Println("🚀 Running server on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("⛔ Invalid method:", r.Method)
		return
	}

	var req struct {
		Fingerprint string `json:"fingerprint"` // may be blank
		UserAgent   string `json:"user_agent"`
		DeviceType  string `json:"device_type"`
		OS          string `json:"os"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Failed to parse JSON: %v\n", err)
		return
	}

	userDB, err := db.OpenSQLite("./secure/data/userdata.db")
	if err != nil {
		log.Fatalf("❌ Failed to open SQLite DB: %v", err)
	}

	// Try to read user by fingerprint
	readSQLBytes, err := os.ReadFile("./secure/queries/read_user.sql")
	if err != nil {
		log.Printf("❌ Failed to read read_user.sql: %v\n", err)
		return
	}
	readStmt, err := userDB.Prepare(string(readSQLBytes))
	if err != nil {
		log.Printf("❌ Failed to prepare read_user SQL: %v\n", err)
		return
	}
	defer readStmt.Close()

	var (
		existingID int64
		fp         string
		ipAddr     string
		ua         string
		devType    string
		osValue    string
		createdAt  string
	)

	err = readStmt.QueryRow(
		sql.Named("fingerprint", req.Fingerprint),
	).Scan(
		&existingID,
		&fp,
		&ipAddr,
		&ua,
		&devType,
		&osValue,
		&createdAt,
	)

	if err == nil {
		// ✅ Found existing user
		log.Printf("✅ Found user by fingerprint: %s (id=%d)\n", fp, existingID)

		json.NewEncoder(w).Encode(struct {
			Found       bool   `json:"found"`
			Fingerprint string `json:"fingerprint"`
		}{
			Found:       true,
			Fingerprint: fp,
		})
		return
	} else if err != sql.ErrNoRows {
		log.Printf("❌ Failed to execute user lookup: %v\n", err)
		return
	}

	// 🆕 Not found → insert new user
	newFingerprint := uuid.NewString()
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	insertSQLBytes, err := os.ReadFile("./secure/queries/insert_user.sql")
	if err != nil {
		log.Printf("❌ Failed to read insert_user.sql: %v\n", err)
		return
	}
	insertStmt, err := userDB.Prepare(string(insertSQLBytes))
	if err != nil {
		log.Printf("❌ Failed to prepare insert_user SQL: %v\n", err)
		return
	}
	defer insertStmt.Close()

	var userID int64
	err = insertStmt.QueryRow(
		sql.Named("fingerprint", newFingerprint),
		sql.Named("ip", ip),
		sql.Named("user_agent", req.UserAgent),
		sql.Named("device_type", req.DeviceType),
		sql.Named("os", req.OS),
	).Scan(&userID)

	if err != nil {
		log.Printf("❌ Failed to insert new user: %v\n", err)
		return
	}

	log.Printf("🆕 Inserted new user: id=%d fingerprint=%s ip=%s\n", userID, newFingerprint, ip)

	json.NewEncoder(w).Encode(struct {
		Found       bool   `json:"found"`
		Fingerprint string `json:"fingerprint"`
	}{
		Found:       false,
		Fingerprint: newFingerprint,
	})
}
