// render.js
function initializeInstructionsModal() {
    const openBtn = document.getElementById('openInstructions');
    const closeBtn = document.getElementById('closeInstructions');
    const modal = document.getElementById('instructionsModal');

    if(openBtn && closeBtn && modal) {
        openBtn.addEventListener('click', () => {
            modal.classList.remove('hidden');
        });
    
        closeBtn.addEventListener('click', () => {
            modal.classList.add('hidden');
        });
    
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.classList.add('hidden');
            }
        });
    }
}

function renderActiveTiles(currentGuess) {
    if (!Array.isArray(currentGuess)) {
        currentGuess = [];
    }
    const activeRow = document.querySelector('.row.active-row');

    if (activeRow) {
        const tiles = Array.from(activeRow.querySelectorAll('.tile'));
        const maxLen = tiles.length;

        const safeGuess = currentGuess.slice(0, maxLen);

        for (let i = 0; i < maxLen; i++) {
            tiles[i].textContent = safeGuess[i] || '';
        }
    }
    
}

function setActiveRow() {
    document.activeElement?.blur();
    window.getSelection()?.removeAllRanges();
    const pyramid = document.querySelector('.pyramid');

    if (pyramid) {
        const rows = Array.from(pyramid.querySelectorAll('.row'));
        rows.forEach(r => r.classList.remove('active-row'));

        const emptyRows = rows.filter(row => {
            const tiles = Array.from(row.querySelectorAll('.tile'));
            return tiles.every(tile => tile.textContent.trim() === '');
        });

        if (emptyRows.length > 0) {
            const activeRow = emptyRows[0];
            activeRow.classList.add('active-row');
        }
    }
}

function renderGuessResult(guessResult) {
    // Support either [{letter,status}, ...] or { letters: [...] }
    let letters = null;

    if (Array.isArray(guessResult)) {
        letters = guessResult;
    } else if (guessResult && Array.isArray(guessResult.letters)) {
        letters = guessResult.letters;
    }

    if (letters && letters.length > 0) {
        const guessLen = letters.length;
        const rows = document.querySelectorAll('.pyramid .row');
        let targetRow = null;

        for (let i = 0; i < rows.length; i++) {
            const row = rows[i];
            if (row.classList.contains('guess-filled')) continue;

            const tiles = row.querySelectorAll('.tile');

            if (tiles.length === guessLen && targetRow === null) {
                targetRow = row;
            }
        }

        if (targetRow !== null) {
            const tiles = targetRow.querySelectorAll('.tile');

            for (let k = 0; k < letters.length; k++) {
                const letterObj = letters[k];
                const tile = tiles[k];

                tile.textContent = letterObj.letter;

                tile.classList.remove('match', 'present', 'miss');

                if (letterObj.status === 'match') {
                    tile.classList.add('match');
                } else if (letterObj.status === 'present') {
                    tile.classList.add('present');
                } else if (letterObj.status === 'miss') {
                    tile.classList.add('miss');
                }
            }

            // Mark this row as filled so it won't be reused
            targetRow.classList.add('guess-filled');
        }
    }
}

function addBlankRow(tileCount) {
    const pyramid = document.querySelector('.pyramid');
    const row = document.createElement('div');
    row.classList.add('row');

    for (let i = 0; i < tileCount; i++) {
        const tile = document.createElement('div');
        tile.classList.add('tile');
        row.appendChild(tile);
    }

    const allRows = Array.from(pyramid.querySelectorAll('.row'));
    const matchingRows = allRows.filter(r => r.children.length === tileCount);

    if (matchingRows.length > 0) {
        const lastMatch = matchingRows[matchingRows.length - 1];
        if (lastMatch.nextSibling) {
            pyramid.insertBefore(row, lastMatch.nextSibling);
        } else {
            pyramid.appendChild(row);
        }
    } else {
        pyramid.appendChild(row);
    }
}

function renderGuessCounter(remainingGuesses) {
    const counterDiv = document.querySelector('.counter');
    if (counterDiv) {
        if (remainingGuesses > 0) {
            counterDiv.textContent = `Incorrect Guesses Remaining: ${remainingGuesses}`;
        } else {
            counterDiv.textContent = 'YOU LOSE';
        }
    }
}


function colorKeyboardKeys(guessResult) {
    if (Array.isArray(guessResult)) {
        const keys = document.querySelectorAll('.keyboard .key');
        guessResult.forEach(({ letter, status }) => {
            const key = Array.from(keys).find(k => k.textContent.toUpperCase() === letter.toUpperCase());
            if (key) {
                const isMatch = key.classList.contains('match');
                const isPresent = key.classList.contains('present');
                const isMiss = key.classList.contains('miss');

                if (status === 'match') {
                    // Upgrade to match: remove present if it exists
                    if (!isMatch) {
                        key.classList.remove('present');
                        key.classList.add('match');
                    }
                } 
                else if (status === 'present') {
                    // Only add present if not already match or present
                    if (!isMatch && !isPresent) {
                        key.classList.remove('miss');
                        key.classList.add('present');
                    }
                } 
                else if (status === 'miss') {
                    // Only mark as miss if no other color is applied
                    if (!isMatch && !isPresent && !isMiss) {
                        key.classList.add('miss');
                    }
                }
            }
        });
    }
}

function clearKeyboardColors() {
    const keys = document.querySelectorAll('.keyboard .key');
    keys.forEach(key => {
        key.classList.remove('match', 'present', 'miss');
    });
}

function showEndMessage(text) {
    document.querySelectorAll('.row').forEach(r => r.classList.remove('active-row'));
    const counterDiv = document.querySelector('.counter');
    if (counterDiv) {
        counterDiv.textContent = '';

        const span = document.createElement('span');
        span.classList.add('end-message');
        span.textContent = text;

        counterDiv.appendChild(span);
    }

    // const copyBtn = document.getElementById('copyResults');
    // if (copyBtn) copyBtn.classList.remove('hidden');

    // Show the new share section
    const shareSection = document.getElementById('shareResultsSection');
    if (shareSection) shareSection.classList.remove('hidden');


    document.querySelector('.end-message')?.classList.add('pulse');

    const pyramid = document.querySelector('.pyramid');
    if (pyramid) pyramid.classList.add('flash-bg');

    document.querySelectorAll('.tile').forEach(tile => {
        tile.classList.add('tile-shake');
    });

    document.body.classList.add('screen-shake');
    setTimeout(() => {
        document.body.classList.remove('screen-shake');
    }, 1000);

    setTimeout(() => {
        if (pyramid) pyramid.classList.remove('flash-bg');
    }, 3000);
}

function buildEmojiResultsFromDOM() {
    const emojiMap = {
        match: "🟣",
        present: "🟠",
        miss: "⚪"
    };

    const rows = Array.from(document.querySelectorAll('.pyramid .row'));

    let result = `Pyramid Puzzle - ${new Date().toLocaleDateString("en-US")}`;

    const tierLengths = [4, 5, 6, 7];
    let achievedTier = 0;

    for (let i = 0; i < tierLengths.length; i++) {
        const length = tierLengths[i];

        const rowsOfThisTier = rows.filter(r => r.querySelectorAll('.tile').length === length);

        const solvedRow = rowsOfThisTier.some(row => {
            const tiles = Array.from(row.querySelectorAll('.tile'));
            return tiles.length === length &&
                    tiles.every(t => t.classList.contains('match'));
        });

        if (solvedRow) achievedTier = i + 1;
    }
    result += ` - ${achievedTier}/4\n`;

    result += 'https://lebron-games.com/\n';

    
    for (const row of rows) {
        const tiles = Array.from(row.querySelectorAll('.tile'));
        if (tiles.length === 0) continue;

        let line = "";

        for (const tile of tiles) {
            if (tile.classList.contains('match')) {
                line += emojiMap.match;
            } else if (tile.classList.contains('present')) {
                line += emojiMap.present;
            } else if (tile.classList.contains('miss')) {
                line += emojiMap.miss;
            } else {
                line += ""; // ignore blanks / incomplete rows
            }
        }

        if (line.length > 0) result += line + "\n";
    }

    return result.trim();
}

function showToast(message) {
    const toast = document.createElement('div');
    toast.textContent = message;
    toast.className = 'toast-message';
    // document.body.appendChild(toast);
    document.querySelector('.share-section').after(toast);

    // Start fade-out after 1.5 seconds
    setTimeout(() => {
        toast.classList.add('fade-out');
    }, 1500);

    // Remove from DOM after animation ends
    toast.addEventListener('transitionend', () => {
        toast.remove();
    });
}

function updateStatsGraphBar(idSolved, idUnsolved, solvedPctRaw) {
    const solved = document.getElementById(idSolved);
    const unsolved = document.getElementById(idUnsolved);
    const label = solved.querySelector(".bar-label");

    const solvedPct = Math.max(0, Math.min(100, solvedPctRaw));
    const unsolvedPct = 100 - solvedPct;

    solved.style.width = `${solvedPct}%`;
    unsolved.style.width = `${unsolvedPct}%`;

    if (solvedPct > 0) {
        label.textContent = `${Math.round(solvedPct)}%`;
        label.style.display = "block";

        // Decide if it fits inside
        if (solvedPct < 19) {
            label.classList.add("outside");
        } else {
            label.classList.remove("outside");
        }
    } else {
        label.textContent = "0%";
        label.style.display = "block";
        label.classList.add("outside");
    }
}


export { initializeInstructionsModal };
export { renderActiveTiles };
export { setActiveRow };
export { renderGuessResult };
export { addBlankRow };
export { renderGuessCounter };
export { colorKeyboardKeys };
export { clearKeyboardColors };
export { showEndMessage };
export { buildEmojiResultsFromDOM };
export { showToast };
export { updateStatsGraphBar };
