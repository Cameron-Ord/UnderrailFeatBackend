package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/calculation"
	"main/db"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	endpoint := vars["endpoint"]
	switch endpoint {

	case "get-all-builds":

	case "get-user-builds":

	case "save-build":

	case "signup":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}

		var signupQuery db.SignupData
		if err := json.Unmarshal(body, &signupQuery); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	case "login":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}

		var loginQuery db.LoginData
		if err := json.Unmarshal(body, &loginQuery); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	case "calculate":

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}

		var data calculation.RequestData
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}
		returnedFeats, err := calculation.PrepareData(data)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(returnedFeats)
		return
	default:
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/{endpoint}", handler).Methods("POST")
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"POST"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)
	http.Handle("/", corsHandler(r))
	port := "6969"
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
