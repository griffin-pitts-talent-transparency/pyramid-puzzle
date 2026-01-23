// leaderboard.js
import { 
    renderSolvedWordsAllTime,
    renderTotalGuessCount,
    renderSolvedGamesAllTime
} from './render.js';

import { getUserAsync } from './auth.js';

async function fetchSolvedWordsAllTime() {
    try {
        const response = await fetch("/api/v1.0/read-solved-words-all-time", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (!response.ok) {
            console.error("Failed to fetch solved words:", response.statusText);
            return;
        }

        const data = await response.json();
        renderSolvedWordsAllTime(data);
    } catch (err) {
        console.error("Error fetching solved words:", err);
    }
}

async function fetchSolvedGamesAllTime() {
    try {
        const response = await fetch("/api/v1.0/read-solved-games-all-time", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (!response.ok) {
            console.error("Failed to fetch solved words:", response.statusText);
            return;
        }

        const data = await response.json();
        renderSolvedGamesAllTime(data);
    } catch (err) {
        console.error("Error fetching solved words:", err);
    }
}

async function fetchTotalGuessCount() {
    try {
        const response = await fetch("/api/v1.0/read-total-guess-count", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (!response.ok) {
            console.error("Failed to fetch total guess count:", response.statusText);
            return;
        }

        const data = await response.json();
        renderTotalGuessCount(data);
    } catch (err) {
        console.error("Error fetching total guess count:", err);
    }
}


/*
    Main
*/
window.addEventListener('DOMContentLoaded', async () => {
    await Promise.all([
        fetchSolvedWordsAllTime(),
        fetchSolvedGamesAllTime(),
        fetchTotalGuessCount(),
        getUserAsync()
    ]);
});
