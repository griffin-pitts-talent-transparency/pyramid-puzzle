// auth.js
import {
    UserManager,
    WebStorageStateStore
} from "/js/vendor/oidc-client-ts.js";

import { updateAuthButtons } from './render.js';

const redirectUri = window.location.origin + "/";

const cognitoAuthConfig = {
    authority: "https://cognito-idp.us-east-2.amazonaws.com/us-east-2_iclUrXbex",
    client_id: "5i5gonu79adh8o7rstbdd2ajdd",
    redirect_uri: redirectUri, // "https://localhost:8080/" // "https://lebron-games.com/",
    silent_redirect_uri: window.location.origin + "/silent-signin.html",
    response_type: "code",
    scope: "openid email phone"
};

const userManager = new UserManager({
    ...cognitoAuthConfig,
    automaticSilentRenew: false,
    userStore: new WebStorageStateStore({ store: window.localStorage}),
    silentRequestTimeoutInSeconds: 20
});

async function signInRedirect() {
    try {
        await userManager.signinRedirect();
    } catch (err) {
        console.error("signinRedirect failed:", err);
    }
}

async function signOutRedirect() {
    await userManager.removeUser();
    await userManager.clearStaleState();

    const logoutUri = window.location.origin + "/";

    const cognitoDomain = "https://us-east-2iclurxbex.auth.us-east-2.amazoncognito.com";
    window.location.href = `${cognitoDomain}/logout?client_id=${cognitoAuthConfig.client_id}&logout_uri=${encodeURIComponent(logoutUri)}`;
    // https://us-east-2iclurxbex.auth.us-east-2.amazoncognito.com/logout?client_id=5i5gonu79adh8o7rstbdd2ajdd&logout_uri=https://localhost:8080/
}

let pendingUserPromise = null;
async function getUserAsync() {
    let userPromise = pendingUserPromise;

    if (!userPromise) {
        pendingUserPromise = (async () => {
            const url = new URL(window.location.href);
            const hasAuthCode = url.searchParams.has("code");
            const hasState = url.searchParams.has("state");
            let user = null;

            // Handle redirect callback (if coming back from AWS login)
            if (hasAuthCode && hasState) {
                try {
                    user = await userManager.signinCallback();
                    const redirectUrl = sessionStorage.getItem("post_signin_redirect") || "/";
                    sessionStorage.removeItem("post_signin_redirect");
                    window.location.replace(redirectUrl);
                    return null; // execution stops due to navigation
                } catch (err) {
                    console.warn("signinCallback failed:", err);
                }
            } else {
                user = await userManager.getUser();

                if (user && user.expired) {
                    try {
                        user = await userManager.signinSilent();
                        console.info("Silent signin successful");
                    } catch (err) {
                        console.warn("Silent signin failed:", err);
                        await userManager.removeUser();
                        user = null;
                    }
                }
            }

            return user;
        })();

        userPromise = pendingUserPromise;
    }

    const result = await userPromise;
    pendingUserPromise = null;

    updateAuthButtons(result);
    return result;
}


/*
    Main
*/
document.addEventListener("DOMContentLoaded", async () => {
    const signInButton = document.getElementById("signInButton");
    if (signInButton) {
        signInButton.addEventListener("click", () => {
            sessionStorage.setItem(
                "post_signin_redirect",
                window.location.href
            );
            window.location.href = "/signin";
        });
    }
    const signOutButton = document.getElementById("signOutButton");
    if (signOutButton) {
        signOutButton.addEventListener("click", async () => {
            await signOutRedirect();
        });
    }
});

export { getUserAsync };
export { signInRedirect };