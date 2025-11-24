// full_game.js
import {
	initializeUser,
    readFingerprint,
	loadTodayGameState
} from './storage.js';

import { loadPartial } from './loadPartials.js';

import {
    initializeInstructionsModal,
    renderActiveTiles,
    setActiveRow,
    renderGuessResult, 
    addBlankRow, 
    renderGuessCounter, 
    colorKeyboardKeys, 
    clearKeyboardColors, 
    showEndMessage,
    buildEmojiResultsFromDOM,
    showToast
} from './render.js';

let currentGuess = []; // For in-progress letter input

async function initializeGameSession() {
	await initializeUser();

	const fingerprint = readFingerprint();

	if (fingerprint) {
		const loadedState = await loadTodayGameState(fingerprint);
		if (loadedState) {
            let gameStatus = "playing";
			for (const guessResult of loadedState.guesses) {
                let count = 0;
                let isCorrect = true;
                if (guessResult && Array.isArray(guessResult)) {
                    count = guessResult.length;
                    renderGuessResult(guessResult);
            
                    // determine if this guess is fully correct
                    for (let i = 0; i < guessResult.length; i++) {
                        const letterObj = guessResult[i];
                        if (letterObj.status !== "match") {
                            isCorrect = false;
                        }
                    }
                } else {
                    isCorrect = false;
                }
            
                colorKeyboardKeys(guessResult);
                if (isCorrect && count === 7) {
                    gameStatus = "won";
                    // showEndMessage("Game over; you WIN!");
                } else if(!isCorrect) {
                    // Only add row if guess is NOT correct
                    addBlankRow(count);
                } else {
                    clearKeyboardColors();
                    if (count < 7) {
                        const resp = await fetch("/api/v1.0/get-next-word-hint", {
                            method: "POST",
                            headers: { "Content-Type": "application/json" },
                            body: JSON.stringify({
                                fingerprint: fingerprint
                            })
                        });
                        const data = await resp.json();
                        colorKeyboardKeys(data.hint);
                    }
                }
            }
            
            if(gameStatus === "playing") {
                const countResp = await fetch("/api/v1.0/get-today-remaining-guesses", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        fingerprint: fingerprint
                    })
                });
                
                const countData = await countResp.json();
    
                if(countData.remainingGuesses > 0) {
                    // player still has guesses remaining
                    renderGuessCounter(countData.remainingGuesses);
                    setActiveRow();
                } else {
                    // gameStatus = "lost";
                    showEndMessage("Game over; you lost...");
                }

            } else if (gameStatus === "won") {
                showEndMessage("Game over; you WON!");
            }
		}
	}
}

async function handleGuess(submittedGuess) {
    const fingerprint = readFingerprint();
    if (fingerprint) {
        const res = await fetch('/api/v1.0/validate-guess', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                fingerprint,
                guess: submittedGuess
            })
        });
    
        if (res.ok) {
            const data = await res.json();

            colorKeyboardKeys(data.guessResult);
            renderGuessResult(data.guessResult);
            renderGuessCounter(data.remainingGuesses);
            
            if(data.remainingGuesses > 0) {
                let count = data.guessResult.length;;
                let isCorrect = true;
                // determine if this guess is fully correct
                for (let i = 0; i < data.guessResult.length; i++) {
                    const letterObj = data.guessResult[i];
                    if (letterObj.status !== "match") {
                        isCorrect = false;
                    }
                }

                if(isCorrect && count === 7) {
                    showEndMessage("Game over; you WON!");
                } else if(isCorrect) {
                    count++;
                    clearKeyboardColors();
                    if (count <= 7) {
                        const resp = await fetch("/api/v1.0/get-next-word-hint", {
                            method: "POST",
                            headers: { "Content-Type": "application/json" },
                            body: JSON.stringify({
                                fingerprint: fingerprint
                            })
                        });
                        const hintResult = await resp.json();
                        colorKeyboardKeys(hintResult.hint);
                        setActiveRow();
                    }
                } else if(count <= 7) {
                    addBlankRow(count);
                    setActiveRow();
                }
            } else {
                showEndMessage("Game over; you lost...");
            }
        }
    }
}

async function handleKeyPress(key) {
    // const activeRow = document.querySelector('.row.active');
    const activeRow = document.querySelector('.row.active-row');
    if (activeRow) {
        const expectedLength = activeRow.querySelectorAll('.tile').length;
        if (/^[a-zA-Z]$/.test(String(key))) {
            if (currentGuess.length < expectedLength) {
                currentGuess.push(key);
                renderActiveTiles(currentGuess);
            }
        } else if (key === '⌫') {
            if (currentGuess.length > 0) {
                currentGuess.pop();
                renderActiveTiles(currentGuess);
            }
        }  else if (key === 'ENTER') {
            if (currentGuess.length === expectedLength) {
                const submitted = currentGuess.join("").toLowerCase().split("");
                currentGuess = []; 
                handleGuess(submitted)
            }
        }
    }
}

function initializeKeyboard() {
    const keys = document.querySelectorAll('.keyboard .key');

    keys.forEach((key) => {
            key.addEventListener('click', async () => {
            const keyValue = key.textContent.trim();
            await handleKeyPress(keyValue);
        });
    });
}

/*
    main
*/
window.addEventListener('DOMContentLoaded', async () => {
    await loadPartial('game-container', '/partials/game.html');
    
    // Initialization
    await initializeGameSession();
    initializeInstructionsModal();
    initializeKeyboard();
    // renderGuessCounter(gameState);

    
    window.addEventListener('keydown', (e) => {
        let key = e.key;
    
        // Normalize desktop keyboard inputs
        if (key === 'Enter') key = 'ENTER';
        else if (key === 'Backspace') key = '⌫';
        else key = key.toUpperCase();
    
        handleKeyPress(key);
    });

    const copyBtn = document.getElementById('copyResults');
    if (copyBtn) {
        copyBtn.addEventListener('click', async () => {
            const share = buildEmojiResultsFromDOM();

            try {
                await navigator.clipboard.writeText(share);
                showToast("Results Copied to Clipboard");
            } catch (e) {
            }
        });
    }
});