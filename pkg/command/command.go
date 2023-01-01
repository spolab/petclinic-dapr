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
package command

import cloudevents "github.com/cloudevents/sdk-go/v2"

const (
	StatusOK = iota
	StatusInvalid
	StatusError
)

type RegisterVetCommand struct {
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type ActorResponse struct {
	Status  int                  `json:"status"`
	Message string               `json:"message,omitempty"`
	Events  []*cloudevents.Event `json:"events,omitempty"`
}
