import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { target: 5, duration: '1m'},
    { target: 15, duration: '1m' },
    { target: 20, duration: '1m' },
    { target: 50, duration: '1m' },
    { target: 50, duration: '2m' },
    { target: 0, duration: '1m' },
  ],
}

export default function () {
  const resp = http.get(`http://127.0.0.1:8765/${__ENV.URL_ALIAS}`, { follow: false, redirects: 0 })
  check(resp, {
    'status is 308': (r) => r.status == 308
  })
}
