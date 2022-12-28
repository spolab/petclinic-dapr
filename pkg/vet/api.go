package vet

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func Register(dapr client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res RegisterVetResponse
		//
		// Invoke the actor
		//
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}
		raw, err := dapr.InvokeActor(ctx, &client.InvokeActorRequest{ActorType: "vet", ActorID: id, Method: "register", Data: bytes})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
		//
		// Parse the response
		//
		err = json.Unmarshal(raw.Data, &res)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, &res)
	}
}
