// storage.js
async function initializeUser() {
	const fingerprint = readFingerprint();

	const payload = {
		fingerprint: fingerprint || null, // Send null if not found
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
		} else if (data.fingerprint) {
			writeFingerprint(data.fingerprint);
		}
	} catch (err) {
	}
}

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