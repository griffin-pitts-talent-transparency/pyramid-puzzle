package word_chain

import "fmt"

func setHardcodedWordChain() {
	CACHED_WORD_CHAIN = wordChainCache{
		Chain: []string{
			"brain",
			"storm",
			"drain",
			"trap",
			"house",
			"cat",
			"rescue",
			"dog",
			"park",
			"bench",
			"press",
		},
	}
}

func manualTest_getNextHintLetterNew() {
	tests := []struct {
		name           string
		guesses        []wordChainGuess
		submittedIndex int
		expected       rune
	}{
		{
			name: "test 1",
			guesses: []wordChainGuess{
				{Word: "xxxxx", Result: false},
			},
			submittedIndex: 1,
			expected:       't', // fill in
		},
		{
			name: "test 2",
			guesses: []wordChainGuess{
				{Word: "xxxxx", Result: false},
				{Word: "yyyyy", Result: false},
			},
			submittedIndex: 1,
			expected:       'o', // fill in
		},
		{
			name: "test 3",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},
				{Word: "aaaaa", Result: false},
			},
			submittedIndex: 2,
			expected:       'r', // fill in
		},
		{
			name: "test 4",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},
				{Word: "drain", Result: true},
				{Word: "aaaa", Result: false},
				{Word: "bbbb", Result: false},
			},
			submittedIndex: 3,
			expected:       'a', // fill in
		},
		{
			name: "test 5",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},
				{Word: "drain", Result: true},
				{Word: "aaaa", Result: false},
				{Word: "bbbb", Result: false},
				{Word: "bbbb", Result: false},
				{Word: "bbbb", Result: false},
			},
			submittedIndex: 4,
			expected:       'o', // fill in
		},
		{
			name: "test 6",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false}, // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false}, // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false}, // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false}, // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false}, // index[4] house (guess 4) - return e
				{Word: "bbbb", Result: false}, // index[5] cat (guess 1) - return a
			},
			submittedIndex: 5,
			expected:       'a', // fill in
		},
		{
			name: "test 7",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false}, // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false}, // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false}, // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false}, // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false}, // index[4] house (guess 4) - return e
				{Word: "bbbb", Result: true},  // index[5] cat (guess 1) - return 0
			},
			submittedIndex: 5,
			expected:       '0', // fill in
		},
		{
			name: "test 8",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return e
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 0
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 0
			},
			submittedIndex: 6,
			expected:       '0', // fill in
		},
		{
			name: "test 9",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return e
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 0
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 0
				{Word: "ppp", Result: false},   // index[7] dog (guess 1) - return o
				{Word: "ppp", Result: false},   // index[7] dog (guess 2) - return g
			},
			submittedIndex: 7,
			expected:       'g', // fill in
		},
		{
			name: "test 10",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return e
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 0
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 0
				{Word: "ppp", Result: false},   // index[7] dog (guess 1) - return o
				{Word: "ppp", Result: false},   // index[7] dog (guess 2) - return g
				{Word: "pppp", Result: false},  // index[8] park (guess 1) - return a
			},
			submittedIndex: 8,
			expected:       'a', // fill in
		},
		{
			name: "test 10",
			guesses: []wordChainGuess{
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 1) - return t
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 2) - return o
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 3) - return r
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 4) - return m
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 1) - return r
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 2) - return a
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 3) - return i
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 4) - return n
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "aaaa", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "aaaa", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "aaaaa", Result: false}, // index[4] house (guess 1) - return o
			},
			submittedIndex: 4,
			expected:       'o', // fill in
		},
		{
			name: "test 11",
			guesses: []wordChainGuess{
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 1) - return t
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 2) - return o
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 3) - return r
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 4) - return m
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 1) - return r
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 2) - return a
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 3) - return i
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 4) - return n
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "aaaa", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "aaaa", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "aaaaa", Result: false}, // index[4] house (guess 1) - return o
				{Word: "aaaaa", Result: false}, // index[4] house (guess 2) - return u
			},
			submittedIndex: 4,
			expected:       'u', // fill in
		},
	}

	for _, t := range tests {
		result := testing_getNextHintLetterNew(t.guesses, t.submittedIndex)
		expectedByte := expectedToByte(t.expected)

		pass := "❌no"
		if result == expectedByte {
			pass = "✅yes"
		}

		fmt.Println(t.name)
		fmt.Printf("expected: %c\n", t.expected)
		if result == 0 {
			fmt.Println("result:   0")
		} else {
			fmt.Printf("result:   %c\n", result)
		}
		fmt.Printf("Pass result: %s\n\n", pass)
	}
}

func expectedToByte(expected rune) byte {
	if expected == '0' {
		return 0
	}
	return byte(expected)
}

func testing_getNextHintLetterNew(guesses []wordChainGuess, submittedIndex int) byte {
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

func manualTest_getNextHintLettersNew() {
	tests := []struct {
		name     string
		guesses  []wordChainGuess
		expected []string
	}{
		{
			name:     "test 1",
			guesses:  []wordChainGuess{},
			expected: []string{"brain", "s", "d", "t", "h", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 2",
			guesses: []wordChainGuess{
				{Word: "xxxxx", Result: false},
				{Word: "yyyyy", Result: false},
			},
			expected: []string{"brain", "sto", "d", "t", "h", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 3",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},
				{Word: "aaaaa", Result: false},
			},
			expected: []string{"brain", "s", "dr", "t", "h", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 4",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},
				{Word: "drain", Result: true},
				{Word: "aaaa", Result: false},
				{Word: "bbbb", Result: false},
			},
			expected: []string{"brain", "s", "d", "tra", "h", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 5",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},
				{Word: "drain", Result: true},
				{Word: "aaaa", Result: false},
				{Word: "bbbb", Result: false},
				{Word: "bbbb", Result: false},
				{Word: "bbbb", Result: false},
			},
			expected: []string{"brain", "s", "d", "trap", "ho", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 6",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false}, // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false}, // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false}, // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false}, // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false}, // index[4] house (guess 4) - return e
				{Word: "bbbb", Result: false}, // index[5] cat (guess 1) - return a
			},
			expected: []string{"brain", "s", "d", "trap", "house", "ca", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 7",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false}, // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false}, // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false}, // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false}, // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false}, // index[4] house (guess 4) - return e
				{Word: "bbbb", Result: true},  // index[5] cat (guess 1) - return 0
			},
			expected: []string{"brain", "s", "d", "trap", "house", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 8",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return e
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 0
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 0
			},
			expected: []string{"brain", "s", "d", "trap", "house", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 9",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return e
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 0
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 0
				{Word: "ppp", Result: false},   // index[7] dog (guess 1) - return o
				{Word: "ppp", Result: false},   // index[7] dog (guess 2) - return g
			},
			expected: []string{"brain", "s", "d", "trap", "house", "c", "r", "dog", "p", "b", "p"},
		},
		{
			name: "test 10",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 0
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 0
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return o
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return u
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return s
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return e
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 0
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 0
				{Word: "ppp", Result: false},   // index[7] dog (guess 1) - return o
				{Word: "ppp", Result: false},   // index[7] dog (guess 2) - return g
				{Word: "pppp", Result: false},  // index[8] park (guess 1) - return a
			},
			expected: []string{"brain", "s", "d", "trap", "house", "c", "r", "dog", "pa", "b", "p"},
		},
		{
			name: "test 11",
			guesses: []wordChainGuess{
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 1) - return t
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 2) - return o
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 3) - return r
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 4) - return m
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 1) - return r
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 2) - return a
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 3) - return i
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 4) - return n
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "aaaa", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "aaaa", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "aaaaa", Result: false}, // index[4] house (guess 1) - return o
			},
			expected: []string{"brain", "storm", "drain", "trap", "ho", "c", "r", "d", "p", "b", "p"},
		},
		{
			name: "test 12",
			guesses: []wordChainGuess{
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 1) - return t
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 2) - return o
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 3) - return r
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 4) - return m
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 1) - return r
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 2) - return a
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 3) - return i
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 4) - return n
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return r
				{Word: "aaaa", Result: false},  // index[3] trap (guess 2) - return a
				{Word: "aaaa", Result: false},  // index[3] trap (guess 3) - return p
				{Word: "aaaaa", Result: false}, // index[4] house (guess 1) - return o
				{Word: "aaaaa", Result: false}, // index[4] house (guess 2) - return u
			},
			expected: []string{"brain", "storm", "drain", "trap", "hou", "c", "r", "d", "p", "b", "p"},
		},
	}

	for _, t := range tests {
		result := testing_getNextHintLettersNew(t.guesses)
		matchResult := matchStringArrays(t.expected, result)
		pass := "❌no"
		if matchResult {
			pass = "✅yes"
		}

		fmt.Println(t.name)
		fmt.Printf("expected: %v\n", t.expected)
		fmt.Printf("result:   %v\n", result)
		fmt.Printf("Pass: %s\n\n", pass)
	}
}

func matchStringArrays(arrOne []string, arrTwo []string) bool {
	if len(arrOne) != len(arrTwo) {
		return false
	}
	for i := 0; i < len(arrOne); i++ {
		if arrOne[i] != arrTwo[i] {
			return false
		}
	}
	return true
}

func testing_getNextHintLettersNew(guesses []wordChainGuess) []string {
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

func manualTest_getActiveWordIndex() {
	tests := []struct {
		name     string
		guesses  []wordChainGuess
		expected int
	}{
		{
			name: "test 1",
			guesses: []wordChainGuess{
				{Word: "xxxxx", Result: false}, // index[1] storm (guess 1) - return 1
			},
			expected: 1, // fill in
		},
		{
			name: "test 2",
			guesses: []wordChainGuess{
				{Word: "xxxxx", Result: false}, // index[1] storm (guess 1) - return 1
				{Word: "yyyyy", Result: false}, // index[1] storm (guess 2) - return 1
			},
			expected: 1, // fill in
		},
		{
			name: "test 3",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 2
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 1) - return 2
			},
			expected: 2, // fill in
		},
		{
			name: "test 4",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 2
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 3
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return 3
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return 3
			},
			expected: 3, // fill in
		},
		{
			name: "test 5",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 2
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 3
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return 3
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return 3
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return 3
				{Word: "bbbb", Result: false}, // index[4] trap (guess 2) - return 4
			},
			expected: 4, // fill in
		},
		{
			name: "test 6",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 2
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 3
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return 3
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return 3
				{Word: "bbbb", Result: false}, // index[3] trap (guess 3) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 1) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 2) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 3) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 4) - return 5
				{Word: "bbbb", Result: false}, // index[5] cat (guess 1) - return 5
			},
			expected: 5, // fill in
		},
		{
			name: "test 7",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true}, // index[1] storm (guess 1) - return 2
				{Word: "drain", Result: true}, // index[2] drain (guess 1) - return 3
				{Word: "aaaa", Result: false}, // index[3] trap (guess 1) - return 3
				{Word: "bbbb", Result: false}, // index[3] trap (guess 2) - return 3
				{Word: "bbbb", Result: false}, // index[3] trap (guess 3) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 1) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 2) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 3) - return 4
				{Word: "bbbb", Result: false}, // index[4] house (guess 4) - return 5
				{Word: "bbbb", Result: true},  // index[5] cat (guess 1) - return 6
			},
			expected: 6, // fill in
		},
		{
			name: "test 8",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 2
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 3
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return 3
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return 3
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return 4 ho
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return 4 hou
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return 4 hous
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return 5 house
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 6
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 7
			},
			expected: 7, // fill in
		},
		{
			name: "test 9",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 2
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 3
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return 3
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return 3
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return 5
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 6
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 7 do
				{Word: "ppp", Result: false},   // index[7] dog (guess 1) - return 7 dog
				{Word: "ppp", Result: false},   // index[7] dog (guess 2) - return 8
			},
			expected: 8, // fill in
		},
		{
			name: "test 10",
			guesses: []wordChainGuess{
				{Word: "storm", Result: true},  // index[1] storm (guess 1) - return 2
				{Word: "drain", Result: true},  // index[2] drain (guess 1) - return 3
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return 3
				{Word: "bbbb", Result: false},  // index[3] trap (guess 2) - return 3
				{Word: "bbbb", Result: false},  // index[3] trap (guess 3) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 1) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 2) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 3) - return 4
				{Word: "bbbb", Result: false},  // index[4] house (guess 4) - return 5
				{Word: "cat", Result: true},    // index[5] cat (guess 1) - return 6
				{Word: "rescue", Result: true}, // index[6] rescue (guess 1) - return 7
				{Word: "ppp", Result: false},   // index[7] dog (guess 1) - return 7
				{Word: "ppp", Result: false},   // index[7] dog (guess 2) - return 8
				{Word: "pppp", Result: false},  // index[8] park (guess 1) - return 8
			},
			expected: 8, // fill in
		},
		{
			name: "test 10",
			guesses: []wordChainGuess{
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 1) - return 1 st
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 2) - return 1 sto
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 3) - return 1 stor
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 4) - return 2 storm
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 1) - return 2 dr
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 2) - return 2 dra
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 3) - return 2 drai
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 4) - return 3 drain
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return 3 tr
				{Word: "aaaa", Result: false},  // index[3] trap (guess 2) - return 3 tra
				{Word: "aaaa", Result: false},  // index[3] trap (guess 3) - return 3 trap
				{Word: "aaaaa", Result: false}, // index[4] house (guess 1) - return 4 ho
			},
			expected: 4, // fill in
		},
		{
			name: "test 11",
			guesses: []wordChainGuess{
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 1) - return 1 st
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 2) - return 1 sto
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 3) - return 1 stor
				{Word: "aaaaa", Result: false}, // index[1] storm (guess 4) - return 2 storm
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 1) - return 2 dr
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 2) - return 2 dra
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 3) - return 2 drai
				{Word: "aaaaa", Result: false}, // index[2] drain (guess 4) - return 3 drain
				{Word: "aaaa", Result: false},  // index[3] trap (guess 1) - return 3
				{Word: "aaaa", Result: false},  // index[3] trap (guess 2) - return 3
				{Word: "aaaa", Result: false},  // index[3] trap (guess 3) - return 4
				{Word: "aaaaa", Result: false}, // index[4] house (guess 1) - return 4
				{Word: "aaaaa", Result: false}, // index[4] house (guess 2) - return 4
			},
			expected: 4, // fill in
		},
	}

	for _, t := range tests {
		result := test_getActiveWordIndex(t.guesses)

		pass := "❌no"
		if result == t.expected {
			pass = "✅yes"
		}

		fmt.Printf("---%s---\n", t.name)
		fmt.Printf("Expected: %d\n", t.expected)
		fmt.Printf("Result:   %d\n", result)
		fmt.Printf("Passed:   %s\n", pass)
	}
}

func test_getActiveWordIndex(guesses []wordChainGuess) int {
	result := 1
	var currentChainWordLength int
	var consecutiveWrongGuesses int
	chainWordIndex := 1
	for i := 0; i < len(guesses); i++ {
		guessResult := guesses[i].Result
		currentChainWordLength = len(CACHED_WORD_CHAIN.Chain[chainWordIndex])
		if guessResult {
			result++
			chainWordIndex++
			consecutiveWrongGuesses = 0
		} else {
			consecutiveWrongGuesses++
			if consecutiveWrongGuesses == (currentChainWordLength - 1) {
				result++
				consecutiveWrongGuesses = 0
				chainWordIndex++
			}
		}
	}

	return result
}
