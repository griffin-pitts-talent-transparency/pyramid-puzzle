// game.js
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
    showToast,
    updateStatsGraphBar
} from './render.js';

let currentGuess = []; // For in-progress letter input

async function initializeGameSession() {
    // const user = await getUserAsync();
	await initializeUser(); // add user as a param

	const fingerprint = readFingerprint();
    const startHint = document.getElementById("startHint");

	if (fingerprint) {
		const loadedState = await loadTodayGameState(fingerprint);
		if (loadedState) {
            let gameStatus = "playing";
            const hasAnyGuesses =
                Array.isArray(loadedState.guesses) &&
                loadedState.guesses.length > 0;

            if (startHint) {
                if (hasAnyGuesses) startHint.classList.add("hidden");
                else startHint.classList.remove("hidden");
            }

			for (const guessResult of loadedState.guesses) {
                let count = 0;
                let isCorrect = true;
                if (guessResult && Array.isArray(guessResult)) {
                    count = guessResult.length;
                    renderGuessResult(guessResult, false);
            
                    isCorrect = isGuessCorrect(guessResult);
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
                    currentGuess = [];
                } else {
                    await Promise.all([
                        revealLossHintAndShowMessage(fingerprint),
                        loadAndDisplayStatsGraph(true)
                    ]);
                    initializeShareToggle();
                }

            } else if (gameStatus === "won") {
                showEndMessage("Game over; you WON!");
                await loadAndDisplayStatsGraph(true);
                initializeShareToggle();
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
            const startHint = document.getElementById("startHint");
            startHint.classList.add("hidden");

            const data = await res.json();
            const guessResult = data.guessResult;

            colorKeyboardKeys(guessResult);
            renderGuessResult(guessResult, true);
            renderGuessCounter(data.remainingGuesses);
            
            if(data.remainingGuesses > 0) {
                let count = guessResult.length;;
                let isCorrect = isGuessCorrect(guessResult);

                if(isCorrect && count === 7) {
                    showEndMessage("Game over; you WON!");
                    await loadAndDisplayStatsGraph();
                    initializeShareToggle();
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
                await Promise.all([
                    revealLossHintAndShowMessage(fingerprint),
                    loadAndDisplayStatsGraph(true)
                ]);
                initializeShareToggle();
            }
        }
    }
}

/*
    Helpers
*/
async function revealLossHintAndShowMessage(fingerprint) {
	try {
		const revealResp = await fetch("/api/v1.0/reveal-next-word", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ fingerprint })
		});

		if (revealResp.ok) {
			const revealData = await revealResp.json();
			const solutionWord = extractWord(revealData.nextWord);
			showEndMessage(`Game over; you lost... The word was "${solutionWord}"`);
		} else {
			showEndMessage("Game over; you lost...");
		}
	} catch {
		showEndMessage("Game over; you lost...");
	}
}

function isGuessCorrect(guessResult) {
	let result = false;
    if (Array.isArray(guessResult)) {
    	result = guessResult.every(letterObj => letterObj.status === "match");
    }
    return result;
}

async function loadAndDisplayStatsGraph(today) {
    const date = today ? getTodayNY() : ""; // blank = all-time
    const data = await fetchStats(date);
    updateStatsGraphBar("bar4Solved", "bar4Unsolved", data.pct4);
    updateStatsGraphBar("bar5Solved", "bar5Unsolved", data.pct5);
    updateStatsGraphBar("bar6Solved", "bar6Unsolved", data.pct6);
    updateStatsGraphBar("bar7Solved", "bar7Unsolved", data.pct7);
}

function getTodayNY() {
    const ny = new Date().toLocaleString("en-US", { timeZone: "America/New_York" });
    const d = new Date(ny);
    return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,"0")}-${String(d.getDate()).padStart(2,"0")}`;
}

async function fetchStats(date) {
    try {
        const res = await fetch("/api/v1.0/get-solved-tier-stats", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ date })
        });
        if (!res.ok) {
            throw new Error(await res.text());
        }
        return await res.json();
    } catch (err) {
        throw err;
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
            await safeHandleKeyPress(keyValue);
        });
    });
}

function extractWord(letterResultArray) {
    if (!Array.isArray(letterResultArray)) return "";
    return letterResultArray.map(x => x.letter).join("").toUpperCase();
}

function initializeShareToggle() {
    const sharePanel = document.getElementById("shareResultsSection");
    const openBtn = document.getElementById("openShare");
    const closeBtn = document.getElementById("closeShare");

    if (!sharePanel || !openBtn || !closeBtn) return;

    // Start in visible state
    sharePanel.classList.remove("hidden");
    openBtn.classList.add("hidden");

    openBtn.addEventListener("click", () => {
        sharePanel.classList.remove("hidden");
        openBtn.classList.add("hidden");
    });

    closeBtn.addEventListener("click", () => {
        sharePanel.classList.add("hidden");
        openBtn.classList.remove("hidden");
    });
}

function getNYOffsetHours() {
    // Get current time in UTC and NY
    const now = new Date();

    // Get offset from UTC in minutes for New York
    const nyOffsetMinutes = new Date(now.toLocaleString("en-US", {
        timeZone: "America/New_York"
    })).getTimezoneOffset();

    // Convert minutes to hours
    return -nyOffsetMinutes / 60;
}

function getNextMidnightNY() {
    const nowUTC = new Date();
    const NY_OFFSET_HOURS = getNYOffsetHours();

    const localOffsetMs = NY_OFFSET_HOURS * 60 * 60 * 1000;
    const nyNow = new Date(nowUTC.getTime() + localOffsetMs);

    const nyMidnight = new Date(nyNow);
    nyMidnight.setUTCHours(0, 0, 0, 0);
    nyMidnight.setUTCDate(nyNow.getUTCDate() + 1);

    const nyMidnightUTC = new Date(nyMidnight.getTime() - localOffsetMs);

    return nyMidnightUTC;
}

function startCountdownTimer() {
    const updateCountdown = () => {
        const now = new Date();
        const nextReset = getNextMidnightNY();
        const diff = nextReset - now;

        const el = document.getElementById("puzzleResetTimer");
        if (!el) return;

        if (diff <= 0) {
            el.textContent = "New puzzle available!";
            return;
        }

        const h = Math.floor(diff / 3600000);
        const m = Math.floor((diff % 3600000) / 60000);
        const s = Math.floor((diff % 60000) / 1000);

        el.textContent = `Next puzzle in: ${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
    };

    updateCountdown();
    setInterval(updateCountdown, 1000);
}


/*
    main
*/
window.addEventListener('DOMContentLoaded', async () => {
    await loadPartial('game-container', '/partials/game.html');
    
    // Initialization
    startCountdownTimer();
    await initializeGameSession();
    initializeInstructionsModal();
    initializeKeyboard();
    
    window.addEventListener('keydown', (e) => {
        let key = e.key;
    
        // Normalize desktop keyboard inputs
        if (key === 'Enter') key = 'ENTER';
        else if (key === 'Backspace') key = '⌫';
        else key = key.toUpperCase();
        safeHandleKeyPress(key);
    });

    /*
        Initialize share buttons
    */
    const copyBtn = document.getElementById("copyToClipboard");
    const smsBtn = document.getElementById("smsShare");
    const twitterBtn = document.getElementById("twitterShare");
    const facebookBtn = document.getElementById("facebookShare");
    // const instagramBtn = document.getElementById("instagramShare");

    const buildShareText = () => buildEmojiResultsFromDOM();
    const encoded = () => encodeURIComponent(buildShareText());

    if (copyBtn) {
        copyBtn.addEventListener("click", () => {
            navigator.clipboard.writeText(buildShareText());
            showToast("Copied!");
        });
    }

    if (smsBtn) {
        smsBtn.addEventListener("click", () => {
            window.location.href = `sms:?body=${encoded()}`;
        });
    }
    
    if (twitterBtn) {
        twitterBtn.addEventListener("click", () => {
            window.open(
                `https://twitter.com/intent/tweet?text=${encoded()}`,
                "_blank"
            );
        });
    }
    
    if (facebookBtn) {
        facebookBtn.addEventListener("click", () => {
            window.open(
                `https://www.facebook.com/sharer/sharer.php?u=https://lebron-games.com&quote=${encoded()}`,
                "_blank"
            );
        });
    }

    document.querySelectorAll('.toggle-option').forEach(button => {
        button.addEventListener('click', async () => {
            const selectedView = button.dataset.view;
            if (selectedView === 'today') {
                await loadAndDisplayStatsGraph(true);
            } else if (selectedView === 'all-time') {
                await loadAndDisplayStatsGraph(false);
            }

            // Update active state
            document.querySelectorAll('.toggle-option').forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
        });
    });
});