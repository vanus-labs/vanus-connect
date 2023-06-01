const cloudEventSender = require("./cloudEventSender");

btn.addEventListener("click", function sendEvent() {
    const eventDetails = {
        type: "com.example.event",
        source: "/example/source",
        data: { message: "Hello, CloudEvent!" },
    };

    const endpoint = "<your_endpoint_url>";

    cloudEventSender(eventDetails, endpoint);
});