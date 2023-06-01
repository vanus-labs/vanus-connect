import fetch from "node-fetch";

const sendHttpRequest = (requestDetails, endpoint) => {
    fetch(endpoint, {
        method: requestDetails.method,
        body: JSON.stringify(requestDetails.body),
        headers: requestDetails.headers,
    })
        .then((response) => {
            if (response.ok) {
                console.log("HTTP request sent successfully.");
            } else {
                console.error(
                    "Failed to send HTTP request:",
                    response.status,
                    response.statusText
                );
            }
        })
        .catch((error) => console.error("Failed to send HTTP request:", error));
};

export default sendHttpRequest;