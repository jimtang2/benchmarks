import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  // stages: [
  //   { duration: '30s', target: 1000 },   
  //   { duration: '30s', target: 2000 },   
  //   { duration: '30s', target: 0 },      
  // ],
  duration: "10s",
  vus: 100,
  thresholds: {
    http_req_duration: ['p(95)<100'],   // 95% of requests under 100ms
    errors: ['rate<0.01'],              // <1% errors
  },
};

export default function () {
  const res = http.get(__ENV.URL);
  check(res, { 'statusOk': (r) => r.status === 200 }) || errorRate.add(1);
  sleep(0.01);
}