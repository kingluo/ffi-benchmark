import http from 'k6/http';

export const options = {
  discardResponseBodies: true,
};

const data = open('/opt/test.data.64k', 'b');
//const data = open('/opt/test.data.1m', 'b');
//const data = open('/opt/test.data.10m', 'b');

export default function () {
  //http.get('http://envoy:10000/get', {headers:{Authorization:'Basic Zm9vbmFtZTp2YWxpZHBhc3N3b3Jk'}});
  http.post('http://envoy:10000/post', data, {
    headers: {
      'Content-Type': 'application/octet-stream',
      Authorization: 'Basic Zm9vbmFtZTp2YWxpZHBhc3N3b3Jk',
      Expect: '',
    }
  });
}

