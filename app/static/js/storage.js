// storage.js
import { getUserAsync } from "./auth.js";
import { updateUserStreak } from "./render.js";

async function initializeUser() {
	const fingerprint = readFingerprint();

	let cognitoSub = null;
    let cognitoUsername = null;
    
	try {
		const user = await getUserAsync();
		if (user?.profile?.sub) {
			cognitoSub = user.profile.sub;
		}
        if (user?.profile?.["cognito:username"]) {
			cognitoUsername = user.profile["cognito:username"];
		}
	} catch (err) {
		// silently ignore
	}

	const payload = {
		fingerprint: fingerprint || null,
		user_agent: navigator.userAgent || "",
		device_type: /Mobi|Android/i.test(navigator.userAgent) ? "mobile" : "desktop",
		os: detectOS(),
		...(cognitoSub && { cognito_sub: cognitoSub }),
        ...(cognitoUsername && { cognito_username: cognitoUsername })
	};

	try {
		const res = await fetch("/api/v1.0/get-user", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(payload)
		});

		if (!res.ok) return;

		const data = await res.json();

        // new user (first time no fingerprint)
        // console.log("isNewUser?: " + data.new_user);
        // console.log("fingerprint: " + fingerprint);
        // console.log("data.fingerprint: " + data.fingerprint);
		if (!fingerprint && data.new_user === true && data.fingerprint) {
			writeFingerprint(data.fingerprint);
		} 
        if (!isNullOrWhitespace(data.fingerprint) && !isNullOrWhitespace(cognitoSub) && data.fingerprint != fingerprint) {
            // user has a different fingerprint in the DB than what is set client-side
            writeFingerprint(data.fingerprint);
        }
        updateUserStreak(data);
	} catch (err) {
		// silent
	}
}

function isNullOrWhitespace(val) {
    return (
        val === null ||
        val === undefined ||
        (typeof val === "string" && val.trim() === "")
    );
}



/*
import { userManager } from './auth.js';

async function initializeUser() {
	const fingerprint = readFingerprint();

	let cognitoSub = null;
	try {
		const user = await userManager.getUser();
		if (user && !user.expired && user.id_token) {
			const payload = JSON.parse(atob(user.id_token.split('.')[1]));
			if (payload.sub) {
				cognitoSub = payload.sub;
			}
		}
	} catch (err) {
		// fail silently
	}

	const payload = {
		cognito_sub: cognitoSub,
		fingerprint: fingerprint || null,
		user_agent: navigator.userAgent || "",
		device_type: /Mobi|Android/i.test(navigator.userAgent) ? "mobile" : "desktop",
		os: detectOS()
	};

	try {
		const res = await fetch("/api/v1.0/get-user", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(payload)
		});

		if (!res.ok) {
			return;
		}

		const data = await res.json();

		if (data.found === true) {
			// no update to fingerprint
		} else if (data.fingerprint) {
			writeFingerprint(data.fingerprint);
		}
	} catch (err) {
		// silent
	}
}*/

function detectOS() {
	const ua = navigator.userAgent;
	if (navigator.userAgentData?.platform) return navigator.userAgentData.platform;
	if (/Windows/i.test(ua)) return "Windows";
	if (/Mac OS/i.test(ua)) return "macOS";
	if (/Android/i.test(ua)) return "Android";
	if (/iPhone|iPad|iPod/i.test(ua)) return "iOS";
	if (/Linux/i.test(ua)) return "Linux";
	return "unknown";
}

function readFingerprint() {
	try {
		return localStorage.getItem("user_fingerprint");
	} catch (err) {
		return null;
	}
}

function writeFingerprint(uuid) {
	try {
		localStorage.setItem("user_fingerprint", uuid);
	} catch (err) {
	}
}

async function loadTodayGameState(fingerprint) {
	let gameState = null;

	if (fingerprint) {
		try {
			const res = await fetch("/api/v1.0/load-today-game", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ fingerprint })
			});

			if (res.ok) {
				const data = await res.json();
				gameState = data.gameState;
			}
		} catch (err) {
            // console.log(err);
		}
	}

	return gameState;
}

async function insertGameState(fingerprint, gameState) {
	if (fingerprint && gameState) {
		try {
			await fetch("/api/insert-game", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					fingerprint: fingerprint,
					state: gameState
				})
			});
		} catch (err) {
		}
	}
}


export { initializeUser };
export { readFingerprint };
export { loadTodayGameState };
export { insertGameState };