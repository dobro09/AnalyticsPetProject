class Analytics {

    constructor() {
        this.endpoint = "http://localhost:8080/api/events";

        this.sessionId = this.getOrCreateSessionId();
    }

    getOrCreateSessionId() {
        let sessionId = sessionStorage.getItem(
            "analytics_session_id"
        );

        if (!sessionId) {
            sessionId = crypto.randomUUID();

            sessionStorage.setItem(
                "analytics_session_id",
                sessionId
            );
        }

        return sessionId;
    }

    track(eventType, element, data = {}) {

        const event = {
            event_type: eventType,

            element:
                element.dataset.analyticsName ||
                element.id ||
                element.tagName.toLowerCase(),

            timestamp: new Date().toISOString(),

            page: window.location.pathname,

            session_id: this.sessionId,

            data: data
        };

        console.log("Analytics event:", event);

        fetch(this.endpoint, {
            method: "POST",

            headers: {
                "Content-Type": "application/json"
            },

            body: JSON.stringify(event)
        })
        .then(response => {

            if (!response.ok) {
                throw new Error(
                    `Analytics request failed: ${response.status}`
                );
            }

        })
        .catch(error => {
            console.error(
                "Failed to send analytics event:",
                error
            );
        });
    }
}


const analytics = new Analytics();