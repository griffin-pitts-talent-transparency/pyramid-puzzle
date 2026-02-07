package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"lebron-games/internal/db"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func HandleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("⛔ Invalid method:", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var incomingRequest struct {
		CognitoSub      string `json:"cognito_sub"`
		CognitoUsername string `json:"cognito_username"`
		Fingerprint     string `json:"fingerprint"`
		UserAgent       string `json:"user_agent"`
		DeviceType      string `json:"device_type"`
		OS              string `json:"os"`
	}

	if err := json.NewDecoder(r.Body).Decode(&incomingRequest); err != nil {
		log.Printf("❌ Failed to parse JSON: %v\n", err)
		return
	}

	userDB, err := db.OpenSQLiteDBByName("userdata")
	if err != nil {
		log.Fatalf("❌ Failed to open SQLite DB: %v", err)
	}
	defer userDB.Close()

	readFingerprintUserStmt, err := db.PrepareSQLStatementByName(userDB, "read_user_by_fingerprint.sql", "userdata")
	if err != nil {
		return
	}
	defer readFingerprintUserStmt.Close()

	readSubUserStmt, err := db.PrepareSQLStatementByName(userDB, "read_user_by_sub.sql", "userdata")
	if err != nil {
		return
	}
	defer readSubUserStmt.Close()

	updateStmt, err := db.PrepareSQLStatementByName(userDB, "update_user.sql", "userdata")
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer updateStmt.Close()

	insertStmt, err := db.PrepareSQLStatementByName(userDB, "insert_user.sql", "userdata")
	if err != nil {
		return
	}
	defer insertStmt.Close()

	insertLoginStreakStmt, err := db.PrepareSQLStatementByName(userDB, "insert_login_streak.sql", "userdata")
	if err != nil {
		return
	}
	defer insertLoginStreakStmt.Close()

	readLoginStreakStmt, err := db.PrepareSQLStatementByName(userDB, "read_login_streak.sql", "userdata")
	if err != nil {
		return
	}
	defer readLoginStreakStmt.Close()

	var (
		userID          int64
		cognitoSub      sql.NullString
		CognitoUsername sql.NullString
		fingerprint     string
		ip              string
		userAgent       string
		deviceType      string
		osValue         string
		createdAt       string
	)

	if incomingRequest.CognitoSub != "" {
		err = readSubUserStmt.QueryRow(
			sql.Named("cognito_sub", incomingRequest.CognitoSub),
		).Scan(
			&userID,
			&cognitoSub,
			&CognitoUsername,
			&fingerprint,
			&ip,
			&userAgent,
			&deviceType,
			&osValue,
			&createdAt,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// log.Println("No user found by cognito_sub")
			} else {
				log.Printf("❌ 'read_sub_user.sql' failed: %v\n", err)
			}
		} else {
			// Successful row scan
			if fingerprint != "" {
				// found = true
			}
		}
	}
	if userID <= 0 && incomingRequest.Fingerprint != "" {
		err = readFingerprintUserStmt.QueryRow(
			sql.Named("fingerprint", incomingRequest.Fingerprint),
		).Scan(
			&userID,
			&cognitoSub,
			&CognitoUsername,
			&fingerprint,
			&ip,
			&userAgent,
			&deviceType,
			&osValue,
			&createdAt,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// log.Println("No user found by fingerprint")
			} else {
				log.Printf("❌ 'read_user_by_fingerprint.sql' failed: %v\n", err)
			}
		} else {
			// Successful row scan
			if fingerprint != "" {
				// found = true
			}
		}
	}

	var newUser = false
	shouldUpdateSub := !cognitoSub.Valid || strings.TrimSpace(cognitoSub.String) == ""
	shouldUpdateUsername := !CognitoUsername.Valid || strings.TrimSpace(CognitoUsername.String) == ""

	incomingHasValidSub := strings.TrimSpace(incomingRequest.CognitoSub) != ""
	incomingHasValidUsername := strings.TrimSpace(incomingRequest.CognitoUsername) != ""

	needsSubUpdate := shouldUpdateSub && incomingHasValidSub
	needsUsernameUpdate := shouldUpdateUsername && incomingHasValidUsername && incomingRequest.CognitoUsername != CognitoUsername.String

	if userID > 0 && (needsSubUpdate || needsUsernameUpdate) {
		// found = true
		// update user
		// log.Printf("updating user")
		_, err = updateStmt.Exec(
			sql.Named("cognito_sub", incomingRequest.CognitoSub),
			sql.Named("cognito_username", incomingRequest.CognitoUsername),
			sql.Named("user_id", userID),
		)
		if err != nil {
			log.Printf("❌ Failed to update user with cognito_sub: %v\n", err)
			// http.Error(w, "Server error", http.StatusInternalServerError)
		}
		// log.Printf("Bound cognito_sub to user id=%d\n", userID)
	} else if userID <= 0 {
		// add user
		var newUserID int64
		newFingerprint := uuid.NewString()
		userIP := r.Header.Get("X-Forwarded-For")
		if userIP == "" {
			userIP = r.RemoteAddr
		}

		err = insertStmt.QueryRow(
			sql.Named("cognito_sub", nullIfEmpty(incomingRequest.CognitoSub)),
			sql.Named("cognito_username", nullIfEmpty(incomingRequest.CognitoUsername)),
			sql.Named("fingerprint", newFingerprint),
			sql.Named("ip", userIP),
			sql.Named("user_agent", incomingRequest.UserAgent),
			sql.Named("device_type", incomingRequest.DeviceType),
			sql.Named("os", incomingRequest.OS),
		).Scan(&newUserID)

		if err != nil {
			log.Printf("❌ Failed to insert new user: %v\n", err)
		} else if newUserID > 0 {
			// new user insert succeeded
			newUser = true
			fingerprint = newFingerprint
		}
	}
	// log.Print(fingerprint)
	// log.Print(newUser)

	// update user login streak (if applicable)
	var loginStreak int64
	if userID > 0 && incomingRequest.CognitoSub != "" {
		_, err = insertLoginStreakStmt.Exec(
			sql.Named("user_id", userID),
		)
		if err != nil {
			log.Printf("❌ Failed to insert login streak: %v\n", err)
		}

		err = readLoginStreakStmt.QueryRow(
			sql.Named("user_id", userID),
		).Scan(&loginStreak)

		if err != nil {
			if err == sql.ErrNoRows {
				// loginStreak = 0 // no streak yet
			} else {
				log.Printf("❌ Failed to read login streak: %v\n", err)
			}
		}
	}
	/*
		if userID > 0 && incomingRequest.CognitoSub != "" {
			_, err = insertLoginStreakStmt.Exec(
				sql.Named("user_id", userID),
			)
			if err != nil {
				log.Printf("❌ Failed to insert login streak: %v\n", err)
			}
		}
	*/

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		NewUser     bool   `json:"new_user"`
		Fingerprint string `json:"fingerprint"`
		LoginStreak int64  `json:"login_streak"`
	}{
		NewUser:     newUser,
		Fingerprint: fingerprint,
		LoginStreak: loginStreak,
	})
}

func nullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
