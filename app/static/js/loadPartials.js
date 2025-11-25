// loadPartials.js

const VERSION = "2.0"; // bump on each deploy

async function loadPartial(id, url) {
    try {
        const versionedUrl = `${url}?v=${VERSION}`;
        const res = await fetch(versionedUrl, {
            cache: "no-store"   // Force fresh load
        });
        if (!res.ok) {
            throw new Error(`Failed to load ${versionedUrl}`);
        }
        const html = await res.text();
        document.getElementById(id).innerHTML = html;
    } catch (err) {
        console.error(err);
    }
}

export { loadPartial };
