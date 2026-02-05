import http from "k6/http";
import {check, sleep} from "k6";

export const options = {
    vus: 100,
    duration: "30s",
    thresholds: {
        http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
        http_req_failed: ["rate<0.01"],   // Less than 1% of requests should fail
    },
};

const BASE_URL = "http://localhost:9092";
const DRIVER_IDS = ["1", "2", "3", "4", "5"];

export default function () {
    const driverId = DRIVER_IDS[Math.floor(Math.random() * DRIVER_IDS.length)];
    const res = http.get(`${BASE_URL}/drivers/${driverId}`, {
        tags: {name: "GetDriverByID"},
    });

    check(res, {
        "response time < 500ms": (r) => r.timings.duration < 500,
        "status is 200": (r) => r.status === 200,
        "has driver data": (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.id !== undefined && data.name !== undefined;
            } catch {
                return false;
            }
        },

    });

    sleep(1);  // Simulate user think time
}
