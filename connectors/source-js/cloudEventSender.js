const cloudEventSender = (eventDetails, endpoint) => {
    // Create a CloudEvent
    const cloudEvent = {
        type: eventDetails.type,
        source: eventDetails.source,
        data: JSON.stringify(eventDetails.data),
        dataContentType: "application/json",
        time: new Date().toISOString(),
    };

    // Serialize the CloudEvent to a JSON string
    const eventJson = JSON.stringify(cloudEvent);

    fetch(endpoint, {
        method: "POST",
        body: eventJson,
        headers: { "Content-Type": "application/cloudevents+json" },
    })
        .then((response) => {
            if (response.ok) {
                console.log("CloudEvent sent successfully.");
            } else {
                console.error(
                    "Failed to send CloudEvent:",
                    response.status,
                    response.statusText
                );
            }
        })
        .catch((error) => console.error("Failed to send CloudEvent:", error));
};

module.exports = cloudEventSender;