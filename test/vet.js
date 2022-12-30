/*
Copyright 2022 Alessandro Santini

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
        const res = http.put("http://vet.petclinic.127.0.0.1.nip.io/CT000O", JSON.stringify(request))
        expect(res.status, 'status').to.equal(203)
    })
}

// Use this only against a local development cluster
export const options = {
    insecureSkipTLSVerify: true,
};