import ws from 'k6/ws';
import { check, sleep, group } from 'k6';
import { Counter, Trend } from 'k6/metrics';

// Custom metrics
const wsConnectionErrors = new Counter('ws_connection_errors');
const wsMessageErrors = new Counter('ws_message_errors');
const wsMessageLatency = new Trend('ws_message_latency');
const wsConnectionDuration = new Trend('ws_connection_duration');

export const options = {
    stages: [
        { duration: '30s', target: 10 },  // Ramp up to 10 connections over 30s
        { duration: '1m30s', target: 10 }, // Stay at 10 connections for 1m30s
        { duration: '20s', target: 0 },    // Ramp down to 0 connections
    ],
    thresholds: {
        'ws_connection_errors': ['count < 5'],
        'ws_message_errors': ['count < 10'],
        'ws_message_latency': ['p(95) < 1000'], // 95th percentile should be less than 1000ms
    },
};

export default function () {
    const url = 'ws://localhost:3000/ws';

    return group('WebSocket Publish Handler', () => {
        const startTime = new Date();
        let messagesSent = 0;

        const res = ws.connect(url, null, (socket) => {
            socket.on('open', () => {
                console.log(`Connection ${__VU} opened`);
            });

            socket.on('message', (msg) => {
                const latency = new Date() - startTime;
                wsMessageLatency.add(latency);

                check(msg, {
                    'response received': (m) => m !== null && m.length > 0,
                });
            });

            socket.on('close', () => {
                console.log(`Connection ${__VU} closed`);
            });

            socket.on('error', (err) => {
                console.log(`Connection ${__VU} error: ${err}`);
                wsConnectionErrors.add(1);
            });

            // Send location data
            for (let i = 0; i < 5; i++) {
                const locationData = generateLocationData();

                try {
                    socket.send(JSON.stringify(locationData));
                    messagesSent++;

                    check(locationData, {
                        'location has id': (l) => l.id !== null && l.id.length > 0,
                        'location has userId': (l) => l.userId !== null && l.userId.length > 0,
                        'location has tripId': (l) => l.tripId !== null && l.tripId.length > 0,
                        'latitude is valid': (l) => l.latitude >= -90 && l.latitude <= 90,
                        'longitude is valid': (l) => l.longitude >= -180 && l.longitude <= 180,
                    });
                } catch (err) {
                    console.log(`Failed to send message: ${err}`);
                    wsMessageErrors.add(1);
                }

                sleep(0.5); // Wait 500ms between messages
            }

            // Wait for responses
            socket.setTimeout(() => {
                socket.close();
            }, 5000);
        });

        const connectionDuration = new Date() - startTime;
        wsConnectionDuration.add(connectionDuration);

        check(res, {
            'connection status 1000': (r) => r && r.status === 1000,
            'messages sent': () => messagesSent === 5,
        });

        if (res.status !== 1000) {
            wsConnectionErrors.add(1);
        }
    });
}

function generateLocationData() {
    const userId = `user_${Math.floor(Math.random() * 1000)}`;
    const tripId = `trip_${Math.floor(Math.random() * 100)}`;
    const id = `location_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

    return {
        id: id,
        latitude: (Math.random() * 180) - 90,
        longitude: (Math.random() * 360) - 180,
        userId: userId,
        tripId: tripId,
    };
}

