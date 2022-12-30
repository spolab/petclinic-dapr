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
package vet

import (
	"encoding/json"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/common/parse"
	"github.com/spolab/petstore/pkg/common/respond"
)

func Register(dapr client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()
		var cmd RegisterVetCommand
		log.Debug().Str("id", id).Msg("begin register")
		//
		// Parse the request
		//
		log.Debug().Str("id", id).Msg("reading request payload")
		if err := parse.JsonFromReader(r.Body, &cmd); err != nil {
			respond.String(w, http.StatusBadRequest, err.Error())
			log.Error().Str("id", id).Err(err).Msg("reading request body")
			return
		}
		//
		// Invoke the actor method
		//
		log.Debug().Str("id", id).Msgf("marshalling command %v", cmd)
		bytes, err := json.Marshal(&cmd)
		if err != nil {
			respond.String(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Debug().Str("id", id).Str("payload", string(bytes)).Msg("executing actor method")
		raw, err := dapr.InvokeActor(ctx, &client.InvokeActorRequest{ActorType: "vet", ActorID: id, Method: "register", Data: bytes})
		if err != nil {
			respond.String(w, http.StatusInternalServerError, err.Error())
			log.Error().Str("id", id).Err(err).Msg("invoking actor method")
			return
		}
		//
		// Parse the response
		//
		log.Debug().Str("id", id).Msg("parsing the response")
		var res RegisterVetResponse
		err = parse.JsonFromBytes(raw.Data, &res)
		if err != nil {
			respond.String(w, http.StatusInternalServerError, err.Error())
			log.Error().Str("id", id).Err(err).Msg("parsing the response")
			return
		}
		respond.JSON(w, http.StatusAccepted, &res)
		log.Debug().Str("id", id).Msg("END register")
	}
}
