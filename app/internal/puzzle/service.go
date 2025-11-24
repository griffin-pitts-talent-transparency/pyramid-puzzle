package puzzle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pyramid_puzzle/internal/secret"
	"pyramid_puzzle/internal/seed"
	"strings"
	"time"
)

type Service struct {
	puzzleDB    *sql.DB
	userDB      *sql.DB
	puzzleWords []string // 4,5,6,7 letter words for today
	puzzleDate  string   // "YYYY-MM-DD"
	// puzzleDB statements
	selectWord      *sql.Stmt
	selectWordCount *sql.Stmt
	// userDB statements
	initializeGameState     *sql.Stmt
	readGameState           *sql.Stmt
	readIncorrectGuessCount *sql.Stmt
	updateGameState         *sql.Stmt
}

type WordResult struct {
	Word       string `json:"word"`
	Length     int    `json:"length"`
	LetterMask int64  `json:"letter_mask"`
}

type LetterResult struct {
	Letter string `json:"letter"`
	Status string `json:"status"` // match, present, miss
}

type GameState struct {
	Guesses [][]LetterResult `json:"guesses"`
}

/*
Service struct constructor
*/
func NewService(puzzleDB *sql.DB, userDB *sql.DB) (*Service, error) {
	s := &Service{
		puzzleDB: puzzleDB,
		userDB:   userDB,
	}

	var err error

	// puzzleDB
	s.selectWord, err = prepareSQL(puzzleDB, "./data/select_word.sql")
	if err != nil {
		return nil, err
	}
	s.selectWordCount, err = prepareSQL(puzzleDB, "./data/select_word_count.sql")
	if err != nil {
		return nil, err
	}

	// userDB
	s.readGameState, err = prepareSQL(userDB, "./secure/queries/read_game_state.sql")
	if err != nil {
		return nil, err
	}
	s.readIncorrectGuessCount, err = prepareSQL(userDB, "./secure/queries/read_incorrect_guess_count.sql")
	if err != nil {
		return nil, err
	}
	s.updateGameState, err = prepareSQL(userDB, "./secure/queries/update_game_state.sql")
	if err != nil {
		return nil, err
	}
	s.initializeGameState, err = prepareSQL(userDB, "./secure/queries/initialize_game_state.sql")
	if err != nil {
		return nil, err
	}

	return s, nil
}

/*
Endpoint connectors
*/
func (s *Service) HandleValidateGuessEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Fingerprint string   `json:"fingerprint"`
		Guess       []string `json:"guess"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// validate guess length
	if len(req.Guess) == 0 {
		http.Error(w, "Missing guess", http.StatusBadRequest)
		return
	}
	if len(req.Guess) < 4 || len(req.Guess) > 7 {
		http.Error(w, "Invalid guess length", http.StatusBadRequest)
		return
	}

	guess := strings.Join(req.Guess, "")
	letterResult, remaining, err := s.ValidateGuess(req.Fingerprint, guess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"guessResult":      letterResult,
		"remainingGuesses": remaining,
	})
}

func (s *Service) HandleGetTodayRemainingGuessesEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Fingerprint string `json:"fingerprint"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// validate fingerprint
	if len(req.Fingerprint) == 0 {
		http.Error(w, "Missing fingerprint", http.StatusBadRequest)
		return
	}

	remaining, err := s.CountTodayRemainingGuesses(req.Fingerprint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"remainingGuesses": remaining,
	})
}

func (s *Service) HandleLoadTodayGameEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Fingerprint string `json:"fingerprint"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	gameState, err := s.GetTodayGameState(req.Fingerprint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if gameState == nil {
		if err := s.InitializeTodayGameState(req.Fingerprint); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// After initializing, reload fresh state
		gameState, err = s.GetTodayGameState(req.Fingerprint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if gameState == nil {
			http.Error(w, "Failed to initialize game state", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"gameState": gameState,
	})
}

func (s *Service) HandleGetNextWordHint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Fingerprint string `json:"fingerprint"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hint, err := s.GetNextWordHint(req.Fingerprint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"hint": hint,
	})
}

/*
Endpoint connector logic
*/
func (s *Service) ValidateGuess(fingerprint string, guess string) ([]LetterResult, int, error) {
	remainingGuesses := 0
	var letterResult []LetterResult
	guessTier := len(guess) - 4
	err := s.BuildCurrentPuzzle()
	if err != nil {
		return nil, remainingGuesses, err
	}

	// Get user's game state for the day using fingerprint
	state, err := s.GetTodayGameState(fingerprint)
	if err != nil {
		return nil, remainingGuesses, err
	}
	if state.Guesses == nil {
		// state not initalized for today
		if err := s.InitializeTodayGameState(fingerprint); err != nil {
			return nil, remainingGuesses, err
		}
		remainingGuesses = 11
	} else {
		remainingGuesses = 11 - CountIncorrectGuesses(state.Guesses)
		if remainingGuesses > 0 {
			// check guess and append to state in DB
			letterResult = computeGuessResult(guess, s.puzzleWords[guessTier])
			if err := s.UpdateGameState(fingerprint, letterResult); err != nil {
				return nil, remainingGuesses, err
			}
			for _, lr := range letterResult {
				if lr.Status != "match" {
					remainingGuesses-- // did not guess word correctly
					break
				}
			}
		}
	}

	return letterResult, remainingGuesses, nil
}

func (s *Service) CountTodayRemainingGuesses(fingerprint string) (int, error) {
	remainingGuesses := 0
	// Today's date
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	var incorrectCount int
	err := s.readIncorrectGuessCount.QueryRow(
		sql.Named("fingerprint", fingerprint),
		sql.Named("date", today),
	).Scan(&incorrectCount)

	if err != nil {
		return remainingGuesses, fmt.Errorf("CountTodayIncorrectGuesses failed: %w", err)
	}

	// Remaining = total - incorrect
	remainingGuesses = 11 - incorrectCount
	if remainingGuesses < 0 {
		remainingGuesses = 0 // hard floor
	}

	return remainingGuesses, nil
}

func (s *Service) GetNextWordHint(fingerprint string) ([]LetterResult, error) {
	var letterResult []LetterResult
	err := s.BuildCurrentPuzzle()
	if err != nil {
		return nil, err
	}

	// Get user's game state for the day using fingerprint
	state, err := s.GetTodayGameState(fingerprint)
	if err != nil {
		return nil, err
	}
	if len(state.Guesses) != 0 {
		// get last all match
		var lastSolved []LetterResult
		if len(state.Guesses) != 0 {
			for i := len(state.Guesses) - 1; i >= 0; i-- {
				guess := state.Guesses[i]
				allMatch := true
				for _, lr := range guess {
					if lr.Status != "match" {
						allMatch = false
						break
					}
				}
				if allMatch {
					lastSolved = guess
					break
				}
			}
		}
		if len(lastSolved) >= 4 && len(lastSolved) < 7 {
			solvedLen := len(lastSolved)
			currIndex := solvedLen - 4
			nextIndex := currIndex + 1

			if nextIndex < len(s.puzzleWords) {
				prevWord := ""
				nextWord := s.puzzleWords[nextIndex]

				for _, lr := range lastSolved {
					prevWord += lr.Letter
				}

				letterResult = make([]LetterResult, 0, len(prevWord))
				nextFreq := make(map[rune]int)
				for _, ch := range nextWord {
					nextFreq[ch]++
				}
				for _, lr := range lastSolved {
					ch := rune(lr.Letter[0])

					status := "miss"
					if nextFreq[ch] > 0 {
						status = "present"
					}

					letterResult = append(letterResult, LetterResult{
						Letter: lr.Letter,
						Status: status,
					})
				}
			}
		}
	}

	return letterResult, nil
}

/*
Helpers
*/
func (s *Service) InitializeTodayGameState(fingerprint string) error {
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	_, err := s.initializeGameState.Exec(
		sql.Named("fingerprint", fingerprint),
		sql.Named("date", today),
	)
	if err != nil {
		return fmt.Errorf("InitializeTodayGameState failed: %w", err)
	}

	return nil
}

func (s *Service) GetTodayGameState(fingerprint string) (*GameState, error) {
	// Today's date
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	// Query DB
	var stateJSON string
	err := s.readGameState.QueryRow(
		sql.Named("fingerprint", fingerprint),
		sql.Named("date", today),
	).Scan(&stateJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			// No state found
			return nil, nil
		}
		return nil, fmt.Errorf("GetTodayGameState failed: %w", err)
	}

	// Parse JSON
	var state GameState
	if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
		return nil, fmt.Errorf("invalid game state JSON: %w", err)
	}

	return &state, nil
}

func (s *Service) UpdateGameState(fingerprint string, guessResult []LetterResult) error {
	// Encode the guess as JSON
	guessJSON, err := json.Marshal(guessResult)
	if err != nil {
		return fmt.Errorf("failed to marshal guess result: %w", err)
	}

	// Today’s date (NY timezone)
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	// Execute prepared update statement
	_, err = s.updateGameState.Exec(
		sql.Named("guess", string(guessJSON)),
		sql.Named("fingerprint", fingerprint),
		sql.Named("date", today),
	)
	if err != nil {
		return fmt.Errorf("UpdateGameState failed: %w", err)
	}

	return nil
}

func (s *Service) BuildCurrentPuzzle() error {
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	if s.puzzleDate == today && len(s.puzzleWords) == 4 {
		return nil
	}

	puzzleKeyStr, err := secret.LoadPuzzleSecret()
	if err != nil {
		return fmt.Errorf("failed to load puzzle key: %v", err)
	}
	puzzleKey := []byte(puzzleKeyStr)

	mainSeed := seed.HMACSeed(puzzleKey, today)
	if mainSeed < 0 {
		mainSeed = -mainSeed
	}

	prevMask := int64(0)
	requiredShared := []int{0, 1, 2, 3}
	tierLabels := []string{"", "five", "six", "seven"}

	words := make([]string, 4)

	for i := 0; i < 4; i++ {
		if i > 0 {
			mainSeed = seed.HMACSeed(puzzleKey, fmt.Sprintf("%d:%s", mainSeed, tierLabels[i]))
			if mainSeed < 0 {
				mainSeed = -mainSeed
			}
		}

		var count int64
		err := s.selectWordCount.QueryRow(
			sql.Named("length", i+4),
			sql.Named("required_shared", requiredShared[i]),
			sql.Named("prev_mask", prevMask),
		).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to count words: %w", err)
		}

		wr, err := s.getWordByLength(
			i+4,
			prevMask,
			requiredShared[i],
			mainSeed%count,
		)
		if err != nil {
			return fmt.Errorf("failed to select word: %w", err)
		}

		words[i] = wr.Word
		prevMask = wr.LetterMask
	}

	s.puzzleWords = words
	s.puzzleDate = today

	return nil
}

func prepareSQL(db *sql.DB, path string) (*sql.Stmt, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	stmt, err := db.Prepare(string(bytes))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare %s: %w", path, err)
	}
	return stmt, nil
}

func CountIncorrectGuesses(guesses [][]LetterResult) int {
	incorrect := 0

	for _, guess := range guesses {
		for _, lr := range guess {
			if lr.Status != "match" {
				incorrect++
				break // only count once per guess
			}
		}
	}

	return incorrect
}

func (s *Service) getWordByLength(length int, prevMask int64, requiredShared int, offset int64) (WordResult, error) {
	params := []any{
		sql.Named("length", length),
		sql.Named("prev_mask", prevMask),
		sql.Named("required_shared", requiredShared),
		sql.Named("offset", offset),
	}

	var result WordResult
	err := s.selectWord.QueryRow(params...).Scan(&result.Word, &result.Length, &result.LetterMask)
	if err != nil {
		return WordResult{}, err
	}

	return result, nil
}

func computeGuessResult(guess string, answer string) []LetterResult {
	result := make([]LetterResult, len(guess))
	letterCounts := make(map[byte]int)

	for i := 0; i < len(answer); i++ {
		letterCounts[answer[i]]++
	}

	// First pass: exact matches
	for i := 0; i < len(guess); i++ {
		g := guess[i]
		a := answer[i]
		if g == a {
			result[i] = LetterResult{Letter: string(g), Status: "match"}
			letterCounts[g]--
		}
	}

	// Second pass: present or miss
	for i := 0; i < len(guess); i++ {
		if result[i].Status != "" {
			continue
		}
		g := guess[i]
		if letterCounts[g] > 0 {
			result[i] = LetterResult{Letter: string(g), Status: "present"}
			letterCounts[g]--
		} else {
			result[i] = LetterResult{Letter: string(g), Status: "miss"}
		}
	}

	return result
}

/*
func (s *Service) PrintFullPuzzle() error {
	lengths := []int{4, 5, 6, 7}
	requiredShared := []int{0, 1, 2, 3}
	tierLabels := []string{"", "five", "six", "seven"}

	loc, _ := time.LoadLocation("America/New_York")
	currentTime := time.Now().In(loc).Format("2006-01-02")

	puzzleKeyStr, err := secret.LoadPuzzleSecret()
	if err != nil {
		return fmt.Errorf("failed to load puzzle key: %v", err)
	}
	puzzleKey := []byte(puzzleKeyStr)
	mainSeed := seed.HMACSeed(puzzleKey, currentTime)
	if mainSeed < 0 {
		mainSeed = -mainSeed
	}

	prevMask := int64(0)
	currSeed := mainSeed

	fmt.Println("🧩 Puzzle for", currentTime)
	for i := 0; i < len(lengths); i++ {
		if i > 0 {
			currSeed = seed.HMACSeed(puzzleKey, fmt.Sprintf("%d:%s", currSeed, tierLabels[i]))
			if currSeed < 0 {
				currSeed = -currSeed
			}
		}

		var count int64
		err := s.selectWordCount.QueryRow(
			sql.Named("length", lengths[i]),
			sql.Named("required_shared", requiredShared[i]),
			sql.Named("prev_mask", prevMask),
		).Scan(&count)
		if err != nil {
			log.Fatalf("Failed to count words for tier %d: %v", i, err)
		}

		word, err := s.getWordByLength(
			lengths[i],
			prevMask,
			requiredShared[i],
			currSeed%count,
		)
		if err != nil {
			log.Fatalf("Failed to select word for tier %d: %v", i, err)
		}

		fmt.Printf("Tier %d (%d-letter): %s\n", i+1, word.Length, word.Word)
		prevMask = word.LetterMask
	}

	return nil
}
*/
