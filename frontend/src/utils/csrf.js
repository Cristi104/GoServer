import ENDPOINT_URL from "./config.js"

async function getCSRFToken() {
    const res = fetch(ENDPOINT_URL + "/api/csrf", {
        method: "GET",
        credentials: "include",
    })

    return res.headers.get("X-CSRF-Token")
}

export default getCSRFToken;
