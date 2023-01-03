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
package api

import (
	"encoding/json"
	"io"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/command"
	"github.com/spolab/petstore/pkg/event"
)

const (
	KeyStateDeleted = "deleted"
)

type VetSnapshot struct {
	Id      string `json:"id"`
	Surname string `json:"surname"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Version int    `json:"version"`
}

func Register(dapr client.Client, broker string, topic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()
		log.Debug().Str("id", id).Msg("begin register")
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			String(w, http.StatusBadRequest, err.Error())
			return
		}
		//
		// Invoke the actor method
		//
		log.Debug().Str("id", id).Str("payload", string(bytes)).Msg("executing actor method")
		raw, err := dapr.InvokeActor(ctx, &client.InvokeActorRequest{ActorType: "vet", ActorID: id, Method: "Register", Data: bytes})
		if err != nil {
			String(w, http.StatusInternalServerError, err.Error())
			log.Error().Str("id", id).Err(err).Msg("invoking actor method")

			return
		}
		//
		// Parse the response
		//
		log.Debug().Str("id", id).Msg("parsing the response")
		var res command.ActorResponse
		err = JsonFromBytes(raw.Data, &res)
		if err != nil {
			String(w, http.StatusInternalServerError, err.Error())
			log.Error().Str("id", id).Err(err).Msg("parsing the response")
			return
		}
		//
		// If the command is OK, take the events sent and publish them
		//
		if res.Status == command.StatusOK {
			for _, event := range res.Events {
				if err := dapr.PublishEvent(ctx, broker, topic, event, client.PublishEventWithRawPayload()); err != nil {
					String(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
		NoContent(w, http.StatusAccepted)
		log.Debug().Str("id", id).Msg("END register")
	}
}

// GetAll returns all the active vets
func GetActive(dapr client.Client, store string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		res, err := dapr.QueryStateAlpha1(ctx, store, "{}", nil)
		if err != nil {
			String(w, http.StatusInternalServerError, err.Error())
			return
		}
		log.Debug().Int("items", len(res.Results)).Msg("response from query")
		response := []*VetSnapshot{}
		for i, item := range res.Results {
			snapshot := &VetSnapshot{}
			if err := json.Unmarshal(item.Value, snapshot); err != nil {
				String(w, http.StatusInternalServerError, err.Error())
				return
			}
			response = append(response, snapshot)
			log.Debug().Str("payload", string(item.Value)).Int("index", i).Msg("response")
		}
		JSON(w, http.StatusOK, response)
	}
}

// Reads the events of interest and updates the read caches as required
func OnEvent(dapr client.Client, store string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapper, err := cloudevents.NewEventFromHTTPRequest(r)
		if err != nil {
			NoContent(w, http.StatusInternalServerError)
			return
		}
		switch wrapper.Type() {
		case event.EventVetRegisteredV1:
			//
			// Unwrap the event
			//
			var event event.VetRegistered
			if err := wrapper.DataAs(&event); err != nil {
				NoContent(w, http.StatusInternalServerError)
				return
			}
			//
			// Serialize state
			//
			snapshot := VetSnapshot{Id: event.Id, Surname: event.Surname, Name: event.Name, Email: event.Email, Phone: event.Phone}
			bytes, err := json.Marshal(&snapshot)
			if err != nil {
				NoContent(w, http.StatusInternalServerError)
				return
			}
			//
			// Store the snapshot
			//
			meta := make(map[string]string)
			meta[KeyStateDeleted] = "false"
			dapr.SaveState(r.Context(), store, event.Id, bytes, meta)
		}
	}
}
