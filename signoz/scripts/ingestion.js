import ws from 'uber-clone/ws';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'uber-clone/metrics';

const messagesSent = new Counter('ws_messages_sent');
const messagesInvalid = new Counter('ws_messages_invalid');
const connectionErrors = new Rate('ws_connection_errors');
const messageLatency = new Trend('ws_message_latency');

export const options = {
    scenarios: {
        steady_drivers: {
            executor: 'constant-vus',
            vus: 10,
            duration: '30s',
            gracefulStop: '5s',
        },
        // rush_hour: {
        //     executor: 'ramping-vus',
        //     startVUs: 0,
        //     stages: [
        //         { duration: '10s', target: 50 },
        //         { duration: '20s', target: 50 },
        //         { duration: '10s', target: 0 },
        //     ],
        //     startTime: '35s',
        // },
    },
    thresholds: {
        ws_connection_errors: ['rate<0.1'],
        ws_messages_sent: ['count>100'],
        ws_message_latency: ['p(95)<100'],
    },
};

function randomLocation() {
    return {
        latitude: 10.762622 + (Math.random() - 0.5) * 0.1,
        longitude: 106.660172 + (Math.random() - 0.5) * 0.1,
        timestamp: new Date().toISOString(),
    };
}

function validDriverLocation(driverId, tripId) {
    return JSON.stringify({
        tripId: tripId,
        driverId: driverId,
        location: randomLocation(),
    });
}

function invalidMessage(type) {
    switch (type) {
        case 'malformed_json':
            return '{ invalid json }';
        case 'missing_fields':
            return JSON.stringify({ driverId: 'driver-1' });
        case 'wrong_types':
            return JSON.stringify({
                tripId: 123,  // should be string
                driverId: null,
                location: 'not an object',
            });
        default:
            return '';
    }
}

export default function () {
    const url = __ENV.WS_URL || 'ws://localhost:13701/ws';
    const driverId = `driver-${__VU}`;
    const tripId = `trip-${__VU}-${Date.now()}`;

    const params = {
        headers: { 'X-Source': 'k6-script' },
    };

    const res = ws.connect(url, params, function (socket) {
        socket.on('open', function () {
            console.log(`VU ${__VU}: Connected as ${driverId}`);

            // Send location updates periodically (every 1 second)
            const intervalId = setInterval(() => {
                
                const start = Date.now();
                const message = validDriverLocation(driverId, tripId);
                socket.send(message);
                messagesSent.add(1);
                messageLatency.add(Date.now() - start);
            }, 1000);

            // Send some invalid messages to test error handling (10% of VUs)
            if (__VU % 10 === 0) {
                setTimeout(() => {
                    socket.send(invalidMessage('malformed_json'));
                    messagesInvalid.add(1);
                }, 2000);

                setTimeout(() => {
                    socket.send(invalidMessage('missing_fields'));
                    messagesInvalid.add(1);
                }, 4000);
            }

            // Close connection after duration
            socket.setTimeout(function () {
                clearInterval(intervalId);
                console.log(`VU ${__VU}: Closing connection`);
                socket.close();
            }, 25000);
        });

        socket.on('message', function (data) {
            console.log(`VU ${__VU}: Received message: ${data}`);
        });

        socket.on('error', function (e) {
            console.error(`VU ${__VU}: Error: ${e.error()}`);
        });

        socket.on('close', function () {
            console.log(`VU ${__VU}: Connection closed`);
        });
    });

    const connected = check(res, {
        'WebSocket connection established (status 101)': (r) => r && r.status === 101,
    });

    if (!connected) {
        connectionErrors.add(1);
    }

    sleep(1); // Small pause between iterations
}