package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

var validate *validator.Validate

var keyVals = make(map[string]string)

type setRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type deleteRequest struct {
	Key string `json:"key" validate:"required"`
}

func main() {
	validate = validator.New()
	r := chi.NewRouter()

	r.Post("/set", setValue)
	r.Get("/get", getValue)
	r.Post("/delete", deleteKey)

	log.Fatal(http.ListenAndServe(":4000", r))
}

func setValue(w http.ResponseWriter, r *http.Request) {

	req := setRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error validating setValue request body: %s", err.Error()), 400)
		return
	}

	keyVals[req.Key] = req.Value

	resp := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		req.Key,
		keyVals[req.Key],
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func getValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key cannot be empty", 400)
		return
	}
	resp := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   key,
		Value: keyVals[key],
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func deleteKey(w http.ResponseWriter, r *http.Request) {
	req := deleteRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error validating deleteKey request body: %s", err.Error()), 400)
		return
	}

	if req.Key == "" {
		http.Error(w, "Key cannot be empty", 400)
		return
	}

	if _, ok := keyVals[req.Key]; !ok {
		w.Write([]byte(fmt.Sprintf("Key does not exist")))
	} else {
		delete(keyVals, req.Key)
		resp := struct {
			Key string `json:"key"`
		}{
			req.Key,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
