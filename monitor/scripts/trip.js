import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 10,
    duration: '30s',
};

export default function () {
    const url = 'http://localhost:13071/trips';
    const payload = JSON.stringify({
        "passengerId": "passenger-123",
        "pickUpLocation": {
            "latitude": 10.77,
            "longitude": 106.69
        }
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer passenger-123',
        },
    };

    const res = http.post(url, payload, params);
    check(res, {
        'is status 200': (r) => r.status === 200,
        'trip id is not empty': (r) => r.json('id') !== '',
    });
    sleep(1);
}

