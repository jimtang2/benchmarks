import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  duration: "30s",
  vus: 600,
  // stages: [
  //   { target: 150, duration: "10s" },
  //   { target: 600, duration: "10s" },
  //   { target: 1200, duration: "10s" }
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
