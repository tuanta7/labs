import http from "k6/http";

export const options = {
  vus: 50,
  duration: "30s",
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:6060";

export default function () {
  http.get(`${BASE_URL}/leak`);
}
