import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  duration: "3m",
  vus: 200,
  // stages: [
  //   { target: 1000, duration: "1m" },
  //   { target: 3000, duration: "2m" },
  //   { target: 500, duration: "1m" }
  // ], 
  thresholds: {
    // http_req_duration: ['p(95)<100'],
    errors: ['rate<0.01'],
  },
};

export default function () {
  const res = http.get(__ENV.URL);
  check(res, { 'check:status-ok': (r) => r.status === 200 }) || errorRate.add(1);
  sleep(0.005);
}
