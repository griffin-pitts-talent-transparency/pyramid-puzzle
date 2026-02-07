package word_chain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"lebron-games/internal/auth"
	"lebron-games/internal/db"
	"log"
	"net/http"
	"strings"
	"time"
)

type wordChainStartingValues struct {
	FirstWord    string   `json:"first_word"`
	FirstLetters []string `json:"first_letters"`
	WordLengths  []int    `json:"word_lengths"`
}

type wordChainGuessResponse struct {
	Result     bool `json:"result"`
	NextLetter byte `json:"next_letter"`
}

type incomingWordChainGuess struct {
	Guess string `json:"guess"`
	Index int    `json:"index"`
}

type wordChainCache struct {
	Date  string
	Chain []string
}

type wordChainGuess struct {
	Word   string `json:"word"`
	Result bool   `json:"result"`
}

type wordChainGameState struct {
	Guesses       []wordChainGuess `json:"guesses"` // {"word": string, "result": boolean}
	WordLengths   []int            `json:"word_lengths"`
	LetterHints   []string         `json:"letter_hints"`
	TotalAttempts int              `json:"total_attempts"`
}

var CACHED_WORD_CHAIN wordChainCache

/*
Endpoint handlers
*/
func HandleInitializeWordChainGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	cognitoSub := auth.ValidateAndExtractCognitoSub(r)
	if cognitoSub != nil {
		/*
			read user's word chain for the day.
			If empty: initialize and return starting word chain.
			Else, build partial game and return
		*/
		loadTodayWordChain()
		userGuesses := readWordChainGuessesBySub(*cognitoSub)
		if len(userGuesses) == 0 { // also checks for null
			initializeWordChainGameStateBySub(*cognitoSub)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wordChainGameState{
			Guesses:       userGuesses,
			WordLengths:   getWordChainWordLengths(),
			LetterHints:   getAllLetterHints(userGuesses),
			TotalAttempts: countIncorrectGuesses(userGuesses),
		})

	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func HandleGetCurrentWordChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	loadTodayWordChain()
	firstWord := ""
	if len(CACHED_WORD_CHAIN.Chain) > 0 {
		firstWord = CACHED_WORD_CHAIN.Chain[0]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wordChainStartingValues{
		FirstWord:    firstWord,
		FirstLetters: getWordChainFirstLetters(),
		WordLengths:  getWordChainWordLengths(),
	})
}

func HandleValidateGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	if !auth.ValidateCognitoUser(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	result := validateGuess(r)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

/*
Endpoint handler delegators
*/
func validateGuess(r *http.Request) wordChainGuessResponse {
	cognitoSub := auth.ValidateAndExtractCognitoSub(r)
	var matchResult bool
	var nextLetter byte
	if cognitoSub != nil {
		loadTodayWordChain()
		var req incomingWordChainGuess
		_ = json.NewDecoder(r.Body).Decode(&req)

		guess := strings.ToLower(req.Guess)

		matchResult = len(CACHED_WORD_CHAIN.Chain) > req.Index &&
			strings.ToLower(CACHED_WORD_CHAIN.Chain[req.Index]) == guess

		appendWordChainGuessBySub(*cognitoSub, guess, matchResult)
		if !matchResult {
			fmt.Println("Match result negative; calculating next hint")
			userGuesses := readWordChainGuessesBySub(*cognitoSub)
			nextLetter = getNextLetterHint(userGuesses, req.Index)
		}
	}
	return wordChainGuessResponse{
		Result:     matchResult,
		NextLetter: nextLetter,
	}
}

/*
Letter hint helpers
*/
func getNextLetterHint(guesses []wordChainGuess, submittedIndex int) byte {
	var result byte
	wordIndex := 1
	guessIndex := 0

	// Consume guesses for all prior words
	for wordIndex < submittedIndex && guessIndex < len(guesses) {
		word := CACHED_WORD_CHAIN.Chain[wordIndex]
		failCount := 0

		for guessIndex < len(guesses) {
			if guesses[guessIndex].Result {
				guessIndex++
				break
			}
			failCount++
			guessIndex++
			if failCount >= len(word)-1 {
				break
			}
		}
		wordIndex++
	}

	word := CACHED_WORD_CHAIN.Chain[submittedIndex]
	failCount := 0
	solved := false

	// IMPORTANT: always inspect at least the first guess for this word
	if guessIndex < len(guesses) && guesses[guessIndex].Result {
		solved = true
	} else {
		for guessIndex < len(guesses) {
			if guesses[guessIndex].Result {
				solved = true
				break
			}
			failCount++
			guessIndex++
			if failCount >= len(word)-1 {
				break
			}
		}
	}

	if !solved && failCount < len(word) {
		result = word[failCount]
	}

	return result
}

func getAllLetterHints(guesses []wordChainGuess) []string {
	results := make([]string, len(CACHED_WORD_CHAIN.Chain))

	// Given word
	results[0] = CACHED_WORD_CHAIN.Chain[0]

	// Default: first letter for all others
	for i := 1; i < len(CACHED_WORD_CHAIN.Chain); i++ {
		results[i] = string(CACHED_WORD_CHAIN.Chain[i][0])
	}

	failCount := 0
	wordIndex := 1

	for i := 0; i < len(guesses) && wordIndex < len(CACHED_WORD_CHAIN.Chain); i++ {
		word := CACHED_WORD_CHAIN.Chain[wordIndex]
		maxFails := len(word) - 1
		g := guesses[i]

		if !g.Result {
			failCount++
		}

		segmentComplete := g.Result || failCount == maxFails

		if segmentComplete {
			n := failCount + 1
			if n > len(word) {
				n = len(word)
			}

			results[wordIndex] = word[:n]

			failCount = 0
			wordIndex++
		}
	}

	// Current in-progress word
	if wordIndex < len(CACHED_WORD_CHAIN.Chain) && failCount > 0 {
		word := CACHED_WORD_CHAIN.Chain[wordIndex]

		n := failCount + 1
		if n > len(word) {
			n = len(word)
		}

		results[wordIndex] = word[:n]
	}

	return results
}

/*
Misc helpers
*/
func getWordChainFirstLetters() []string {
	result := []string{}

	if len(CACHED_WORD_CHAIN.Chain) > 0 {
		result = make([]string, len(CACHED_WORD_CHAIN.Chain))
		for i, word := range CACHED_WORD_CHAIN.Chain {
			if len(word) > 0 {
				result[i] = string(word[0])
			}
		}
	}

	return result
}

func getWordChainWordLengths() []int {
	result := []int{}

	if len(CACHED_WORD_CHAIN.Chain) > 0 {
		result = make([]int, len(CACHED_WORD_CHAIN.Chain))
		for i, word := range CACHED_WORD_CHAIN.Chain {
			result[i] = len(word)
		}
	}

	return result
}

func countIncorrectGuesses(guesses []wordChainGuess) int {
	var count int
	for _, g := range guesses {
		if !g.Result {
			count++
		}
	}
	return count
}

/*
DB helpers
*/
func readWordChainGuessesBySub(cognitoSub string) []wordChainGuess {
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	var rawGuesses []map[string]interface{}
	var guesses []wordChainGuess

	userdataDB, err := db.OpenSQLiteDBByName("userdata")
	if err == nil {
		defer userdataDB.Close()

		stmt, err := db.PrepareSQLStatementByName(
			userdataDB,
			"read_word_chain_game_state_by_sub.sql",
			"userdata",
		)
		if err == nil {
			defer stmt.Close()

			var jsonBlob string
			err = stmt.QueryRow(
				sql.Named("cognito_sub", cognitoSub),
				sql.Named("puzzle_date", today),
			).Scan(&jsonBlob)
			if err == nil {
				if err := json.Unmarshal([]byte(jsonBlob), &rawGuesses); err != nil {
					log.Println("❌ Failed to unmarshal raw guesses:", err)
				} else {
					for _, raw := range rawGuesses {
						word, _ := raw["word"].(string)
						var result bool
						switch v := raw["result"].(type) {
						case float64:
							result = v != 0
						case bool:
							result = v
						}
						guesses = append(guesses, wordChainGuess{
							Word:   word,
							Result: result,
						})
					}
				}
			}
		}
	}

	if len(guesses) == 0 {
		return nil
	}
	return guesses
}

func initializeWordChainGameStateBySub(cognitoSub string) {
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	userdataDB, err := db.OpenSQLiteDBByName("userdata")
	if err == nil {
		defer userdataDB.Close()

		stmt, err := db.PrepareSQLStatementByName(
			userdataDB,
			"initialize_word_chain_game_state_by_sub.sql",
			"userdata",
		)
		if err == nil {
			defer stmt.Close()
			_, _ = stmt.Exec(
				sql.Named("cognito_sub", cognitoSub),
				sql.Named("puzzle_date", today),
			)
		}
	}
}

func loadTodayWordChain() {
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	if CACHED_WORD_CHAIN.Date != today || CACHED_WORD_CHAIN.Chain == nil {
		wordChainDB, err := db.OpenSQLiteDBByName("word_chain")
		if err != nil {
			log.Fatalf("❌ Failed to open word_chain DB: %v", err)
		}
		defer wordChainDB.Close()

		stmt, err := db.PrepareSQLStatementByName(wordChainDB, "read_word_chain.sql", "word_chain")
		if err == nil {
			defer stmt.Close()

			var jsonBlob string
			err = stmt.QueryRow(today).Scan(&jsonBlob)
			if err == nil {
				var chain []string
				err = json.Unmarshal([]byte(jsonBlob), &chain)
				if err == nil {
					CACHED_WORD_CHAIN = wordChainCache{
						Date:  today,
						Chain: chain,
					}
				}
			}
		}
	}
}

func appendWordChainGuessBySub(cognitoSub string, word string, result bool) {
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	userdataDB, err := db.OpenSQLiteDBByName("userdata")
	if err == nil {
		defer userdataDB.Close()

		stmt, err := db.PrepareSQLStatementByName(
			userdataDB,
			"update_word_chain_game_state_by_sub.sql",
			"userdata",
		)
		if err == nil {
			defer stmt.Close()
			_, _ = stmt.Exec(
				sql.Named("word", word),
				sql.Named("result", result),
				sql.Named("cognito_sub", cognitoSub),
				sql.Named("puzzle_date", today),
			)

		}
	}
}
