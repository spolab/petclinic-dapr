import http from 'k6/http';
import { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.2/index.js';

export default function() {
    describe('Register a vet', () => {
        let request = {
            name: "name",
            surname: "surname",
            phone: "phone",
            email: "email@email.com",
        }            
        const res = http.put("http://vet.petclinic.127.0.0.1.nip.io/CT000O", request)
        expect(res.status, 'status').to.equal(203)
    })
}

// Use this only against a local development cluster
export const options = {
    insecureSkipTLSVerify: true,
};