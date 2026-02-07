import { getUserAsync } from "./auth.js"; 

let CURRENT_WORD_INDEX = 0;
const NUM_WORDS = 11;

function renderLetterPush(wordIndex, letter, type = "") {
    const tiles = document.querySelectorAll(`#word-${wordIndex} .word-chain-tile`);
    for (const tile of tiles) {
        if (tile.textContent === "") {
            tile.textContent = letter;
            if (type != "") {
                tile.classList.add(type);
            }
            return;
        }
    }
}

function renderLetterPop(wordIndex) {
    const tiles = document.querySelectorAll(`#word-${wordIndex} .word-chain-tile`);
    for (let i = tiles.length - 1; i >= 0; i--) {
        if(!tiles[i].classList.contains("word-chain-hint-letter") && tiles[i].textContent !== "") {
            tiles[i].textContent = "";
            tiles[i].classList.remove("word-chain-hint-letter", "word-chain-solved-letter");
            return;
        }
    }
}

function renderCurrentWord() {
    const allWords = document.querySelectorAll(".word-chain-current-word");
    for (let i = 0; i < allWords.length; i++) {
        allWords[i].classList.remove("word-chain-current-word");
    }

    const current = document.getElementById(`word-${CURRENT_WORD_INDEX}`);
    if (current) {
        current.classList.add("word-chain-current-word");
    }
}

function renderCorrectWord() {
    const activeWord = document.querySelector(".word-chain-current-word");
    if (activeWord) {
        activeWord.classList.remove("word-chain-current-word");
        activeWord.classList.add("word-chain-correct-word");
    }
}

function renderShowAuthGate() {
    document.getElementById("authGate")?.classList.remove("hidden");
    document.getElementById("wordChainGameContainer")?.classList.add("hidden");
}

function renderHideAuthGate() {
    document.getElementById("authGate")?.classList.add("hidden");
    document.getElementById("wordChainGameContainer")?.classList.remove("hidden");
}

async function getAndRequireAuthenticatedUserAsync() {
    let user = await getUserAsync();

    if (!user || user.expired || !user.id_token) {
        renderShowAuthGate();
        user = null;
    } else {
        renderHideAuthGate();
    }

    return user;
}

function initializeWordTiles(wordIndex, length) {
    const container = document.getElementById(`word-${wordIndex}`);
    container.innerHTML = "";

    for (let i = 0; i < length; i++) {
        const tile = document.createElement("div");
        tile.className = "word-chain-tile";
        container.appendChild(tile);
    }
}

function extractGuess() {
    let result = "";
    const tiles = document.querySelectorAll(`#word-${CURRENT_WORD_INDEX} .word-chain-tile`);
    for (const tile of tiles) {
        result += tile.textContent.toLowerCase();
    }
    return result;
}

async function handleGuessSubmitted() {
    const activeRow = document.querySelectorAll(`#word-${CURRENT_WORD_INDEX} .word-chain-tile`);
    const playerGuess = extractGuess();
    console.log("SUBMITTED PLAYER GUESS IS: " + playerGuess);
    console.log("CURRENT WORD INDEX IS: " + CURRENT_WORD_INDEX);
    if (playerGuess.length === activeRow.length) {
        const user = await getAndRequireAuthenticatedUserAsync();
        if (user && !user.expired && user.id_token) {
            try {
                const res = await fetch("/api/v1.0/word-chain-validate-guess", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${user.id_token}`
                    },
                    body: JSON.stringify({
                        guess: playerGuess,
                        index: CURRENT_WORD_INDEX
                    })
                });

                const data = await res.json();
                console.log("GUESS VALIDATED; JSON:");
                console.log(JSON.stringify(data));
                if(data) {
                    if (data.result) {
                        console.log("CORRECT GUESS - RENDERING CORRECT GUESS");
                        renderCorrectWord();
                        if(CURRENT_WORD_INDEX === NUM_WORDS) {
                            // win
                            console.log("win");
                        } else {
                            CURRENT_WORD_INDEX++;
                            // renderActiveRow();
                        }
                    } else {
                        console.log("INCORRECT GUESS - RENDERING INCORRECT GUESS");
                        for(let i = 0; i < playerGuess.length; i++) {
                            renderLetterPop(CURRENT_WORD_INDEX);
                        }
                        const nextLetter = String.fromCharCode(data.next_letter);
                        if (/^[a-z]$/.test(nextLetter)) {
                            renderLetterPush(CURRENT_WORD_INDEX, nextLetter, "word-chain-hint-letter");
                            // handle the case where the entire word is now hint letters
                            const tiles = document.querySelectorAll(`#word-${CURRENT_WORD_INDEX} .word-chain-tile`);
                            const allHint = Array.from(tiles).every(tile =>
                                tile.classList.contains("word-chain-hint-letter") && tile.textContent.trim() !== "");

                            if (allHint) {
                                console.log("ALL HINT LETTERS - ADVANCING WORD INDEX");
                                CURRENT_WORD_INDEX++;
                                if (CURRENT_WORD_INDEX === NUM_WORDS) {
                                    console.log("win");
                                } else {
                                    // renderActiveRow();
                                }
                            }
                        } else {
                            console.warn("Invalid hint letter received, skipping render:", data.next_letter);
                        }
                    }
                } else {
                    // TODO: handle where data is null
                    console.log("ERROR - VALIDATE GUESS DATA IS NULL");
                }
            } catch {
            }
            if(CURRENT_WORD_INDEX < NUM_WORDS) {
                renderCurrentWord();
            }
        }
    }
}

let lastKeyTime = 0;
async function safeHandleKeyPress(key) {
    const now = performance.now();
    if (now - lastKeyTime < 25) return; // prevent ghost double input
    lastKeyTime = now;
    await handleKeyPress(key);
}

async function handleKeyPress(key) {
    if (/^[a-zA-Z]$/.test(String(key))) {
        renderLetterPush(CURRENT_WORD_INDEX, key);
    } else if (key === '⌫') {
        renderLetterPop(CURRENT_WORD_INDEX);
    }  else if (key === 'ENTER') {
        if(CURRENT_WORD_INDEX < NUM_WORDS) {
            handleGuessSubmitted();
        }
    }
}

function initializeKeyboard() {
    const keys = document.querySelectorAll('.keyboard .key');

    keys.forEach((key) => {
        key.addEventListener('click', async () => {
            const keyValue = key.textContent.trim();
            await safeHandleKeyPress(keyValue);
        });
    });

    // desktop keyboard
    window.addEventListener('keydown', (e) => {
        let key = e.key;

        if (key === 'Enter') key = 'ENTER';
        else if (key === 'Backspace') key = '⌫';
        else key = key.toUpperCase();

        safeHandleKeyPress(key);
    });
}

function initializeGame(gameState) {
    if (gameState) {
        const solvedGuesses = Array.isArray(gameState.guesses)
            ? gameState.guesses.filter(g => g.result === true).map(g => g.word) : [];

        CURRENT_WORD_INDEX = 0;
        const totalWords = gameState.word_lengths.length;
        for (let i = 0; i < totalWords; i++) {
            let wordLength = gameState.word_lengths[i];
            initializeWordTiles(i, wordLength);
            let letterHint = gameState.letter_hints[i];
            for (let j = 0; j < wordLength; j++) {
                let letter = letterHint[j];
                if(letter) {
                    renderLetterPush(i, letterHint[j], "word-chain-hint-letter");
                }
            }
        }
        let solvedGuessIndex = 0;
        for (let i = 0; i < totalWords; i++) {
            if(solvedGuessIndex < solvedGuesses.length) {
                let wordLength = gameState.word_lengths[i];
                let hintTilesLength = (document.querySelectorAll(`#word-${i} .word-chain-tile.word-chain-hint-letter`)).length;
                if(wordLength > hintTilesLength) {
                    // check if there is a solvedGuess to push. If so, push the remaining letters.
                    // This will be a substring of the hint letter(s). If the word is 'apple', the hint
                    // letters might be 'ap', and the solved word will be 'apple' entirely, so you need substring 'pple'
                    // and push them in order
                    let solvedWord = solvedGuesses[solvedGuessIndex];
                    for (let j = hintTilesLength; j < solvedWord.length; j++) {
                        renderLetterPush(i, solvedWord[j], "word-chain-solved-letter");
                    }
                    solvedGuessIndex++;
                }
            }
        }
        solvedGuessIndex = 0;
        for (let i = 0; i < totalWords; i++) {
            let wordLength = gameState.word_lengths[i];
            let hintTilesLength = (document.querySelectorAll(`#word-${i} .word-chain-tile.word-chain-hint-letter`)).length;
            if(hintTilesLength === wordLength) {
                CURRENT_WORD_INDEX++;
            } else if(solvedGuessIndex < solvedGuesses.length) {
                CURRENT_WORD_INDEX++;
                solvedGuessIndex++;
            }

        }
        // renderActiveRow();
        renderCurrentWord();
    }
}

/*
    Main
*/
document.addEventListener("DOMContentLoaded", async () => {
    // await getCurrentWordChain();
    initializeKeyboard();

    const user = await getAndRequireAuthenticatedUserAsync();
    if (user && !user.expired && user.id_token) {
        try {
            const res = await fetch("/api/v1.0/word-chain-initialize-game", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${user.id_token}`
                }
            });

            if (res.ok) {
                const gameState = await res.json();
                console.log(JSON.stringify(gameState));
                initializeGame(gameState);
            }
        } catch (err) {
            console.error("Initialization failed:", err);
        }
    }

    const signInButton = document.getElementById("authGateLoginButton");
    if (signInButton) {
        signInButton.addEventListener("click", () => {
            sessionStorage.setItem(
                "post_signin_redirect",
                window.location.href
            );
            window.location.href = "/signin";
        });
    }
});