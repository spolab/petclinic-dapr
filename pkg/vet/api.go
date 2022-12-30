package vet

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func Register(dapr client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var res RegisterVetResponse
		ctx := r.Context()
		log.Debug().Str("id", id).Msg("START register")
		//
		// Invoke the actor
		//
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			log.Error().Str("id", id).Err(err).Msg("reading request body")
			return
		}
		raw, err := dapr.InvokeActor(ctx, &client.InvokeActorRequest{ActorType: "vet", ActorID: id, Method: "register", Data: bytes})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			log.Error().Str("id", id).Err(err).Msg("invoking actor")
			return
		}
		//
		// Parse the response
		//
		err = json.Unmarshal(raw.Data, &res)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			log.Error().Str("id", id).Err(err).Msg("parsing the response")

			return
		}
		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, &res)
		log.Debug().Str("id", id).Msg("END register")
	}
}
